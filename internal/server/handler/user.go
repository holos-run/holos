package handler

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/user"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewUserHandler returns a new UserService implementation.
func NewUserHandler(db *ent.Client) *UserHandler {
	return &UserHandler{db: db}
}

// UserHandler implements the UserService interface.
type UserHandler struct {
	db *ent.Client
}

func (h *UserHandler) GetUserClaims(
	ctx context.Context,
	req *connect.Request[holos.GetUserClaimsRequest],
) (*connect.Response[holos.GetUserClaimsResponse], error) {
	authnIdentity, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	res := connect.NewResponse(&holos.GetUserClaimsResponse{
		Iss:           authnIdentity.Issuer(),
		Sub:           authnIdentity.Subject(),
		Email:         authnIdentity.Email(),
		EmailVerified: authnIdentity.Verified(),
		Name:          authnIdentity.Name(),
	})
	return res, nil
}

func (h *UserHandler) GetUser(
	ctx context.Context,
	req *connect.Request[holos.GetUserRequest],
) (*connect.Response[holos.GetUserResponse], error) {
	email := req.Msg.GetEmail()
	if email == "" {
		id, err := authn.FromContext(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
		email = id.Email()
	}
	dbUser, err := getUser(ctx, h.db, email)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
	}

	res := connect.NewResponse(&holos.GetUserResponse{User: UserToRPC(dbUser)})
	return res, nil
}

func (h *UserHandler) RegisterUser(
	ctx context.Context,
	req *connect.Request[holos.RegisterUserRequest],
) (*connect.Response[holos.RegisterUserResponse], error) {
	log := logger.FromContext(ctx).With("issue", 127)
	oidc, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	var name string
	if req.Msg.Name != nil {
		name = req.Msg.GetName()
	} else {
		name = oidc.Name()
	}

	var createdUser *ent.User
	err = WithTx(ctx, h.db, func(tx *ent.Tx) error {
		createdUser, err = createUser(ctx, tx.Client(), name, oidc)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "could not save transaction", "err", err)
		return nil, err
	}

	log = log.With("user.id", createdUser.ID, "user.name", createdUser.Name)
	log.InfoContext(ctx, "registered user", "event", "registration")

	res := connect.NewResponse(&holos.RegisterUserResponse{
		User: UserToRPC(createdUser),
	})
	return res, nil
}

// UserToRPC returns an *holos.User adapted from *ent.User u.
func UserToRPC(u *ent.User) *holos.User {
	iamUser := holos.User{
		Id:    u.ID.String(),
		Email: u.Email,
		Name:  u.Name,
		Timestamps: &holos.Timestamps{
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		},
	}
	return &iamUser
}

func getUser(ctx context.Context, client *ent.Client, email string) (*ent.User, error) {
	log := logger.FromContext(ctx)
	user, err := client.User.Query().Where(user.Email(email)).Only(ctx)
	if err != nil {
		log.DebugContext(ctx, "could not get user", "err", err, "email", email)
		return nil, errors.Wrap(err)
	}
	return user, nil
}

func createUser(ctx context.Context, client *ent.Client, name string, claims authn.Identity) (*ent.User, error) {
	log := logger.FromContext(ctx)
	// Create the user, error if it already exists
	user, err := client.User.
		Create().
		SetEmail(claims.Email()).
		SetIss(claims.Issuer()).
		SetSub(claims.Subject()).
		SetName(name).
		Save(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		log.ErrorContext(ctx, "could not create user", "err", err)
		return user, err
	}

	log = log.With("user", user)
	log.InfoContext(ctx, "created user")

	return user, nil
}

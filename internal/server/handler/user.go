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

func (h *UserHandler) GetCallerClaims(
	ctx context.Context,
	req *connect.Request[holos.GetCallerClaimsRequest],
) (*connect.Response[holos.GetCallerClaimsResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	res := connect.NewResponse(&holos.GetCallerClaimsResponse{
		Claims: &holos.Claims{
			Iss:           authnID.Issuer(),
			Sub:           authnID.Subject(),
			Email:         authnID.Email(),
			EmailVerified: authnID.Verified(),
			Name:          authnID.Name(),
			Groups:        authnID.Groups(),
			GivenName:     authnID.GivenName(),
			FamilyName:    authnID.FamilyName(),
			Picture:       authnID.Picture(),
		},
	})
	return res, nil
}

func (h *UserHandler) GetCallerUser(
	ctx context.Context,
	req *connect.Request[holos.GetCallerUserRequest],
) (*connect.Response[holos.GetCallerUserResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	dbUser, err := getUser(ctx, h.db, authnID.Issuer(), authnID.Subject())
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
	}
	res := connect.NewResponse(&holos.GetCallerUserResponse{User: UserToRPC(dbUser)})
	return res, nil
}

func (h *UserHandler) CreateCallerUser(
	ctx context.Context,
	req *connect.Request[holos.CreateCallerUserRequest],
) (*connect.Response[holos.CreateCallerUserResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	var createdUser *ent.User
	err = WithTx(ctx, h.db, func(tx *ent.Tx) error {
		createdUser, err = createUser(ctx, tx.Client(), authnID.Name(), authnID)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "could not save transaction", "err", err)
		return nil, err
	}

	res := connect.NewResponse(&holos.CreateCallerUserResponse{
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

func getUser(ctx context.Context, client *ent.Client, iss string, sub string) (*ent.User, error) {
	log := logger.FromContext(ctx)
	user, err := client.User.Query().
		Where(
			user.Iss(iss),
			user.Sub(sub),
		).
		Only(ctx)
	if err != nil {
		log.DebugContext(ctx, "could not get user", "err", err, "iss", iss, "sub", sub)
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

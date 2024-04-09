package handler

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/server/ent"
	"github.com/holos-run/holos/internal/server/ent/useridentity"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/pkg/errors"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
)

func createUser(ctx context.Context, client *ent.Client, name string, claims authn.Identity) (*ent.User, error) {
	log := logger.FromContext(ctx).With("issue", 127)
	// Create the user, error if it already exists
	user, err := client.User.
		Create().
		SetName(name).
		SetEmail(claims.Email()).
		SetEmailVerified(claims.Verified()).
		Save(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		log.ErrorContext(ctx, "could not create user", "err", err)
		return user, err
	}

	log = log.With("user", user)
	log.DebugContext(ctx, "created user")

	link, err := client.UserIdentity.Create().
		SetIss(claims.Issuer()).
		SetSub(claims.Subject()).
		SetEmail(claims.Email()).
		SetEmailVerified(claims.Verified()).
		SetName(claims.Name()).
		SetUserID(user.ID).
		Save(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		log.ErrorContext(ctx, "could not create link", "err", err)
		return user, err
	}
	log = log.With("link", link)
	log.DebugContext(ctx, "created link")
	return user, nil
}

func updateUser(ctx context.Context, user *ent.User, name *string) (*ent.User, error) {
	if name == nil || *name == user.Name {
		return user, nil
	}
	updated, err := user.Update().SetName(*name).Save(ctx)
	return updated, errors.Wrap(err)
}

func (h *HolosHandler) RegisterUser(
	ctx context.Context,
	req *connect.Request[holos.RegisterUserRequest],
) (*connect.Response[holos.RegisterUserResponse], error) {
	log := logger.FromContext(ctx).With("issue", 127)
	oidc, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	// Check if the user exists
	entUser, err := h.db.UserIdentity.
		Query().
		Where(
			useridentity.And(
				useridentity.Iss(oidc.Issuer()),
				useridentity.Sub(oidc.Subject()),
			),
		).
		QueryUser().
		Only(ctx)
	if err == nil {
		log.DebugContext(ctx, "already registered", "status", "exists")
		user, err := updateUser(ctx, entUser, req.Msg.Name)
		if err != nil {
			slog.ErrorContext(ctx, "could not update", "err", err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		res := connect.NewResponse(
			&holos.RegisterUserResponse{
				User:          UserToRPC(user),
				AlreadyExists: true,
			})
		return res, nil
	}
	var notFoundError *ent.NotFoundError
	if !errors.As(err, &notFoundError) {
		err = connect.NewError(connect.CodeInternal, errors.Wrap(err))
		log.ErrorContext(ctx, "could not register", "err", err)
		return nil, err
	}
	log.DebugContext(ctx, "not found needs to be created", "err", err)

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

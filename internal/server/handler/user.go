package handler

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
)

func createUser(ctx context.Context, client *ent.Client, name string, claims authn.Identity) (*ent.User, error) {
	log := logger.FromContext(ctx)
	// Create the user, error if it already exists
	user, err := client.User.
		Create().
		SetEmail(claims.Email()).
		SetIss(claims.Issuer()).
		SetSub(claims.Subject()).
		SetName(claims.Name()).
		Save(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		log.ErrorContext(ctx, "could not create user", "err", err)
		return user, err
	}

	log = log.With("user", user)
	log.DebugContext(ctx, "created user")

	return user, nil
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

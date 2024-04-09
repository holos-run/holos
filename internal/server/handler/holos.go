package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/server/ent"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	holos "github.com/holos-run/holos/internal/server/service/gen/holos/v1alpha1"
	"github.com/holos-run/holos/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewHolosHandler(db *ent.Client) *HolosHandler {
	return &HolosHandler{db: db}
}

// HolosHandler implements the connect service handler interface.
type HolosHandler struct {
	db *ent.Client
}

func (h *HolosHandler) GetUserClaims(
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

package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/user"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
)

// NewPlatformHandler returns a new PlatformService implementation.
func NewPlatformHandler(db *ent.Client) *PlatformHandler {
	return &PlatformHandler{db: db}
}

// PlatformHandler implements the PlatformService interface.
type PlatformHandler struct {
	db *ent.Client
}

func (h *PlatformHandler) GetPlatforms(
	ctx context.Context,
	req *connect.Request[holos.GetPlatformsRequest],
) (*connect.Response[holos.GetPlatformsResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	dbUser, err := h.db.User.Query().
		Where(
			user.Iss(authnID.Issuer()),
			user.Sub(authnID.Subject()),
		).
		WithOrganizations(func(q *ent.OrganizationQuery) {
			q.WithPlatforms()
		}).
		Only(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	reqOrgID, err := uuid.FromString(req.Msg.OrgId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	var reqDBOrg *ent.Organization
	wantOrgIDs := make([]uuid.UUID, 0, len(dbUser.Edges.Organizations))
	for _, org := range dbUser.Edges.Organizations {
		if org.ID == reqOrgID {
			reqDBOrg = org
			break
		}
		wantOrgIDs = append(wantOrgIDs, org.ID)
	}

	if reqDBOrg == nil {
		return nil, connect.NewError(
			connect.CodeNotFound,
			fmt.Errorf(
				"org_id not found\n\thave (%+v)\n\twant (%+v)",
				reqOrgID,
				wantOrgIDs))
	}

	rpcPlatforms := make([]*holos.Platform, 0, len(reqDBOrg.Edges.Platforms))

	for _, platform := range reqDBOrg.Edges.Platforms {
		rpcPlatforms = append(rpcPlatforms, &holos.Platform{
			Id:          platform.ID.String(),
			Name:        platform.Name,
			DisplayName: platform.DisplayName,
		})
	}

	resp := connect.NewResponse(&holos.GetPlatformsResponse{
		OrgId:     reqDBOrg.ID.String(),
		Platforms: rpcPlatforms,
	})

	return resp, nil
}

func (h *PlatformHandler) AddPlatform(
	ctx context.Context,
	req *connect.Request[holos.AddPlatformRequest],
) (*connect.Response[holos.GetPlatformsResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	_ = authnID
	err = errors.Wrap(fmt.Errorf("not implemented"))
	return nil, connect.NewError(connect.CodeUnimplemented, err)
}

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
	"google.golang.org/protobuf/types/known/timestamppb"
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
	_, reqDBOrg, err := getAuthnUsersOrg(ctx, req.Msg.OrgId, h.db)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return getPlatformsResponse(reqDBOrg), nil
}

func (h *PlatformHandler) AddPlatform(
	ctx context.Context,
	req *connect.Request[holos.AddPlatformRequest],
) (*connect.Response[holos.GetPlatformsResponse], error) {
	dbUser, dbOrg, err := getAuthnUsersOrg(ctx, req.Msg.OrgId, h.db)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	platform, err := h.db.Platform.Create().
		SetOrgID(dbOrg.ID).
		SetCreatorID(dbUser.ID).
		SetName(req.Msg.Name).
		SetDisplayName(req.Msg.DisplayName).
		Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	resp := getPlatformsResponse(dbOrg)
	resp.Msg.Platforms = append(resp.Msg.Platforms, PlatformToRPC(platform))

	return resp, nil
}

func PlatformToRPC(platform *ent.Platform) *holos.Platform {
	return &holos.Platform{
		Id:          platform.ID.String(),
		Name:        platform.Name,
		DisplayName: platform.DisplayName,
		Timestamps: &holos.Timestamps{
			CreatedAt: timestamppb.New(platform.CreatedAt),
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		},
		Creator: &holos.Creator{
			Id: platform.CreatorID.String(),
		},
	}
}

func getAuthnUsersOrg(ctx context.Context, orgID string, db *ent.Client) (*ent.User, *ent.Organization, error) {
	reqOrgID, err := uuid.FromString(orgID)
	if err != nil {
		err = errors.Wrap(fmt.Errorf("%w:\n\thave (%+v)\n\twant (uuid string)", err, orgID))
		return nil, nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	dbUser, err := db.User.Query().
		Where(
			user.Iss(authnID.Issuer()),
			user.Sub(authnID.Subject()),
		).
		WithOrganizations(func(q *ent.OrganizationQuery) {
			q.WithPlatforms()
		}).
		Only(ctx)
	if err != nil {
		return nil, nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
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
		return nil, nil, connect.NewError(
			connect.CodePermissionDenied,
			fmt.Errorf(
				"not an organization member\n\thave (%+v)\n\twant (%+v)",
				reqOrgID,
				wantOrgIDs))
	}

	return dbUser, reqDBOrg, nil
}

func getPlatformsResponse(reqDBOrg *ent.Organization) *connect.Response[holos.GetPlatformsResponse] {
	// one extra in case a new platform is appended.
	rpcPlatforms := make([]*holos.Platform, 0, 1+len(reqDBOrg.Edges.Platforms))
	for _, platform := range reqDBOrg.Edges.Platforms {
		rpcPlatforms = append(rpcPlatforms, PlatformToRPC(platform))
	}

	return connect.NewResponse(&holos.GetPlatformsResponse{
		OrgId:     reqDBOrg.ID.String(),
		Platforms: rpcPlatforms,
	})
}

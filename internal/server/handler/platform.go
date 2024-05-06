package handler

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/organization"
	entplatform "github.com/holos-run/holos/internal/ent/platform"
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

func (h *PlatformHandler) ListPlatforms(
	ctx context.Context,
	req *connect.Request[holos.ListPlatformsRequest],
) (*connect.Response[holos.ListPlatformsResponse], error) {
	_, reqDBOrg, err := getAuthnUsersOrg(ctx, req.Msg.OrgId, h.db)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	resp := &holos.ListPlatformsResponse{Platforms: rpcPlatforms(reqDBOrg)}
	return connect.NewResponse(resp), nil
}

func (h *PlatformHandler) AddPlatform(
	ctx context.Context,
	req *connect.Request[holos.AddPlatformRequest],
) (*connect.Response[holos.AddPlatformResponse], error) {
	dbUser, dbOrg, err := getAuthnUsersOrg(ctx, req.Msg.Platform.OrgId, h.db)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	platform, err := h.db.Platform.Create().
		SetOrgID(dbOrg.ID).
		SetCreatorID(dbUser.ID).
		SetName(req.Msg.Platform.Name).
		SetDisplayName(req.Msg.Platform.DisplayName).
		Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	resp := &holos.AddPlatformResponse{Platforms: rpcPlatforms(dbOrg)}
	resp.Platforms = append(resp.Platforms, PlatformToRPC(platform))

	return connect.NewResponse(resp), nil
}

func (h *PlatformHandler) getPlatform(ctx context.Context, id string, uid authn.Identity) (*ent.Platform, error) {
	platformID, err := uuid.FromString(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	p, err := h.db.Platform.Query().
		Where(entplatform.ID(platformID)).
		Where(entplatform.HasOrganizationWith(
			organization.HasUsersWith(
				user.Iss(uid.Issuer()),
				user.Sub(uid.Subject()),
			))).
		Only(ctx)
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
	}

	return p, nil
}

func (h *PlatformHandler) GetPlatform(ctx context.Context, req *connect.Request[holos.GetPlatformRequest]) (*connect.Response[holos.GetPlatformResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	p, err := h.getPlatform(ctx, req.Msg.PlatformId, authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return connect.NewResponse(&holos.GetPlatformResponse{Platform: PlatformToRPC(p)}), nil
}

func (h *PlatformHandler) PutModel(ctx context.Context, req *connect.Request[holos.PutModelRequest]) (*connect.Response[holos.PutModelResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	p, err := h.getPlatform(ctx, req.Msg.GetPlatformId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	slog.WarnContext(ctx, "todo: validate the platform config against cue definitions", "action", "todo", "cue", len(p.Cue))

	_, err = p.Update().
		SetModel(&holos.Model{Model: req.Msg.GetModel()}).
		Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	return connect.NewResponse(&holos.PutModelResponse{Model: req.Msg.Model}), nil
}

func (h *PlatformHandler) GetModel(ctx context.Context, req *connect.Request[holos.GetModelRequest]) (*connect.Response[holos.GetModelResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	p, err := h.getPlatform(ctx, req.Msg.PlatformId, authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return connect.NewResponse(&holos.GetModelResponse{Model: p.Model.Model}), nil
}

func (h *PlatformHandler) GetForm(ctx context.Context, req *connect.Request[holos.GetFormRequest]) (*connect.Response[holos.GetFormResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	p, err := h.getPlatform(ctx, req.Msg.GetPlatformId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return connect.NewResponse(&holos.GetFormResponse{Fields: p.Form.GetFields(), Model: p.Model.GetModel()}), nil
}

func (h *PlatformHandler) PutForm(ctx context.Context, req *connect.Request[holos.PutFormRequest]) (*connect.Response[holos.PutFormResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	p, err := h.getPlatform(ctx, req.Msg.GetPlatformId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	_, err = p.Update().
		SetForm(&holos.Form{Fields: req.Msg.GetFields()}).
		Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	resp := &holos.PutFormResponse{Fields: req.Msg.GetFields()}
	return connect.NewResponse(resp), nil
}

func PlatformToRPC(platform *ent.Platform) *holos.Platform {
	return &holos.Platform{
		Id:          platform.ID.String(),
		Name:        platform.Name,
		DisplayName: platform.DisplayName,
		OrgId:       platform.OrgID.String(),
		Spec:        &holos.PlatformSpec{Model: platform.Model.GetModel()},
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

	// Check the user is a member of the organization.
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

func rpcPlatforms(reqDBOrg *ent.Organization) []*holos.Platform {
	if reqDBOrg == nil {
		return nil
	}
	// one extra in case a new platform is appended.
	platforms := make([]*holos.Platform, 0, 1+len(reqDBOrg.Edges.Platforms))
	for _, platform := range reqDBOrg.Edges.Platforms {
		platforms = append(platforms, PlatformToRPC(platform))
	}
	return platforms
}

package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/organization"
	entplatform "github.com/holos-run/holos/internal/ent/platform"
	"github.com/holos-run/holos/internal/ent/user"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	storage "github.com/holos-run/holos/service/gen/holos/storage/v1alpha1"
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

// CreatePlatform implements the PlatformService CreatePlatform rpc method.
func (h *PlatformHandler) CreatePlatform(ctx context.Context, req *connect.Request[platform.CreatePlatformRequest]) (*connect.Response[platform.CreatePlatformResponse], error) {
	if req == nil || req.Msg == nil || req.Msg.Platform == nil || req.Msg.Platform.Owner == nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing platform owner org")))
	}
	orgID := req.Msg.Platform.Owner.GetOrgId()
	dbUser, dbOrg, err := getAuthnUsersOrg(ctx, orgID, h.db)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	builder := h.db.Platform.Create().
		SetOrgID(dbOrg.ID).
		SetCreatorID(dbUser.ID).
		SetName(req.Msg.Platform.Name).
		SetDisplayName(req.Msg.Platform.GetDisplayName())
	entPlatform, err := builder.Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	resp := &platform.CreatePlatformResponse{
		Platform: PlatformToRPC(entPlatform),
	}

	return connect.NewResponse(resp), nil
}

// ListPlatforms implements the PlatformService ListPlatforms rpc method.
func (h *PlatformHandler) ListPlatforms(ctx context.Context, req *connect.Request[platform.ListPlatformsRequest]) (*connect.Response[platform.ListPlatformsResponse], error) {
	if req == nil || req.Msg == nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no message in request")))
	}
	_, reqDBOrg, err := getAuthnUsersOrg(ctx, req.Msg.OrgId, h.db)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	resp := &platform.ListPlatformsResponse{Platforms: rpcPlatforms(reqDBOrg)}
	return connect.NewResponse(resp), nil
}

// GetPlatform implements the PlatformService GetPlatform rpc method.
func (h *PlatformHandler) GetPlatform(ctx context.Context, req *connect.Request[platform.GetPlatformRequest]) (*connect.Response[platform.GetPlatformResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	p, err := h.getPlatform(ctx, req.Msg.GetPlatformId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// JEFFTODO: FieldMask in the PlatformToRPC() method.

	return connect.NewResponse(&platform.GetPlatformResponse{Platform: PlatformToRPC(p)}), nil
}

// getEditor ensures the user identity stored in the context is a member of the
// organization.  Useful to get the editor uuid for mutations.  orgID must be a
// valid uuid string.
func getEditor(ctx context.Context, db *ent.Client, authnID authn.Identity, orgID string) (*ent.User, error) {
	orgUUID, err := uuid.FromString(orgID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	editor, err := db.User.Query().
		Where(
			user.And(
				user.Iss(authnID.Issuer()),
				user.Sub(authnID.Subject()),
				user.HasOrganizationsWith(organization.ID(orgUUID)),
			),
		).Only(ctx)
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
	}
	return editor, nil
}

func (h *PlatformHandler) UpdatePlatform(ctx context.Context, req *connect.Request[platform.UpdatePlatformRequest]) (*connect.Response[platform.UpdatePlatformResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	editor, err := getEditor(ctx, h.db, authnID, req.Msg.GetPlatform().GetOwner().GetOrgId())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	p, err := h.getPlatform(ctx, req.Msg.GetPlatform().GetId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// JEFFTODO: Field Mask
	builder := p.Update().
		SetName(req.Msg.GetPlatform().GetName()).
		SetDisplayName(req.Msg.GetPlatform().GetDisplayName()).
		SetForm(&storage.Form{
			Fields: req.Msg.GetPlatform().GetSpec().GetForm().GetFields(),
		}).
		SetModel(&storage.Model{
			Model: req.Msg.GetPlatform().GetSpec().GetModel(),
		}).
		SetEditor(editor)

	entPlatform, err := builder.Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	resp := &platform.UpdatePlatformResponse{
		Platform: PlatformToRPC(entPlatform),
	}

	return connect.NewResponse(resp), nil

}

// getPlatform returns a platform by id ensuring the request comes from an
// identity that is a member of the organization owning the platform.
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

func PlatformToRPC(entity *ent.Platform) *platform.Platform {
	return &platform.Platform{
		Id: entity.ID.String(),
		Owner: &platform.Owner{
			Owner: &platform.Owner_OrgId{
				OrgId: entity.OrgID.String(),
			},
		},
		Name:        entity.Name,
		DisplayName: &entity.DisplayName,
		Spec: &platform.Spec{
			Model: entity.Model.GetModel(),
			Form: &platform.Form{
				Fields: entity.Form.GetFields(),
			},
		},
		Detail: &object.Detail{
			CreatedBy: &object.ResourceEditor{
				Editor: &object.ResourceEditor_UserId{
					UserId: entity.CreatedByID.String(),
				},
			},
			CreatedAt: timestamppb.New(entity.CreatedAt),
			UpdatedBy: &object.ResourceEditor{
				Editor: &object.ResourceEditor_UserId{
					UserId: entity.UpdatedByID.String(),
				},
			},
			UpdatedAt: timestamppb.New(entity.UpdatedAt),
		},
	}
}

func getAuthnUsersOrg(ctx context.Context, orgID string, db *ent.Client) (*ent.User, *ent.Organization, error) {
	reqOrgID, err := uuid.FromString(orgID)
	if err != nil {
		err = errors.Wrap(fmt.Errorf("invalid org id: %w:\n\thave (%+v)\n\twant (uuid string)", err, orgID))
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

func rpcPlatforms(reqDBOrg *ent.Organization) []*platform.Platform {
	if reqDBOrg == nil {
		return nil
	}
	// one extra in case a new platform is appended.
	platforms := make([]*platform.Platform, 0, 1+len(reqDBOrg.Edges.Platforms))
	for _, platform := range reqDBOrg.Edges.Platforms {
		platforms = append(platforms, PlatformToRPC(platform))
	}
	return platforms
}

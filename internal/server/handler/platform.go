package handler

import (
	"context"
	"fmt"
	"reflect"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/organization"
	entplatform "github.com/holos-run/holos/internal/ent/platform"
	"github.com/holos-run/holos/internal/ent/user"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/strings"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	storage "github.com/holos-run/holos/service/gen/holos/storage/v1alpha1"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
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

// GetPlatform implements the PlatformService GetPlatform rpc method.
func (h *PlatformHandler) GetPlatform(ctx context.Context, req *connect.Request[platform.GetPlatformRequest]) (*connect.Response[platform.GetPlatformResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	dbPlatform, err := h.getPlatform(ctx, req.Msg.GetPlatformId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	rpcPlatform := PlatformToRPC(dbPlatform)

	mask, err := fieldmask_utils.MaskFromProtoFieldMask(req.Msg.GetFieldMask(), strings.PascalCase)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}
	var rpcPlatformMasked platform.Platform
	if err = fieldmask_utils.StructToStruct(mask, rpcPlatform, &rpcPlatformMasked); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	return connect.NewResponse(&platform.GetPlatformResponse{Platform: &rpcPlatformMasked}), nil
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

	mask, err := fieldmask_utils.MaskFromProtoFieldMask(req.Msg.GetFieldMask(), strings.PascalCase)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	platforms, err := rpcPlatforms(reqDBOrg, mask)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	resp := &platform.ListPlatformsResponse{Platforms: platforms}
	return connect.NewResponse(resp), nil
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

func (h *PlatformHandler) UpdatePlatform(
	ctx context.Context,
	req *connect.Request[platform.UpdatePlatformRequest],
) (*connect.Response[platform.UpdatePlatformResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	// Update mask is required to avoid an older client accidentally writing over
	// fields added to the update operation.
	if req.Msg.GetUpdateMask() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(errors.New("missing update mask")))
	}

	// Refer to https://github.com/mennanov/fieldmask-utils/blob/v1.1.2/README.md#naming-function
	mask, err := fieldmask_utils.MaskFromProtoFieldMask(req.Msg.GetUpdateMask(), strings.PascalCase)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	ops := make(map[string]interface{})
	if err := fieldmask_utils.StructToMap(mask, req.Msg.GetUpdate(), ops, fieldmask_utils.WithMapVisitor(newVisitor(ctx))); err != nil {
		return nil, errors.Wrap(err)
	}

	p, err := h.getPlatform(ctx, req.Msg.GetUpdate().GetPlatformId(), authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	log := logger.FromContext(ctx).With("platform_id", p.ID.String(), "org_id", p.OrgID.String())

	if len(ops) == 0 {
		err = errors.New("nothing to do: provide fields to update in the mask")
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	editor, err := getEditor(ctx, h.db, authnID, p.OrgID.String())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	log = log.With("op", "update", "editor", editor.Email)

	builder := p.Update()
	builder.SetEditor(editor)
	for field := range ops {
		log := log.With("field", field)
		switch field {
		case "Name":
			name := req.Msg.GetUpdate().GetName()
			log.InfoContext(ctx, "update", field, name)
			builder.SetName(name)
		case "DisplayName":
			name := req.Msg.GetUpdate().GetDisplayName()
			log.InfoContext(ctx, "update", field, name)
			builder.SetDisplayName(name)
		case "Model":
			log.InfoContext(ctx, "update")
			builder.SetModel(&storage.Model{Model: req.Msg.GetUpdate().GetModel()})
		case "Form":
			log.InfoContext(ctx, "update")
			builder.SetForm(&storage.Form{FieldConfigs: req.Msg.GetUpdate().GetForm().GetFieldConfigs()})
		default:
			err := errors.Wrap(errors.New("could not update: unknown field " + field))
			log.ErrorContext(ctx, "could not update", "err", err)
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	dbPlatform, err := builder.Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	resp := platform.UpdatePlatformResponse{
		Platform: PlatformToRPC(dbPlatform),
	}

	return connect.NewResponse(&resp), nil
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
			return nil, connect.NewError(connect.CodeNotFound, err)
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, err)
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
				FieldConfigs: entity.Form.GetFieldConfigs(),
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

func rpcPlatforms(reqDBOrg *ent.Organization, mask fieldmask_utils.Mask) ([]*platform.Platform, error) {
	if reqDBOrg == nil {
		return nil, nil
	}

	// one extra in case a new platform is appended.
	platforms := make([]*platform.Platform, 0, 1+len(reqDBOrg.Edges.Platforms))
	for _, dbPlatform := range reqDBOrg.Edges.Platforms {
		var platformMasked platform.Platform
		if err := fieldmask_utils.StructToStruct(mask, PlatformToRPC(dbPlatform), &platformMasked); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
		}
		platforms = append(platforms, &platformMasked)
	}
	return platforms, nil
}

// newVisitor returns a new fieldmask visitor function for use with
// fieldmask_utils.StructToMap.  This is here largely as an placeholder to
// remember we can mutate the value if we want.
func newVisitor(ctx context.Context) func(filter fieldmask_utils.FieldFilter, src, dst reflect.Value, srcFieldName, dstFieldName string, srcFieldValue reflect.Value) fieldmask_utils.MapVisitorResult {
	log := logger.FromContext(ctx)
	return func(filter fieldmask_utils.FieldFilter, src, dst reflect.Value, srcFieldName, dstFieldName string, srcFieldValue reflect.Value) fieldmask_utils.MapVisitorResult {
		log.DebugContext(ctx, "visitor", "srcFieldName", srcFieldName, "dstFieldName", dstFieldName)
		return fieldmask_utils.MapVisitorResult{
			UpdatedDst: &dst,
		}
	}
}

package handler

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"unicode"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/user"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	holos "github.com/holos-run/holos/service/gen/holos/organization/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewOrganizationHandler returns a new OrganizationService implementation.
func NewOrganizationHandler(db *ent.Client) *OrganizationHandler {
	return &OrganizationHandler{db: db}
}

// OrganizationHandler implements the OrganizationService interface.
type OrganizationHandler struct {
	db *ent.Client
}

func (h *OrganizationHandler) ListOrganizations(ctx context.Context, req *connect.Request[holos.ListOrganizationsRequest]) (*connect.Response[holos.ListOrganizationsResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	dbUser, err := h.db.User.Query().
		Where(
			user.Iss(authnID.Issuer()),
			user.Sub(authnID.Subject()),
		).
		WithOrganizations().
		Only(ctx)
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
	}

	rpcOrgs := make([]*holos.Organization, 0, len(dbUser.Edges.Organizations))
	for _, dbOrg := range dbUser.Edges.Organizations {
		rpcOrgs = append(rpcOrgs, OrganizationToRPC(dbOrg))
	}

	// JEFFTODO: FieldMask
	res := connect.NewResponse(&holos.ListOrganizationsResponse{
		User: &object.UserRef{
			User: &object.UserRef_UserId{
				UserId: dbUser.ID.String(),
			},
		},
		Organizations: rpcOrgs,
	})
	return res, nil
}

func (h *OrganizationHandler) CreateOrganization(
	ctx context.Context,
	req *connect.Request[holos.CreateOrganizationRequest],
) (*connect.Response[holos.CreateOrganizationResponse], error) {
	log := logger.FromContext(ctx)
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	dbUser, err := getUser(ctx, h.db, authnID)
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
	}

	var dbOrg *ent.Organization
	err = WithTx(ctx, h.db, func(tx *ent.Tx) (err error) {
		dbOrg, err = tx.Organization.Create().
			SetName(cleanAndAppendRandom(authnID.Name())).
			SetDisplayName(authnID.GivenName() + "'s Org").
			SetCreatorID(dbUser.ID).
			SetEditorID(dbUser.ID).
			Save(ctx)
		if err != nil {
			return err
		}
		return tx.Organization.UpdateOne(dbOrg).AddUsers(dbUser).Exec(ctx)
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err))
	}
	log = log.With("organization", dbOrg)

	log.InfoContext(ctx, "created organization")

	res := connect.NewResponse(&holos.CreateOrganizationResponse{
		Organization: OrganizationToRPC(dbOrg),
	})
	return res, nil
}

func cleanAndAppendRandom(s string) string {
	mapping := func(r rune) rune {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}
		return -1
	}
	cleaned := strings.Map(mapping, s)
	randNum := rand.Intn(900_000) + 100_000
	return fmt.Sprintf("%s-%06d", cleaned, randNum)
}

// OrganizationToRPC returns an *holos.Organization adapted from *ent.Organization u.
func OrganizationToRPC(org *ent.Organization) *holos.Organization {
	orgID := org.ID.String()
	rpcEntity := holos.Organization{
		OrgId:       &orgID,
		Name:        org.Name,
		DisplayName: &org.DisplayName,
		Detail: &object.Detail{
			CreatedBy: &object.ResourceEditor{
				Editor: &object.ResourceEditor_UserId{
					UserId: org.CreatedByID.String(),
				},
			},
			CreatedAt: timestamppb.New(org.CreatedAt),
			UpdatedBy: &object.ResourceEditor{
				Editor: &object.ResourceEditor_UserId{
					UserId: org.UpdatedByID.String(),
				},
			},
			UpdatedAt: timestamppb.New(org.UpdatedAt),
		},
	}
	return &rpcEntity
}

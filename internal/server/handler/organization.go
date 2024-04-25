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
	"github.com/holos-run/holos/internal/server/middleware/authn"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
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

func (h *OrganizationHandler) GetCallerOrganizations(
	ctx context.Context,
	req *connect.Request[holos.GetCallerOrganizationsRequest],
) (*connect.Response[holos.GetCallerOrganizationsResponse], error) {
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

	res := connect.NewResponse(&holos.GetCallerOrganizationsResponse{
		User:          UserToRPC(dbUser),
		Organizations: rpcOrgs,
	})
	return res, nil
}

func (h *OrganizationHandler) CreateCallerOrganization(
	ctx context.Context,
	req *connect.Request[holos.CreateCallerOrganizationRequest],
) (*connect.Response[holos.GetCallerOrganizationsResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	// todo get user by iss, sub
	dbUser, err := getUser(ctx, h.db, authnID.Email())
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		} else {
			return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		}
	}

	var org *ent.Organization
	err = WithTx(ctx, h.db, func(tx *ent.Tx) (err error) {
		org, err = h.db.Organization.Create().
			SetName(cleanAndAppendRandom(authnID.Name())).
			SetDisplayName(authnID.GivenName() + "'s Org").
			SetCreatorID(dbUser.ID).
			Save(ctx)
		if err != nil {
			return err
		}
		dbUser, err = dbUser.Update().
			AddOrganizations(org).
			Save(ctx)
		return err
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err))
	}

	// TODO: prefetch organizations
	dbOrgs, err := dbUser.QueryOrganizations().All(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err))
	}
	rpcOrgs := make([]*holos.Organization, 0, len(dbOrgs))
	for _, dbOrg := range dbOrgs {
		rpcOrgs = append(rpcOrgs, OrganizationToRPC(dbOrg))
	}

	res := connect.NewResponse(&holos.GetCallerOrganizationsResponse{
		User:          UserToRPC(dbUser),
		Organizations: rpcOrgs,
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
	rpcEntity := holos.Organization{
		Id:          org.ID.String(),
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Timestamps: &holos.Timestamps{
			CreatedAt: timestamppb.New(org.CreatedAt),
			UpdatedAt: timestamppb.New(org.UpdatedAt),
		},
	}
	return &rpcEntity
}

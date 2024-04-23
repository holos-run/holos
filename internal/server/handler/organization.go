package handler

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/organization"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
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

func (h *OrganizationHandler) GetOrganization(
	ctx context.Context,
	req *connect.Request[holos.GetOrganizationRequest],
) (*connect.Response[holos.GetOrganizationResponse], error) {
	name := req.Msg.GetName()

	entity, err := h.db.Organization.Query().Where(organization.Name(name)).Only(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
	}

	res := connect.NewResponse(&holos.GetOrganizationResponse{Organization: OrganizationToRPC(entity)})
	return res, nil
}

func (h *OrganizationHandler) RegisterOrganization(
	ctx context.Context,
	req *connect.Request[holos.RegisterOrganizationRequest],
) (*connect.Response[holos.RegisterOrganizationResponse], error) {
	log := logger.FromContext(ctx)
	creator, err := getAuthenticatedUser(ctx, h.db)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	var createdOrganization *ent.Organization
	err = WithTx(ctx, h.db, func(tx *ent.Tx) error {
		createdOrganization, err = createOrganization(ctx, tx.Client(), req.Msg.GetName(), req.Msg.GetDisplayName(), creator.ID)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "could not save transaction", "err", err)
		return nil, err
	}

	log = log.With("organization", createdOrganization)
	log.InfoContext(ctx, "registered", "event", "registration", "resource", "organization")

	res := connect.NewResponse(&holos.RegisterOrganizationResponse{
		Organization: OrganizationToRPC(createdOrganization),
	})
	return res, nil
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

func createOrganization(ctx context.Context, client *ent.Client, name string, displayName string, creatorID uuid.UUID) (*ent.Organization, error) {
	log := logger.FromContext(ctx)
	// Create the user, error if it already exists
	entity, err := client.Organization.
		Create().
		SetName(name).
		SetDisplayName(displayName).
		SetCreatorID(creatorID).
		Save(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
		log.ErrorContext(ctx, "could not create user", "err", err)
		return entity, err
	}

	log = log.With("organization", entity)
	log.InfoContext(ctx, "created")

	return entity, nil
}

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/ent/user"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/strings"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	org "github.com/holos-run/holos/service/gen/holos/organization/v1alpha1"
	storage "github.com/holos-run/holos/service/gen/holos/storage/v1alpha1"
	holos "github.com/holos-run/holos/service/gen/holos/user/v1alpha1"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewUserHandler returns a new UserService implementation.
func NewUserHandler(db *ent.Client) *UserHandler {
	return &UserHandler{db: db}
}

// UserHandler implements the UserService interface.
type UserHandler struct {
	db *ent.Client
}

func (h *UserHandler) GetUser(ctx context.Context, req *connect.Request[holos.GetUserRequest]) (*connect.Response[holos.GetUserResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	if req.Msg.GetUser() != nil || req.Msg.GetFieldMask() != nil {
		return nil, connect.NewError(connect.CodeUnimplemented, errors.Wrap(fmt.Errorf("not implemented: make an empty request instead")))
	}

	dbUser, err := getUser(ctx, h.db, authnID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	rpcUser := UserToRPC(dbUser)
	if pic := authnID.Picture(); pic != "" {
		rpcUser.Picture = &pic
	}

	return connect.NewResponse(&holos.GetUserResponse{User: rpcUser}), nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *connect.Request[holos.CreateUserRequest]) (*connect.Response[holos.CreateUserResponse], error) {
	_, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	var createdUser *ent.User
	if rpcUser := req.Msg.GetUser(); rpcUser != nil {
		createdUser, err = h.createUser(ctx, h.db, rpcUser)
	} else {
		createdUser, err = h.createCallerUser(ctx)
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}

	res := connect.NewResponse(&holos.CreateUserResponse{
		User: UserToRPC(createdUser),
	})
	return res, nil
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *connect.Request[holos.RegisterUserRequest]) (*connect.Response[holos.RegisterUserResponse], error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	var dbUser *ent.User
	var dbOrg *ent.Organization
	var rpcUser holos.User
	var rpcOrg org.Organization

	userMask, err := fieldmask_utils.MaskFromProtoFieldMask(req.Msg.GetUserMask(), strings.PascalCase)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}
	orgMask, err := fieldmask_utils.MaskFromProtoFieldMask(req.Msg.GetOrganizationMask(), strings.PascalCase)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	// Server assigns IDs.
	userID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeInternal, err))
	}
	orgID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeInternal, err))
	}
	bareID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeInternal, err))
	}
	refID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeInternal, err))
	}

	// Perform registration in a single transaction.
	if err := WithTx(ctx, h.db, func(tx *ent.Tx) (err error) {
		// Create the user
		dbUser, err = tx.User.Create().
			SetID(userID).
			SetEmail(authnID.Email()).
			SetIss(authnID.Issuer()).
			SetSub(authnID.Subject()).
			SetName(authnID.Name()).
			Save(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				rpcErr := connect.NewError(connect.CodeAlreadyExists, errors.New("user already registered"))
				errInfo := &errdetails.ErrorInfo{
					Reason:   "USER_EXISTS",
					Domain:   "user.holos.run",
					Metadata: map[string]string{"email": authnID.Email()},
				}
				if detail, detailErr := connect.NewErrorDetail(errInfo); detailErr != nil {
					logger.FromContext(ctx).ErrorContext(ctx, detailErr.Error(), "err", detailErr)
				} else {
					rpcErr.AddDetail(detail)
				}
				return errors.Wrap(rpcErr)
			}
			return errors.Wrap(err)
		}
		// Create the org
		dbOrg, err = tx.Organization.Create().
			SetID(orgID).
			SetName(cleanAndAppendRandom(authnID.Name())).
			SetDisplayName(authnID.GivenName() + "'s Org").
			SetCreatorID(userID).
			SetEditorID(userID).
			AddUserIDs(userID).
			Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}
		// Create the platforms.
		decoder := json.NewDecoder(bytes.NewReader([]byte(BareForm)))
		decoder.DisallowUnknownFields()
		var form storage.Form
		if err := decoder.Decode(&form); err != nil {
			return errors.Wrap(err)
		}

		decoder = json.NewDecoder(bytes.NewReader([]byte(Model)))
		decoder.DisallowUnknownFields()
		var model storage.Model
		if err := decoder.Decode(&model); err != nil {
			return errors.Wrap(err)
		}

		// Add the platforms.
		err = tx.Platform.Create().
			SetID(bareID).
			SetName("bare").
			SetDisplayName("Bare Platform").
			SetForm(&form).
			SetModel(&model).
			SetCreatorID(userID).
			SetEditorID(userID).
			SetOrgID(orgID).
			Exec(ctx)
		if err != nil {
			return errors.Wrap(err)
		}
		err = tx.Platform.Create().
			SetID(refID).
			SetName("reference").
			SetDisplayName("Holos Reference Platform").
			SetForm(&form).
			SetModel(&model).
			SetCreatorID(userID).
			SetEditorID(userID).
			SetOrgID(orgID).
			Exec(ctx)
		if err != nil {
			return errors.Wrap(err)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(connect.NewError(connect.CodeFailedPrecondition, err))
	}

	if err = fieldmask_utils.StructToStruct(userMask, UserToRPC(dbUser), &rpcUser); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}
	if err = fieldmask_utils.StructToStruct(orgMask, OrganizationToRPC(dbOrg), &rpcOrg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Wrap(err))
	}

	msg := holos.RegisterUserResponse{
		User:         &rpcUser,
		Organization: &rpcOrg,
	}

	return connect.NewResponse(&msg), nil
}

// UserToRPC returns an *holos.User adapted from *ent.User u.
func UserToRPC(entity *ent.User) *holos.User {
	uid := entity.ID.String()
	iamUser := holos.User{
		Id:    &uid,
		Email: entity.Email,
		Name:  &entity.Name,
		Subject: &object.Subject{
			Iss: entity.Iss,
			Sub: entity.Sub,
		},
		Detail: &object.Detail{
			CreatedBy: &object.ResourceEditor{
				Editor: &object.ResourceEditor_UserId{
					UserId: uid,
				},
			},
			CreatedAt: timestamppb.New(entity.CreatedAt),
			UpdatedBy: &object.ResourceEditor{
				Editor: &object.ResourceEditor_UserId{
					UserId: uid,
				},
			},
			UpdatedAt: timestamppb.New(entity.UpdatedAt),
		},
	}
	return &iamUser
}

func getUser(ctx context.Context, client *ent.Client, authnID authn.Identity) (*ent.User, error) {
	user, err := client.User.Query().
		Where(
			user.Iss(authnID.Issuer()),
			user.Sub(authnID.Subject()),
		).
		Only(ctx)
	if err != nil {
		if ent.MaskNotFound(err) == nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.Wrap(err))
		}
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	return user, nil
}

func (h *UserHandler) createUser(ctx context.Context, client *ent.Client, rpcUser *holos.User) (*ent.User, error) {
	log := logger.FromContext(ctx)
	var dbUser *ent.User
	dbUser, err := client.User.Create().
		SetName(rpcUser.GetName()).
		SetIss(rpcUser.GetSubject().GetIss()).
		SetSub(rpcUser.GetSubject().GetSub()).
		SetEmail(rpcUser.GetEmail()).
		Save(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Wrap(err))
	}
	log = log.With("user", dbUser)
	log.InfoContext(ctx, "created user")
	return dbUser, nil
}

func (h *UserHandler) createCallerUser(ctx context.Context) (*ent.User, error) {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	emailVerified := authnID.Verified()
	name := authnID.Name()
	givenName := authnID.GivenName()
	familyName := authnID.FamilyName()
	picture := authnID.Picture()

	rpcUser := holos.User{
		Subject: &object.Subject{
			Iss: authnID.Issuer(),
			Sub: authnID.Subject(),
		},
		Email:         authnID.Email(),
		EmailVerified: &emailVerified,
		Name:          &name,
		GivenName:     &givenName,
		FamilyName:    &familyName,
		Groups:        authnID.Groups(),
		Picture:       &picture,
	}

	var createdUser *ent.User
	err = WithTx(ctx, h.db, func(tx *ent.Tx) error {
		createdUser, err = h.createUser(ctx, tx.Client(), &rpcUser)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "could not save transaction", "err", err)
		return nil, err
	}

	return createdUser.Unwrap(), nil
}

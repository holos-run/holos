package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
)

const AdminEmail = "jeff@openinfrastructure.co"

// NewSystemHandler returns a new SystemService implementation.
func NewSystemHandler(db *ent.Client) *SystemHandler {
	return &SystemHandler{db: db}
}

// SystemHandler implements the PlatformService interface.
type SystemHandler struct {
	db *ent.Client
}

func (h *SystemHandler) checkAdmin(ctx context.Context) error {
	authnID, err := authn.FromContext(ctx)
	if err != nil {
		return connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	if authnID.Email() != AdminEmail {
		err := fmt.Errorf("not an admin:\n\thave (%+v)\n\twant (%+v)", authnID.Email(), AdminEmail)
		return connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}
	return nil
}

func (h *SystemHandler) DropTables(ctx context.Context, req *connect.Request[holos.EmptyRequest]) (*connect.Response[holos.EmptyResponse], error) {
	if err := h.checkAdmin(ctx); err != nil {
		return nil, err
	}
	log := logger.FromContext(ctx)

	if err := WithTx(ctx, h.db, func(tx *ent.Tx) (err error) {
		var n int
		if n, err = tx.Platform.Delete().Exec(ctx); err != nil {
			return errors.Wrap(err)
		}
		log.WarnContext(ctx, "deleted platforms", "count", n)
		if n, err = tx.Organization.Delete().Exec(ctx); err != nil {
			return errors.Wrap(err)
		}
		log.WarnContext(ctx, "deleted organizations", "count", n)
		if n, err = tx.User.Delete().Exec(ctx); err != nil {
			return errors.Wrap(err)
		}
		log.WarnContext(ctx, "deleted users", "count", n)
		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	return connect.NewResponse(&holos.EmptyResponse{}), nil
}

func (h *SystemHandler) SeedDatabase(ctx context.Context, req *connect.Request[holos.EmptyRequest]) (*connect.Response[holos.EmptyResponse], error) {
	if err := h.checkAdmin(ctx); err != nil {
		return nil, err
	}

	if err := WithTx(ctx, h.db, func(tx *ent.Tx) (err error) {
		jeff, err := tx.User.Create().
			SetEmail("jeff@openinfrastructure.co").
			SetIss("https://login.ois.run").
			SetSub("261773693724656988").
			SetName("Jeff McCune").
			Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}
		nate, err := tx.User.Create().
			SetEmail("nate@openinfrastructure.co").
			SetIss("https://login.ois.run").
			SetSub("261775487611699776").
			SetName("Nate McCurdy").
			Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}
		gary, err := tx.User.Create().
			SetEmail("gary@openinfrastructure.co").
			SetIss("https://login.ois.run").
			SetSub("261775531836441152").
			SetName("Gary Larizza").
			Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}

		// Create the org
		org, err := tx.Organization.Create().
			SetName("ois").
			SetDisplayName("Open Infrastructure Services").
			SetCreator(jeff).
			Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}

		// Add org memebers
		org, err = org.Update().AddUsers(jeff, gary, nate).Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}

		var hf holos.PlatformForm
		if err := json.Unmarshal([]byte(BareForm), &hf); err != nil {
			return errors.Wrap(err)
		}

		// Add a platform
		err = tx.Platform.Create().
			SetName("bare").
			SetDisplayName("Bare Platform").
			SetConfigForm(&hf).
			SetCreator(jeff).
			SetOrgID(org.ID).
			Exec(ctx)
		if err != nil {
			return errors.Wrap(err)
		}

		stuff := []string{"Jeff", "Gary", "Nate"}
		for _, name := range stuff {
			err := tx.Platform.Create().
				SetName(strings.ToLower(name)).
				SetDisplayName(name + "'s Platform").
				SetConfigForm(&hf).
				SetCreator(jeff).
				SetOrgID(org.ID).
				Exec(ctx)
			if err != nil {
				return errors.Wrap(err)
			}
		}

		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.Wrap(err))
	}

	return connect.NewResponse(&holos.EmptyResponse{}), nil
}

const BareForm = `{
  "kind": "PlatformForm",
  "spec": {
    "sections": [
      {
        "name": "org",
        "description": "Organization config values are used to derive more specific configuration values throughout the platform.",
        "displayName": "Organization",
        "fieldConfigs": [
          {
            "key": "name",
            "type": "input",
            "props": {
              "label": "Name",
              "required": true,
              "description": "DNS label, e.g. 'example'",
              "placeholder": "example"
            }
          },
          {
            "key": "domain",
            "type": "input",
            "props": {
              "label": "Domain",
              "required": true,
              "description": "DNS domain, e.g. 'example.com'",
              "placeholder": "example.com"
            }
          },
          {
            "key": "displayName",
            "type": "input",
            "props": {
              "label": "Display Name",
              "required": true,
              "description": "Display name, e.g. 'Example Organization'",
              "placeholder": "Example Organization"
            }
          },
          {
            "key": "contactEmail",
            "type": "input",
            "props": {
              "label": "Contact Email",
              "required": true,
              "description": "Technical contact email address",
              "placeholder": "platform-team@example.com"
            }
          }
        ]
      }
    ]
  },
  "metadata": {
    "name": "bare"
  },
  "apiVersion": "forms.holos.run/v1alpha1"
}`

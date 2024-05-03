package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	psvc "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
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
			SetID(uuid.FromStringOrNil("018f36fb-e3f2-7f7f-a72f-ce48eb16c82d")).
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
			SetID(uuid.FromStringOrNil("018f36fb-e3f7-7f7f-a1c5-c85fb735d215")).
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

		var form psvc.Form
		if err := json.Unmarshal([]byte(BareForm), &form); err != nil {
			return errors.Wrap(err)
		}

		var model psvc.Model
		if err := json.Unmarshal([]byte(Model), &model); err != nil {
			return errors.Wrap(err)
		}

		// Add a platform
		err = tx.Platform.Create().
			SetID(uuid.FromStringOrNil("018f36fb-e3ff-7f7f-a5d1-7ca2bf499e94")).
			SetName("bare").
			SetDisplayName("Bare Platform").
			SetForm(&form).
			SetModel(&model).
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
				SetForm(&form).
				SetModel(&model).
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

const Model = `{"model":{}}`

const BareForm = `
{
  "fields": [
    {
      "key": "org",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Organization",
        "description": "Organization config values are used to derive more specific configuration values throughout the platform."
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "name",
          "type": "input",
          "props": {
            "label": "Name",
            "description": "DNS label, e.g. 'example'",
            "pattern": "[a-z][0-9a-z]{4,29}",
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited."
            }
          },
          "resetOnHide": true
        },
        {
          "key": "domain",
          "type": "input",
          "props": {
            "label": "Domain",
            "placeholder": "example.com",
            "description": "DNS domain, e.g. 'example.com'"
          },
          "resetOnHide": true
        },
        {
          "key": "displayName",
          "type": "input",
          "props": {
            "label": "Display Name",
            "placeholder": "Example Organization",
            "description": "Display name, e.g. 'Example Organization'"
          },
          "resetOnHide": true
        },
        {
          "key": "contactEmail",
          "type": "input",
          "props": {
            "label": "Contact Email",
            "placeholder": "platform-team@example.com",
            "description": "Technical contact email address"
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "privacy",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Data Privacy",
        "description": "Configure data privacy aspects of the platform."
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "country",
          "type": "select",
          "props": {
            "label": "Select Planet",
            "description": "Juridiction of applicable data privacy laws.",
            "options": [
              {
                "value": "mercury",
                "label": "Mercury"
              },
              {
                "value": "venus",
                "label": "Venus"
              },
              {
                "value": "earth",
                "label": "Earth"
              },
              {
                "value": "mars",
                "label": "Mars"
              },
              {
                "value": "jupiter",
                "label": "Jupiter"
              },
              {
                "value": "saturn",
                "label": "Saturn"
              },
              {
                "value": "uranus",
                "label": "Uranus"
              },
              {
                "value": "neptune",
                "label": "Neptune"
              }
            ]
          },
          "resetOnHide": true
        },
        {
          "key": "regions",
          "type": "select",
          "props": {
            "label": "Select Regions",
            "description": "Select the regions this platform operates in.",
            "multiple": true,
            "selectAllOption": "Select All",
            "options": [
              {
                "value": "us-east-2",
                "label": "Ohio"
              },
              {
                "value": "us-west-2",
                "label": "Oregon"
              },
              {
                "value": "eu-west-1",
                "label": "Ireland"
              },
              {
                "value": "eu-west-2",
                "label": "London",
                "disabled": true
              }
            ]
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "terms",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Terms and Conditions",
        "description": "Example of a boolean checkbox."
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "didAgree",
          "type": "checkbox",
          "props": {
            "label": "Accept terms",
            "description": "In order to proceed, please accept terms",
            "pattern": "true",
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "Please accept the terms"
            }
          },
          "resetOnHide": true
        }
      ]
    }
  ]
}
`

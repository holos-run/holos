package handler

import (
	"bytes"
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
	holosstrings "github.com/holos-run/holos/internal/strings"
	storage "github.com/holos-run/holos/service/gen/holos/storage/v1alpha1"
	system "github.com/holos-run/holos/service/gen/holos/system/v1alpha1"
	"github.com/holos-run/holos/version"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
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

func (h *SystemHandler) GetVersion(ctx context.Context, req *connect.Request[system.GetVersionRequest]) (*connect.Response[system.GetVersionResponse], error) {
	_, err := authn.FromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.Wrap(err))
	}

	mask, err := fieldmask_utils.MaskFromProtoFieldMask(req.Msg.GetFieldMask(), holosstrings.PascalCase)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	srcVersion := version.NewVersionInfo()
	var rpcVersion system.Version
	if err := fieldmask_utils.StructToStruct(mask, &srcVersion, &rpcVersion); err != nil {
		return nil, errors.Wrap(err)
	}

	return connect.NewResponse(&system.GetVersionResponse{Version: &rpcVersion}), nil
}

func (h *SystemHandler) DropTables(ctx context.Context, req *connect.Request[system.DropTablesRequest]) (*connect.Response[system.DropTablesResponse], error) {
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

	return connect.NewResponse(&system.DropTablesResponse{}), nil
}

func (h *SystemHandler) SeedDatabase(ctx context.Context, req *connect.Request[system.SeedDatabaseRequest]) (*connect.Response[system.SeedDatabaseResponse], error) {
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
			SetEditor(jeff).
			Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}

		// Add org memebers
		org, err = org.Update().AddUsers(jeff, gary, nate).Save(ctx)
		if err != nil {
			return errors.Wrap(err)
		}

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

		// Add a platform
		err = tx.Platform.Create().
			SetID(uuid.FromStringOrNil("018f36fb-e3ff-7f7f-a5d1-7ca2bf499e94")).
			SetName("bare").
			SetDisplayName("Bare Platform").
			SetForm(&form).
			SetModel(&model).
			SetCreator(jeff).
			SetEditor(jeff).
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
				SetEditor(jeff).
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

	return connect.NewResponse(&system.SeedDatabaseResponse{}), nil
}

const Model = `{"model":{}}`

const BareForm = `
{
  "field_configs": [
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
            "pattern": "^[a-z]([0-9a-z]|-){1,28}[0-9a-z]$",
            "minLength": 3,
            "maxLength": 30,
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "It must be 3 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited."
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
            "minLength": 3,
            "maxLength": 100,
            "description": "DNS domain, e.g. 'example.com'",
            "required": true
          },
          "resetOnHide": true
        },
        {
          "key": "displayName",
          "type": "input",
          "props": {
            "label": "Display Name",
            "placeholder": "Example Organization",
            "description": "Display name, e.g. 'Example Organization'",
            "maxLength": 100,
            "required": true
          },
          "resetOnHide": true
        },
        {
          "key": "contactEmail",
          "type": "input",
          "props": {
            "label": "Contact Email",
            "placeholder": "platform-team@example.com",
            "description": "Technical contact email address",
            "required": true
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "cloud",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Cloud Providers",
        "description": "Select the services that provide resources for the platform."
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "providers",
          "type": "select",
          "props": {
            "label": "Select Providers",
            "description": "Select the cloud providers the platform builds upon.",
            "multiple": true,
            "selectAllOption": "Select All",
            "options": [
              {
                "value": "aws",
                "label": "Amazon Web Services"
              },
              {
                "value": "gcp",
                "label": "Google Cloud Platform"
              },
              {
                "value": "azure",
                "label": "Microsoft Azure"
              },
              {
                "value": "cloudflare",
                "label": "Cloudflare"
              },
              {
                "value": "github",
                "label": "GitHub"
              },
              {
                "value": "ois",
                "label": "Open Infrastructure Services"
              },
              {
                "value": "onprem",
                "label": "On Premises",
                "disabled": true
              }
            ]
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "aws",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Amazon Web Services",
        "description": "Provide the information necessary for Holos to manage AWS resources to provide the platform."
      },
      "expressions": {
        "hide": "!formState.model.cloud?.providers?.includes(\"aws\")"
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "primaryRoleARN",
          "type": "input",
          "props": {
            "label": "Holos Admin Role ARN",
            "description": "Enter the AWS Role ARN Holos will use to bootstrap resources.  For example, arn:aws:iam::123456789012:role/HolosAdminAccess",
            "pattern": "^arn:.*",
            "minLength": 4,
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "Must be a valid ARN.  Refer to https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html"
            }
          },
          "resetOnHide": true
        },
        {
          "key": "regions",
          "type": "select",
          "props": {
            "label": "Select Regions",
            "description": "Select the AWS regions this platform operates in.",
            "multiple": true,
            "required": true,
            "selectAllOption": "Select All",
            "options": [
              {
                "value": "us-east-1",
                "label": "N. Virginia (us-east-1)"
              },
              {
                "value": "us-east-2",
                "label": "Ohio (us-east-2)"
              },
              {
                "value": "us-west-1",
                "label": "N. California (us-west-1)"
              },
              {
                "value": "us-west-2",
                "label": "Oregon (us-west-2)"
              },
              {
                "value": "us-gov-west1",
                "label": "US GovCloud West (us-gov-west1)"
              },
              {
                "value": "us-gov-east1",
                "label": "US GovCloud East (us-gov-east1)"
              },
              {
                "value": "ca-central-1",
                "label": "Canada (ca-central-1)"
              },
              {
                "value": "eu-north-1",
                "label": "Stockholm (eu-north-1)"
              },
              {
                "value": "eu-west-1",
                "label": "Ireland (eu-west-1)"
              },
              {
                "value": "eu-west-2",
                "label": "London (eu-west-2)"
              },
              {
                "value": "eu-west-3",
                "label": "Paris (eu-west-3)"
              },
              {
                "value": "eu-central-1",
                "label": "Frankfurt (eu-central-1)"
              },
              {
                "value": "eu-south-1",
                "label": "Milan (eu-south-1)"
              },
              {
                "value": "af-south-1",
                "label": "Cape Town (af-south-1)"
              },
              {
                "value": "ap-northeast-1",
                "label": "Tokyo (ap-northeast-1)"
              },
              {
                "value": "ap-northeast-2",
                "label": "Seoul (ap-northeast-2)"
              },
              {
                "value": "ap-northeast-3",
                "label": "Osaka (ap-northeast-3)"
              },
              {
                "value": "ap-southeast-1",
                "label": "Singapore (ap-southeast-1)"
              },
              {
                "value": "ap-southeast-2",
                "label": "Sydney (ap-southeast-2)"
              },
              {
                "value": "ap-east-1",
                "label": "Hong Kong (ap-east-1)"
              },
              {
                "value": "ap-south-1",
                "label": "Mumbai (ap-south-1)"
              },
              {
                "value": "me-south-1",
                "label": "Bahrain (me-south-1)"
              },
              {
                "value": "sa-east-1",
                "label": "São Paulo (sa-east-1)"
              },
              {
                "value": "cn-north-1",
                "label": "Bejing (cn-north-1)"
              },
              {
                "value": "cn-northwest-1",
                "label": "Ningxia (cn-northwest-1)"
              },
              {
                "value": "ap-southeast-3",
                "label": "Jakarta (ap-southeast-3)"
              }
            ]
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "gcp",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Google Cloud Platform",
        "description": "Use this form to configure platform level GCP settings."
      },
      "expressions": {
        "hide": "!formState.model.cloud?.providers?.includes(\"gcp\")"
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "regions",
          "type": "select",
          "props": {
            "label": "Select Regions",
            "description": "Select the GCP regions this platform operates in.",
            "multiple": true,
            "selectAllOption": "Select All",
            "options": [
              {
                "value": "africa-south1",
                "label": "africa-south1"
              },
              {
                "value": "asia-east1",
                "label": "asia-east1"
              },
              {
                "value": "asia-east2",
                "label": "asia-east2"
              },
              {
                "value": "asia-northeast1",
                "label": "asia-northeast1"
              },
              {
                "value": "asia-northeast2",
                "label": "asia-northeast2"
              },
              {
                "value": "asia-northeast3",
                "label": "asia-northeast3"
              },
              {
                "value": "asia-south1",
                "label": "asia-south1"
              },
              {
                "value": "asia-south2",
                "label": "asia-south2"
              },
              {
                "value": "asia-southeast1",
                "label": "asia-southeast1"
              },
              {
                "value": "asia-southeast2",
                "label": "asia-southeast2"
              },
              {
                "value": "australia-southeast1",
                "label": "australia-southeast1"
              },
              {
                "value": "australia-southeast2",
                "label": "australia-southeast2"
              },
              {
                "value": "europe-central2",
                "label": "europe-central2"
              },
              {
                "value": "europe-north1",
                "label": "europe-north1"
              },
              {
                "value": "europe-southwest1",
                "label": "europe-southwest1"
              },
              {
                "value": "europe-west1",
                "label": "europe-west1"
              },
              {
                "value": "europe-west10",
                "label": "europe-west10"
              },
              {
                "value": "europe-west12",
                "label": "europe-west12"
              },
              {
                "value": "europe-west2",
                "label": "europe-west2"
              },
              {
                "value": "europe-west3",
                "label": "europe-west3"
              },
              {
                "value": "europe-west4",
                "label": "europe-west4"
              },
              {
                "value": "europe-west6",
                "label": "europe-west6"
              },
              {
                "value": "europe-west8",
                "label": "europe-west8"
              },
              {
                "value": "europe-west9",
                "label": "europe-west9"
              },
              {
                "value": "me-central1",
                "label": "me-central1"
              },
              {
                "value": "me-central2",
                "label": "me-central2"
              },
              {
                "value": "me-west1",
                "label": "me-west1"
              },
              {
                "value": "northamerica-northeast1",
                "label": "northamerica-northeast1"
              },
              {
                "value": "northamerica-northeast2",
                "label": "northamerica-northeast2"
              },
              {
                "value": "southamerica-east1",
                "label": "southamerica-east1"
              },
              {
                "value": "southamerica-west1",
                "label": "southamerica-west1"
              },
              {
                "value": "us-central1",
                "label": "us-central1"
              },
              {
                "value": "us-east1",
                "label": "us-east1"
              },
              {
                "value": "us-east4",
                "label": "us-east4"
              },
              {
                "value": "us-east5",
                "label": "us-east5"
              },
              {
                "value": "us-south1",
                "label": "us-south1"
              },
              {
                "value": "us-west1",
                "label": "us-west1"
              },
              {
                "value": "us-west2",
                "label": "us-west2"
              },
              {
                "value": "us-west3",
                "label": "us-west3"
              },
              {
                "value": "us-west4",
                "label": "us-west4"
              }
            ]
          },
          "resetOnHide": true
        },
        {
          "key": "gcpProjectID",
          "type": "input",
          "props": {
            "label": "Project ID",
            "description": "Enter the project id where the provisioner cluster resides.",
            "pattern": "^[a-z]([0-9a-z]|-){1,28}[0-9a-z]$",
            "minLength": 6,
            "maxLength": 30,
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "It must be 3 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited."
            }
          },
          "resetOnHide": true
        },
        {
          "key": "gcpProjectNumber",
          "type": "input",
          "props": {
            "label": "Project Number",
            "type": "number",
            "description": "Enter the project number where the provisioner cluster resides.",
            "pattern": "^[0-9]+$",
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "Must be a valid project number."
            }
          },
          "resetOnHide": true
        },
        {
          "key": "provisionerCABundle",
          "type": "input",
          "props": {
            "label": "Provisioner CA Bundle",
            "description": "Enter the provisioner cluster ca bundle.  kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.certificate-authority-data}'",
            "pattern": "^[0-9a-zA-Z]+=*$",
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "Must be a base64 encoded pem encoded certificate bundle."
            }
          },
          "resetOnHide": true
        },
        {
          "key": "provisionerURL",
          "type": "input",
          "props": {
            "label": "Provisioner URL",
            "description": "Enter the URL of the provisioner cluster API endpoint.  kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.server}'",
            "pattern": "^https://.*$",
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "Must be a https:// URL."
            }
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "cloudflare",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Cloudflare",
        "description": "Cloudflare is primarily used for DNS automation."
      },
      "expressions": {
        "hide": "!formState.model.cloud?.providers?.includes(\"cloudflare\")"
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "email",
          "type": "input",
          "props": {
            "label": "Account Email",
            "description": "Enter the Cloudflare email address to manage DNS",
            "minLength": 3,
            "required": true
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "github",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "GitHub",
        "description": "GitHub is primarily used to host Git repositories and execute Actions workflows."
      },
      "expressions": {
        "hide": "!formState.model.cloud?.providers?.includes(\"github\")"
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "primaryOrg",
          "type": "input",
          "props": {
            "label": "Organization",
            "description": "Enter the primary GitHub organization associed with the platform.",
            "pattern": "^(?!-)(?!.*--)([a-zA-Z0-9]|-){1,39}$",
            "minLength": 1,
            "maxLength": 39,
            "required": true
          },
          "validation": {
            "messages": {
              "pattern": "All characters must be either a hyphen or alphanumeric.  Cannot start with a hyphen.  Cannot include consecutive hyphens."
            }
          },
          "resetOnHide": true
        }
      ]
    },
    {
      "key": "backups",
      "wrappers": [
        "holos-panel"
      ],
      "props": {
        "label": "Backups",
        "description": "Configure platform level data backup settings.  Requires AWS."
      },
      "resetOnHide": true,
      "fieldGroup": [
        {
          "key": "s3bucket",
          "type": "select",
          "props": {
            "label": "S3 Bucket Region",
            "description": "Select the S3 Bucket Region.",
            "multiple": true,
            "options": [
              {
                "value": "us-east-1",
                "label": "N. Virginia (us-east-1)"
              },
              {
                "value": "us-east-2",
                "label": "Ohio (us-east-2)"
              },
              {
                "value": "us-west-1",
                "label": "N. California (us-west-1)"
              },
              {
                "value": "us-west-2",
                "label": "Oregon (us-west-2)"
              },
              {
                "value": "us-gov-west1",
                "label": "US GovCloud West (us-gov-west1)"
              },
              {
                "value": "us-gov-east1",
                "label": "US GovCloud East (us-gov-east1)"
              },
              {
                "value": "ca-central-1",
                "label": "Canada (ca-central-1)"
              },
              {
                "value": "eu-north-1",
                "label": "Stockholm (eu-north-1)"
              },
              {
                "value": "eu-west-1",
                "label": "Ireland (eu-west-1)"
              },
              {
                "value": "eu-west-2",
                "label": "London (eu-west-2)"
              },
              {
                "value": "eu-west-3",
                "label": "Paris (eu-west-3)"
              },
              {
                "value": "eu-central-1",
                "label": "Frankfurt (eu-central-1)"
              },
              {
                "value": "eu-south-1",
                "label": "Milan (eu-south-1)"
              },
              {
                "value": "af-south-1",
                "label": "Cape Town (af-south-1)"
              },
              {
                "value": "ap-northeast-1",
                "label": "Tokyo (ap-northeast-1)"
              },
              {
                "value": "ap-northeast-2",
                "label": "Seoul (ap-northeast-2)"
              },
              {
                "value": "ap-northeast-3",
                "label": "Osaka (ap-northeast-3)"
              },
              {
                "value": "ap-southeast-1",
                "label": "Singapore (ap-southeast-1)"
              },
              {
                "value": "ap-southeast-2",
                "label": "Sydney (ap-southeast-2)"
              },
              {
                "value": "ap-east-1",
                "label": "Hong Kong (ap-east-1)"
              },
              {
                "value": "ap-south-1",
                "label": "Mumbai (ap-south-1)"
              },
              {
                "value": "me-south-1",
                "label": "Bahrain (me-south-1)"
              },
              {
                "value": "sa-east-1",
                "label": "São Paulo (sa-east-1)"
              },
              {
                "value": "cn-north-1",
                "label": "Bejing (cn-north-1)"
              },
              {
                "value": "cn-northwest-1",
                "label": "Ningxia (cn-northwest-1)"
              },
              {
                "value": "ap-southeast-3",
                "label": "Jakarta (ap-southeast-3)"
              }
            ]
          },
          "expressions": {
            "props.disabled": "!formState.model.cloud?.providers?.includes(\"aws\")",
            "props.required": "formState.model.cloud?.providers?.includes(\"aws\")",
            "props.description": "formState.model.cloud?.providers?.includes(\"aws\") ? 'Select the S3 Bucket Region.' : 'Enable AWS in the Cloud Provider section to configure backups.'"
          },
          "resetOnHide": true
        }
      ]
    }
  ]
}
`

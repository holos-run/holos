// Package client provides configuration and convenience methods for making API calls to the holos server.
package client

import (
	"context"
	"flag"
	"time"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/token"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/organization/v1alpha1/organizationconnect"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/platform/v1alpha1/platformconnect"
	"github.com/holos-run/holos/service/gen/holos/user/v1alpha1/userconnect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
)

type PlatformMutation struct {
	Name        string
	DisplayName string
	flagSet     *flag.FlagSet
}

func (pm *PlatformMutation) FlagSet() *flag.FlagSet {
	if pm == nil {
		return nil
	}
	if pm.flagSet != nil {
		return pm.flagSet
	}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&pm.Name, "name", "example", "platform name")
	fs.StringVar(&pm.DisplayName, "display-name", "Example Platform", "platform display name")
	pm.flagSet = fs
	return fs
}

func New(cfg *Config) *Client {
	t := token.NewClient(cfg.Token())
	s := cfg.Client().Server()
	return &Client{
		cfg:    cfg,
		usrSvc: userconnect.NewUserServiceClient(t, s),
		orgSvc: organizationconnect.NewOrganizationServiceClient(t, s),
		pltSvc: platformconnect.NewPlatformServiceClient(t, s),
	}
}

// Client provides convenience methods for making API calls to the holos server.
type Client struct {
	cfg    *Config
	usrSvc userconnect.UserServiceClient
	pltSvc platformconnect.PlatformServiceClient
	orgSvc organizationconnect.OrganizationServiceClient
}

func (c *Client) Platforms(ctx context.Context, orgID string) ([]*platform.Platform, error) {
	if c == nil {
		return nil, errors.New("no service client")
	}
	req := &platform.ListPlatformsRequest{
		OrgId: orgID,
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"id", "name", "displayName"},
		},
	}
	resp, err := c.pltSvc.ListPlatforms(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}
	return resp.Msg.GetPlatforms(), nil
}

func (c *Client) UpdateForm(ctx context.Context, platformID string, form *object.Form) error {
	start := time.Now()
	req := &platform.UpdatePlatformRequest{
		PlatformId: platformID,
		Update:     &platform.PlatformMutation{Form: form},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"form"}},
	}
	_, err := c.pltSvc.UpdatePlatform(ctx, connect.NewRequest(req))
	if err != nil {
		return errors.Wrap(err)
	}
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "updated platform", "platform_id", platformID, "duration", time.Since(start))
	return nil
}

func (c *Client) UpdatePlatformModel(ctx context.Context, platformID string, model *structpb.Struct) error {
	start := time.Now()
	req := &platform.UpdatePlatformRequest{
		PlatformId: platformID,
		Update:     &platform.PlatformMutation{Model: model},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"model"}},
	}
	_, err := c.pltSvc.UpdatePlatform(ctx, connect.NewRequest(req))
	if err != nil {
		return errors.Wrap(err)
	}
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "updated platform", "platform_id", platformID, "duration", time.Since(start))
	return nil
}

// PlatformModel gets the platform model from the PlatformService.
func (c *Client) PlatformModel(ctx context.Context, platformID string) (*structpb.Struct, error) {
	start := time.Now()
	req := &platform.GetPlatformRequest{
		PlatformId: platformID,
		FieldMask:  &fieldmaskpb.FieldMask{Paths: []string{"spec.model"}},
	}
	pf, err := c.pltSvc.GetPlatform(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "get platform", "platform_id", platformID, "duration", time.Since(start))
	return pf.Msg.GetPlatform().GetSpec().GetModel(), nil
}

func (c *Client) CreatePlatform(ctx context.Context, pm PlatformMutation) (*platform.Platform, error) {
	log := logger.FromContext(ctx).With("platform", pm.Name)
	start := time.Now()
	req := &platform.CreatePlatformRequest{
		OrgId: c.cfg.context.OrgID,
		Create: &platform.PlatformMutation{
			Name:        &pm.Name,
			DisplayName: &pm.DisplayName,
		},
	}
	pf, err := c.pltSvc.CreatePlatform(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	log = log.With("platform_id", pf.Msg.GetPlatform().GetId())
	log.DebugContext(ctx, "create platform", "duration", time.Since(start))
	return pf.Msg.GetPlatform(), nil
}

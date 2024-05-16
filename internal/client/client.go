// Package client provides configuration and convenience methods for making API calls to the holos server.
package client

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/token"
	"github.com/holos-run/holos/service/gen/holos/organization/v1alpha1/organizationconnect"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/platform/v1alpha1/platformconnect"
	"github.com/holos-run/holos/service/gen/holos/user/v1alpha1/userconnect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

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

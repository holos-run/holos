package register

import (
	"context"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/token"
	org "github.com/holos-run/holos/service/gen/holos/organization/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/organization/v1alpha1/organizationconnect"
	user "github.com/holos-run/holos/service/gen/holos/user/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/user/v1alpha1/userconnect"
)

// User registers the user with the holos server.
func User(ctx context.Context, cfg *client.Config) error {
	log := logger.FromContext(ctx)
	client := userconnect.NewUserServiceClient(token.NewClient(cfg.Token()), cfg.Client().Server())

	var err error
	var u *user.User
	var o *org.Organization
	cc := &holos.ClientContext{}

	u, err = getUser(ctx, client)
	if err != nil {
		if connect.CodeOf(err) != connect.CodeNotFound {
			return errors.Wrap(err)
		}
		if u, o, err = registerUser(ctx, client); err != nil {
			return errors.Wrap(err)
		}
		// Save the registration context
		cc.OrgID = o.GetOrgId()
		cc.UserID = u.GetId()
		if err := cc.Save(ctx); err != nil {
			return errors.Wrap(err)
		}
		log.InfoContext(ctx, "created user", "email", u.GetEmail(), "id", u.GetId())
	}

	if cc.Exists() {
		if err := cc.Load(ctx); err != nil {
			return errors.Wrap(err)
		}
	}

	// Ensure the current user id gets saved.
	cc.UserID = u.GetId()

	// Ensure an org ID gets saved.
	if cc.OrgID == "" {
		org, err := getOrg(ctx, cfg)
		if err != nil {
			return errors.Wrap(err)
		}
		cc.OrgID = org.GetOrgId()
	}

	// One last save, we know we have the user id and org id at this point.
	if err := cc.Save(ctx); err != nil {
		return errors.Wrap(err)
	}

	log.InfoContext(ctx, "user", "email", u.GetEmail(), "user_id", cc.UserID, "org_id", cc.OrgID)
	return nil
}

func getUser(ctx context.Context, client userconnect.UserServiceClient) (*user.User, error) {
	req := connect.NewRequest(&user.GetUserRequest{})
	resp, err := client.GetUser(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return resp.Msg.GetUser(), nil
}

// getOrg returns the first organization returned from the ListOrganizations rpc
// method.
func getOrg(ctx context.Context, cfg *client.Config) (*org.Organization, error) {
	client := organizationconnect.NewOrganizationServiceClient(token.NewClient(cfg.Token()), cfg.Client().Server())
	req := connect.NewRequest(&org.ListOrganizationsRequest{})
	resp, err := client.ListOrganizations(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	orgs := resp.Msg.GetOrganizations()
	if len(orgs) == 0 {
		return nil, nil
	} else {
		return orgs[0], nil
	}

}

func registerUser(ctx context.Context, client userconnect.UserServiceClient) (*user.User, *org.Organization, error) {
	req := connect.NewRequest(&user.RegisterUserRequest{})
	resp, err := client.RegisterUser(ctx, req)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}
	return resp.Msg.GetUser(), resp.Msg.GetOrganization(), nil
}

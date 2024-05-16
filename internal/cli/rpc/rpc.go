package rpc

import (
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/token"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/platform/v1alpha1/platformconnect"
	"github.com/spf13/cobra"
)

// Remove these
const BarePlatformID = "018f36fb-e3ff-7f7f-a5d1-7ca2bf499e94"

type Config struct {
	holos  *holos.Config
	client *holos.ClientConfig
	token  *token.Config
}

// New returns a new rpc command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("rpc")
	cmd.Short = "remote procedure calls"

	config := &Config{
		holos:  cfg,
		client: holos.NewClientConfig(),
		token:  token.NewConfig(),
	}

	cmd.PersistentFlags().AddGoFlagSet(config.client.FlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.token.FlagSet())
	cmd.AddCommand(NewPlatformModel(config))
	return cmd
}

func NewPlatformModel(cfg *Config) *cobra.Command {
	cmd := command.New("platform-model")
	cmd.Short = "get the platform model"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		log := logger.FromContext(ctx)
		// client := platformconnect.NewPlatformServiceClient(token.NewClient(cfg.token), cfg.client.Server())
		client := platformconnect.NewPlatformServiceClient(token.NewClient(cfg.token), cfg.client.Server())
		// JEFFTODO - FieldMask
		log.WarnContext(ctx, "JEFFTODO use the fieldmask to get only the model")
		req := connect.NewRequest(&platform.GetPlatformRequest{PlatformId: BarePlatformID})
		res, err := client.GetPlatform(ctx, req)
		if err != nil {
			return errors.Wrap(fmt.Errorf("could not get platform: %w", err))
		}
		if res == nil || res.Msg == nil || res.Msg.Platform == nil || res.Msg.Platform.Spec == nil {
			return errors.Wrap(fmt.Errorf("response missing platform spec:\n\thave: (%+v)\n\twant: (response message platform.spec)", res))
		}
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(res.Msg.Platform.Spec.Model); err != nil {
			return errors.Wrap(fmt.Errorf("could not encode json: %w", err))
		}
		return err
	}
	return cmd
}

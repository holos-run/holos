package rpc

import (
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/token"
	h "github.com/holos-run/holos/service/gen/holos/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/v1alpha1/holosconnect"
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
		ctx := cmd.Context()
		client := holosconnect.NewPlatformServiceClient(token.NewClient(cfg.token), cfg.client.Server())
		req := connect.NewRequest(&h.PlatformServiceGetModelRequest{PlatformId: BarePlatformID})
		res, err := client.GetModel(ctx, req)
		if err != nil {
			return errors.Wrap(fmt.Errorf("could not get model: %w", err))
		}
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(res.Msg.Model); err != nil {
			return errors.Wrap(fmt.Errorf("could not encode json: %w", err))
		}
		return err
	}
	return cmd
}

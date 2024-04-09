package kv

import (
	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/errors"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// New returns the kv root command for the cli
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("kv")
	cmd.Short = "work with secrets in the provisioner cluster"
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}
	// flags
	cmd.PersistentFlags().SortFlags = false
	cmd.PersistentFlags().AddGoFlagSet(cfg.KVFlagSet())
	// subcommands
	cmd.AddCommand(newGetCmd(cfg))
	cmd.AddCommand(newListCmd(cfg))
	cmd.AddCommand(newPutCmd(cfg))
	return cmd
}

func newClientSet(cfg *holos.Config) (*kubernetes.Clientset, error) {
	kcfg, err := clientcmd.BuildConfigFromFlags("", cfg.KVKubeconfig())
	if err != nil {
		return nil, errors.Wrap(err)
	}
	clientset, err := kubernetes.NewForConfig(kcfg)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return clientset, nil
}

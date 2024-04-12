package kv

import (
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/secret"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListCmd(cfg *holos.Config) *cobra.Command {
	cmd := command.New("list")
	cmd.Args = cobra.NoArgs
	cmd.Short = "list secrets"
	cmd.Flags().SortFlags = false
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())
	cmd.RunE = makeListRunFunc(cfg)

	return cmd
}

func makeListRunFunc(cfg *holos.Config) command.RunFunc {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		cs, err := newClientSet(cfg)
		if err != nil {
			return err
		}
		selector := metav1.ListOptions{LabelSelector: secret.NameLabel}
		secrets, err := cs.CoreV1().Secrets(cfg.KVNamespace()).List(ctx, selector)
		if err != nil {
			return errors.Wrap(err)
		}
		labels := make(map[string]bool)
		for _, s := range secrets.Items {
			if value, ok := s.Labels[secret.NameLabel]; ok {
				labels[value] = true
			}
		}
		for label := range labels {
			cfg.Println(label)
		}
		return nil
	}
}

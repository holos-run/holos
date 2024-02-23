package kv

import (
	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListCmd(cfg *config.Config) *cobra.Command {
	cmd := command.New("list")
	cmd.Args = cobra.NoArgs
	cmd.Short = "list secrets"
	cmd.Flags().SortFlags = false
	cmd.RunE = makeListRunFunc(cfg)

	return cmd
}

func makeListRunFunc(cfg *config.Config) command.RunFunc {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		cs, err := newClientSet(cfg)
		if err != nil {
			return err
		}
		selector := metav1.ListOptions{LabelSelector: NameLabel}
		secrets, err := cs.CoreV1().Secrets(cfg.KVNamespace()).List(ctx, selector)
		if err != nil {
			return wrapper.Wrap(err)
		}
		labels := make(map[string]bool)
		for _, secret := range secrets.Items {
			if value, ok := secret.Labels[NameLabel]; ok {
				labels[value] = true
			}
		}
		for label := range labels {
			cfg.Println(label)
		}
		return nil
	}
}

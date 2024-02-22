package cli

import (
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sort"
)

const NameLabel = "holos.run/secret.name"

// newKVRootCmd returns the kv root command for the cli
func newKVRootCmd(cfg *config.Config) *cobra.Command {
	cmd := newCmd("kv")
	cmd.Short = "work with secrets in the provisioner cluster"
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}
	// flags
	cmd.PersistentFlags().SortFlags = false
	cmd.PersistentFlags().AddGoFlagSet(cfg.KVFlagSet())
	// subcommands
	cmd.AddCommand(newKVGetCmd(cfg))
	return cmd
}

func newKVGetCmd(cfg *config.Config) *cobra.Command {
	cmd := newCmd("get")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "print secret data in txtar format"
	cmd.Flags().SortFlags = false
	cmd.RunE = makeKVGetRunFunc(cfg)

	return cmd
}

func makeKVGetRunFunc(cfg *config.Config) runFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		kcfg, err := clientcmd.BuildConfigFromFlags("", cfg.KVKubeconfig())
		if err != nil {
			return wrapper.Wrap(err)
		}
		clientset, err := kubernetes.NewForConfig(kcfg)
		if err != nil {
			return wrapper.Wrap(err)
		}

		for _, name := range args {
			nlog := log.With(NameLabel, name)
			opts := metav1.ListOptions{
				LabelSelector: NameLabel + "=" + name,
			}
			list, err := clientset.CoreV1().Secrets("secrets").List(ctx, opts)
			if err != nil {
				return wrapper.Wrap(err)
			}
			nlog.DebugContext(ctx, "results", "len", len(list.Items))
			if len(list.Items) < 1 {
				continue
			}

			sort.Slice(list.Items, func(i, j int) bool {
				return list.Items[i].CreationTimestamp.Before(&list.Items[j].CreationTimestamp)
			})

			// most recent secret is the one we want.
			secret := list.Items[len(list.Items)-1]

			for k, v := range secret.Data {
				nlog.DebugContext(ctx, "data", "name", secret.Name, "key", k, "len", len(v))
			}

			if len(secret.Data) > 0 {
				cfg.Println(secret.Name)
			}
			for k, v := range secret.Data {
				cfg.Printf("-- %s --\n", k)
				cfg.Write(ensureNewline(v))
			}
		}
		return nil
	}
}

func ensureNewline(b []byte) []byte {
	if len(b) > 0 && b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	return b
}

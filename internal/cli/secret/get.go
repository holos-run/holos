package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const printFlagName = "print-key"

func NewGetCmd(hc *holos.Config) *cobra.Command {
	cmd := command.New("secrets NAME [--to-file=destination]")
	cmd.Aliases = []string{"secret"}
	cmd.Args = cobra.MinimumNArgs(0)
	cmd.Short = "get secrets stored in the management cluster"

	cfg, flagSet := newConfig()
	flagSet.Var(&cfg.files, "to-file", "extract files from the secret")
	cfg.printFile = flagSet.String(printFlagName, "", "print one key from the secret")
	cfg.extract = flagSet.Bool("extract-all", false, "extract all files from the secret")
	cfg.extractTo = flagSet.String("extract-to", ".", "extract to directory")

	cmd.Flags().SortFlags = false
	cmd.Flags().AddFlagSet(flagSet)
	cmd.RunE = makeGetRunFunc(hc, cfg)
	return cmd
}

func makeGetRunFunc(hc *holos.Config, cfg *config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		namespace := *cfg.namespace
		ctx := cmd.Context()
		log := logger.FromContext(ctx).With("namespace", namespace)

		cs, err := hc.ProvisionerClientset()
		if err != nil {
			return err
		}

		// List secrets if no arguments.
		if len(args) == 0 {
			return listSecrets(cmd.Context(), hc, namespace)
		}

		// Get each secret.
		for _, secretName := range args {
			log := log.With(NameLabel, secretName)
			opts := metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", NameLabel, secretName),
			}
			list, err := cs.CoreV1().Secrets(namespace).List(ctx, opts)
			if err != nil {
				return errors.Wrap(err)
			}

			log.DebugContext(ctx, "results", "len", len(list.Items))
			if len(list.Items) < 1 {
				return errors.Wrap(fmt.Errorf("not found: %v", secretName))
			}

			// Sort oldest first.
			sort.Slice(list.Items, func(i, j int) bool {
				return list.Items[i].CreationTimestamp.Before(&list.Items[j].CreationTimestamp)
			})

			// Get the most recent.
			secret := list.Items[len(list.Items)-1]
			log = log.With("secret", secret.Name)

			// Extract the data keys (file names).
			keys := make([]string, 0, len(secret.Data))
			for k, v := range secret.Data {
				keys = append(keys, k)
				log.DebugContext(ctx, "data", "name", secret.Name, "key", k, "len", len(v))
			}

			// Extract specified files or all files.
			toExtract := cfg.files
			if *cfg.extract {
				toExtract = keys
			}

			printFile := *cfg.printFile
			if len(toExtract) == 0 {
				if printFile == "" { // print all data keys
					data := make(map[string]string)
					for k, v := range secret.Data {
						data[k] = string(v)
					}
					b, err := json.MarshalIndent(data, "", "  ")
					if err != nil {
						return errors.Wrap(err)
					}
					b = util.EnsureNewline(b)
					hc.Write(b)
				} else { // print named data keys keys
					if data, found := secret.Data[printFile]; found {
						hc.Write(data)
					} else {
						err := fmt.Errorf("cannot print: want %s have %v: did you mean --extract-all or --%s=name", printFile, keys, printFlagName)
						return errors.Wrap(err)
					}
				}
			}

			// Iterate over --to-file values.
			for _, name := range toExtract {
				data, found := secret.Data[name]
				if !found {
					err := fmt.Errorf("%s not found in %v", name, keys)
					return errors.Wrap(err)
				}
				path := filepath.Join(*cfg.extractTo, name)
				if err := os.WriteFile(path, data, 0666); err != nil {
					return errors.Wrap(fmt.Errorf("could not write %s: %w", path, err))
				}
				log.InfoContext(ctx, "wrote: "+path, "name", name, "bytes", len(data))
			}
		}

		return nil
	}
}

// listSecrets lists holos secrets in the provisioner cluster
func listSecrets(ctx context.Context, hc *holos.Config, namespace string) error {
	cs, err := hc.ProvisionerClientset()
	if err != nil {
		return err
	}
	selector := metav1.ListOptions{LabelSelector: NameLabel}
	secrets, err := cs.CoreV1().Secrets(namespace).List(ctx, selector)
	if err != nil {
		return errors.Wrap(err)
	}
	secretNames := make(map[string]bool)
	for _, secret := range secrets.Items {
		if labelValue, ok := secret.Labels[NameLabel]; ok {
			secretNames[labelValue] = true
		}
	}
	for secretName := range secretNames {
		hc.Println(secretName)
	}
	return nil
}

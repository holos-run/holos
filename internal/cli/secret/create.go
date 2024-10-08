package secret

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/util/hash"
	"sigs.k8s.io/yaml"
)

func NewCreateCmd(hc *holos.Config) *cobra.Command {
	cmd := command.New("secret NAME [--from-file=source]")
	cmd.Aliases = []string{"secrets", "sec"}
	cmd.Args = cobra.ExactArgs(1)
	cmd.Short = "Create a holos secret from files or directories"

	cfg, flagSet := newConfig()
	flagSet.Var(&cfg.files, "from-file", "store files as keys in the secret")
	cfg.dryRun = flagSet.Bool("dry-run", false, "dry run")
	cfg.appendHash = flagSet.Bool("append-hash", true, "append hash to kubernetes secret name")
	cfg.dataStdin = flagSet.Bool("data-stdin", false, "read data field as json from stdin if")
	cfg.trimTrailingNewlines = flagSet.Bool("trim-trailing-newlines", true, "trim trailing newlines if true")

	cmd.Flags().SortFlags = false
	cmd.Flags().AddFlagSet(flagSet)
	cmd.RunE = makeCreateRunFunc(hc, cfg)
	return cmd

}

func makeCreateRunFunc(hc *holos.Config, cfg *config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		secretName := args[0]
		secret := &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: *cfg.namespace,
				Labels:    map[string]string{NameLabel: secretName},
			},
			Data: make(secretData),
		}

		if *cfg.cluster != "" {
			clusterPrefix := fmt.Sprintf("%s-", *cfg.cluster)
			if !strings.HasPrefix(secretName, clusterPrefix) {
				const msg = "missing cluster name prefix"
				log.WarnContext(ctx, msg, "have", secretName, "want", clusterPrefix)
			}
		}

		if *cfg.dataStdin {
			log.InfoContext(ctx, "reading data keys from stdin...")
			var obj map[string]string
			data, err := io.ReadAll(hc.Stdin())
			if err != nil {
				return errors.Wrap(err)
			}
			err = yaml.Unmarshal(data, &obj)
			if err != nil {
				return errors.Wrap(err)
			}
			for k, v := range obj {
				secret.Data[k] = []byte(v)
			}
		}

		for _, file := range cfg.files {
			if err := filepath.WalkDir(file, makeWalkFunc(secret.Data, file, *cfg.trimTrailingNewlines)); err != nil {
				return errors.Wrap(err)
			}
		}

		if owner := os.Getenv("USER"); owner != "" {
			secret.Labels[OwnerLabel] = owner
		}
		if *cfg.cluster != "" {
			secret.Labels[ClusterLabel] = *cfg.cluster
		}

		if *cfg.appendHash {
			if secretHash, err := hash.SecretHash(secret); err != nil {
				return errors.Wrap(err)
			} else {
				secret.Name = fmt.Sprintf("%s-%s", secret.Name, secretHash)
			}
		}

		if *cfg.dryRun {
			out, err := yaml.Marshal(secret)
			if err != nil {
				return errors.Wrap(err)
			}
			hc.Write(out)
			return nil
		}

		cs, err := hc.ProvisionerClientset()
		if err != nil {
			return errors.Wrap(err)
		}
		secret, err = cs.CoreV1().
			Secrets(*cfg.namespace).
			Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(err)
		}

		log.InfoContext(ctx, "created: "+secret.Name, "secret", secret.Name, "name", secretName, "namespace", secret.Namespace)
		return nil
	}
}

func makeWalkFunc(data secretData, root string, trimNewlines bool) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Depth is the count of path separators from the root
		depth := strings.Count(path[len(root):], string(filepath.Separator))

		if depth > 1 {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			key := filepath.Base(path)
			if data[key], err = os.ReadFile(path); err != nil {
				return errors.Wrap(err)
			}
			if trimNewlines {
				data[key] = bytes.TrimRight(data[key], "\r\n")
			}
		}

		return nil
	}
}

package kv

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/cli/secret"
	"github.com/holos-run/holos/pkg/errors"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/spf13/cobra"
	"golang.org/x/tools/txtar"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/util/hash"
	"sigs.k8s.io/yaml"
)

type putConfig struct {
	secretName *string
	file       *string
	dryRun     *bool
}

func newPutCmd(cfg *holos.Config) *cobra.Command {
	cmd := command.New("put")
	cmd.Args = cobra.MinimumNArgs(0)
	cmd.Short = "put a secret from stdin or file args"
	cmd.Flags().SortFlags = false

	pcfg := putConfig{}
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	pcfg.secretName = flagSet.String("name", "", "secret name to use instead of txtar comment")
	pcfg.file = flagSet.String("file", "", "file name to use instead of txtar path")
	pcfg.dryRun = flagSet.Bool("dry-run", false, "print to standard output instead of creating")

	cmd.Flags().AddGoFlagSet(flagSet)
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())
	cmd.RunE = makePutRunFunc(cfg, pcfg)

	return cmd
}

func makePutRunFunc(cfg *holos.Config, pcfg putConfig) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		a := &txtar.Archive{}

		// Add stdin to the archive.
		if len(args) == 0 {
			data, err := io.ReadAll(cfg.Stdin())
			if err != nil {
				return errors.Wrap(err)
			}

			if *pcfg.file != "" {
				file := txtar.File{
					Name: *pcfg.file,
					Data: data,
				}
				a.Files = append(a.Files, file)
			} else {
				a = txtar.Parse(data)
			}
		}

		// Do we have a secret name?
		if *pcfg.secretName != "" {
			a.Comment = []byte(*pcfg.secretName)
		}
		if len(a.Comment) == 0 {
			// Use the first argument if not
			if len(args) > 0 {
				a.Comment = []byte(filepath.Base(args[0]))
			} else {
				err := fmt.Errorf("missing secret name from name, args, or txtar comment")
				return errors.Wrap(err)
			}
		}

		head, _, _ := bytes.Cut(a.Comment, []byte("\n"))
		secretName := string(head)

		// Add files from the filesystem to the archive
		for _, name := range args {
			if err := filepath.WalkDir(name, makeWalkFunc(a, name)); err != nil {
				return errors.Wrap(err)
			}
		}

		log := logger.FromContext(cmd.Context())
		ctx := cmd.Context()

		// Nothing to do?
		if len(a.Files) == 0 {
			log.WarnContext(ctx, "nothing to do")
			return nil
		}

		// Create the secret.
		secret, err := createSecret(ctx, cfg, pcfg, a, secretName)
		if err != nil {
			return errors.Wrap(err)
		}

		if *pcfg.dryRun {
			data, err := yaml.Marshal(secret)
			if err != nil {
				return errors.Wrap(err)
			}
			cfg.Println(string(data))
			return nil
		}

		// Make the API call
		cs, err := newClientSet(cfg)
		if err != nil {
			return errors.Wrap(err)
		}

		secret, err = cs.CoreV1().Secrets(cfg.KVNamespace()).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(err)
		}

		log.InfoContext(ctx, "created: "+secret.Name, "secret", secret.Name, "name", secretName, "namespace", secret.Namespace)
		return nil
	}
}

func createSecret(ctx context.Context, cfg *holos.Config, pcfg putConfig, a *txtar.Archive, secretName string) (*v1.Secret, error) {
	secretData := make(map[string][]byte)
	for _, file := range a.Files {
		secretData[file.Name] = file.Data
	}

	labels := map[string]string{secret.NameLabel: secretName}
	if owner := os.Getenv("USER"); owner != "" {
		labels[secret.OwnerLabel] = owner
	}
	if cluster := cfg.ClusterName(); cluster != "" {
		labels[secret.ClusterLabel] = cluster
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   secretName,
			Labels: labels,
		},
		Data: secretData,
	}

	secretHash, err := hash.SecretHash(secret)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	secret.Name = fmt.Sprintf("%s-%s", secret.Name, secretHash)

	return secret, nil
}

func makeWalkFunc(a *txtar.Archive, rootDir string) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Depth is the count of path separators from the root
		depth := strings.Count(path[len(rootDir):], string(filepath.Separator))

		if depth > 1 {
			if d.IsDir() {
				return filepath.SkipDir
			}
		}

		if !d.IsDir() {
			if file, err := file(path); err != nil {
				return errors.Wrap(err)
			} else {
				file.Name = filepath.Base(path)
				a.Files = append(a.Files, file)
			}
		}

		return nil
	}
}

func file(path string) (file txtar.File, err error) {
	file.Name = path
	file.Data, err = os.ReadFile(path)
	return
}

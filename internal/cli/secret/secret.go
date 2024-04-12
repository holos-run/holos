package secret

import (
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/pflag"
)

const NameLabel = "holos.run/secret.name"
const OwnerLabel = "holos.run/owner.name"
const ClusterLabel = "holos.run/cluster.name"

type secretData map[string][]byte

type config struct {
	files                holos.StringSlice
	printFile            *string
	extract              *bool
	dryRun               *bool
	appendHash           *bool
	dataStdin            *bool
	trimTrailingNewlines *bool
	cluster              *string
	namespace            *string
	extractTo            *string
}

func newConfig() (*config, *pflag.FlagSet) {
	cfg := &config{}
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
	cfg.namespace = flagSet.StringP("namespace", "n", holos.DefaultProvisionerNamespace, "namespace in the provisioner cluster")
	cfg.cluster = flagSet.String("cluster-name", "", "cluster name selector")
	return cfg, flagSet
}

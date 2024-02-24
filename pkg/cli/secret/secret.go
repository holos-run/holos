package secret

import (
	"flag"
	"github.com/holos-run/holos/pkg/holos"
)

const NameLabel = "holos.run/secret.name"
const OwnerLabel = "holos.run/owner.name"
const ClusterLabel = "holos.run/cluster.name"

type secretData map[string][]byte

type config struct {
	files     holos.StringSlice
	printFile *string
	extract   *bool
	dryRun    *bool
	cluster   *string
	namespace *string
}

func newConfig() (*config, *flag.FlagSet) {
	cfg := &config{}
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg.namespace = flagSet.String("namespace", holos.DefaultProvisionerNamespace, "namespace in the provisioner cluster")
	cfg.cluster = flagSet.String("cluster-name", "", "cluster name selector")
	return cfg, flagSet
}

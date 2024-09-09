package holos

import (
	"encoding/json"
	dto "github.com/holos-run/holos/service/gen/holos/object/v1alpha1:object"
)

// _ClusterName is the --cluster-name flag value provided by the holos cli.
_ClusterName: string @tag(cluster, type=string)

// _PlatformConfig represents all of the data passed from holos to cue, used to
// carry the platform and project models.
_PlatformConfig:     dto.#PlatformConfig & json.Unmarshal(_PlatformConfigJSON)
_PlatformConfigJSON: string | *"{}" @tag(platform_config, type=string)

package holos

import (
	"encoding/json"
	dto "github.com/holos-run/holos/service/gen/holos/object/v1alpha1:object"
)

// Note, tags should have a reasonable default value to easily use cue eval and
// cue export without needing to make a bunch of decisions about tag values.

// _ClusterName is the --cluster-name flag value provided by the holos cli.
_ClusterName: string | *"no-name" @tag(cluster, type=string)

// _PlatformConfig represents all of the data passed from holos to cue, used to
// carry the platform and project models.
_PlatformConfig:     dto.#PlatformConfig & json.Unmarshal(_PlatformConfigJSON)
_PlatformConfigJSON: string | *"{}" @tag(platform_config, type=string)

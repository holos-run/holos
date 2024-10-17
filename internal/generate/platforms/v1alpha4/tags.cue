package holos

import (
	"encoding/json"
	dto "github.com/holos-run/holos/service/gen/holos/object/v1alpha1:object"
)

// Note, tags should have a reasonable default value to easily use cue eval and
// cue export without needing to make a bunch of decisions about tag values.

// _PlatformConfig represents all of the data passed from holos to cue, used to
// carry the platform and project models.
_PlatformConfig:     dto.#PlatformConfig & json.Unmarshal(_PlatformConfigJSON)
_PlatformConfigJSON: string | *"{}" @tag(holos_platform_config, type=string)

_Tags: {
	// Reserved tags
	name:      string | *"no-name"      @tag(holos_name, type=string)
	cluster:   string | *"no-cluster"   @tag(holos_cluster, type=string)
	component: string | *"no-component" @tag(holos_component, type=string)
	// User defined tags
	environment: string | *"no-environment" @tag(environment, type=string)
	project:     string | *"no-project"     @tag(project, type=string)
}

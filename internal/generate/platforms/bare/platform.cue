package holos

import "encoding/json"

import v1 "github.com/holos-run/holos/api/v1alpha1"

import dto "github.com/holos-run/holos/service/gen/holos/object/v1alpha1:object"

// _PlatformConfig represents all of the data passed from holos to cue.
// Intended to carry the platform model and project models.
_PlatformConfig:     dto.#PlatformConfig & json.Unmarshal(_PlatformConfigJSON)
_PlatformConfigJSON: string | *"{}" @tag(platform_config, type=string)

// _Platform provides a platform resource to the holos cli for rendering.  The
// field is hidden because most components need to refer to platform data,
// specifically the platform model and the project models.  The platform
// resource itself is output once when rendering the entire platform, see the
// platform/ subdirectory.
_Platform: v1.#Platform & {
	metadata: {
		name: string | *"bare" @tag(platform_name, type=string)
	}

	// spec is the platform specification
	spec: {
		// model represents the web form values provided by the user.
		model: _PlatformConfig.platform_model
		components: [for c in _components {c}]

		_components: [string]: v1.#PlatformSpecComponent
		_components: {
			for WorkloadCluster in _Clusters.Workload {
				"\(WorkloadCluster)-configmap": {
					path:    "components/configmap"
					cluster: WorkloadCluster
				}
			}
		}
	}
}

// _Clusters represents the clusters in the platform.  The default values are
// intended to be provided by the user in a file which is not written over by
// `holos generate`.
_Clusters: {
	Workload: [...string] | *["mycluster"]
}

package holos

import "encoding/json"

import v1 "github.com/holos-run/holos/api/v1alpha1"

import dto "github.com/holos-run/holos/service/gen/holos/object/v1alpha1:object"

// #PlatformConfig represents all of the data passed from holos to cue, used to
// cary the platform and project models.
#PlatformConfig:     dto.#PlatformConfig & json.Unmarshal(_PlatformConfigJSON)
_PlatformConfigJSON: string | *"{}" @tag(platform_config, type=string)

// #Cluster represents a single cluster in the platform.
#Cluster: name: string

// #Fleet represents a grouping of similar clusters.  A platform is usually
// composed of a workload fleet and a management fleet.
#Fleet: {
	name: string
  clusters: [Name=string]: #Cluster & { name: Name }
}

// #Fleets represents all the fleets in the platform.
#Fleets: [Name=string]: #Fleet & { name: Name }

// #Platform represents a platform build.
#Platform: {
  Name: string | *"holos" @tag(platform_name, type=string)

  // Components represent the components to include in the platform.
  Components: [string]: v1.#PlatformSpecComponent

  Output: v1.#Platform & {
    metadata: name: Name

    spec: {
      // model represents the web form values provided by the user.
      model: #PlatformConfig.platform_model
      components: [for c in Components {c}]
    }
  }
}

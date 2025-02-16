package holos

import (
	"path"
	"encoding/json"
	"holos.example/schema/my-app/deployment"
)

parameters: {
	config: _ @tag(config)
}
config: deployment.#Config & json.Unmarshal(parameters.config)

// component represents the holos component definition, which produces a
// BuildPlan for holos to execute, rendering the manifests.
component: #Helm & {
	// See step1.cue and step3.cue for where values are mixed in.
	Chart: {
		version: "0.1.0"
		repository: {
			name: "my-app"
			url:  "https://chart.holos.example/my-app"
		}
	}
}

// holos represents the output for the holos command line to process.  The holos
// command line processes a BuildPlan to render the helm chart component.
//
// Use the holos show buildplans command to see the BuildPlans that holos render
// platform renders.
holos: component.BuildPlan

#ComponentConfig: {
	Name:          _
	OutputBaseDir: _

	_ArgoAppName: "\(config.customer)-\(config.application)-\(config.cluster)"
	_GitOpsArtifact: artifact: path.Join(["apps", "\(_ArgoAppName)-application.gen.yaml"], path.Unix)
	let ResourcesPath = path.Join(["deploy", OutputBaseDir, "components", Name], path.Unix)
	_ArgoApplication: spec: source: path: ResourcesPath
}

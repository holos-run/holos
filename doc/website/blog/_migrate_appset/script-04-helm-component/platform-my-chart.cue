package holos

// Imports ./config/environments/*.cue as the environments cue package.  The
// package exposes ./config/environments/**/config.json files via the
// environments.config struct
import "holos.example/config/environments"

// Manage my-chart for each of the three environments.  Platform components are
// rendered by the holos render platform command.
//
// Use the following command command to inspect the Platform spec holos render
// platform processes.
//
//  holos show platform
//
// CONFIG represents each migrated environments/**/config.json file.
for CONFIG in environments.config {
	// Add one holos component for each config.json file to the
	// Platform.spec.components list.
	Platform: Components: "\(CONFIG.env)-my-chart": #MyChart & {
		parameters: {
			env:     CONFIG.env
			region:  CONFIG.region
			type:    CONFIG.type
			version: CONFIG.version
			chart:   CONFIG.chart
		}
	}
}

// #MyChart defines a re-usable way to manage my-chart across qa, staging, and
// production.
#MyChart: {
	name: "my-chart"
	path: "components/my-chart"
	// CUE supports constraints, here we constrain environment to one of three
	// possible values.
	parameters: {
		// Embed the config.json schema (env, region, type, version, chart fields)
		environments.#Config

		// Define the env field here as any value (_) so we can refer to it.
		// Otherwise cue complains #MyChart.parameters.outputBaseDir: reference
		// "env" not found
		env: _
		// Write output manifests organized by environment env in this case refers
		// to the env field defined by environments.#Config
		outputBaseDir: "environments/\(env)"
	}
	// CUE supports string substitution.
	annotations: "app.holos.run/description": "my-chart \(parameters.chart) for environment \(parameters.env)"
	// Selector labels
	labels: env: parameters.env
}

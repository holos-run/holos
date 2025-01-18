@if(flatten && step1 || step2 || step3)
package holos

import (
	"encoding/json"

	// The deployment config.json configs reside in this package.
	"holos.example/config/my-app/deployment"
)

// #MyApp defines a re-usable way to manage my-app for each deployment config.
#MyApp: {
	// _config represents the concrete values from the deployment config.json
	_config: _

	// Render components/my-app by default.  See render-values.cue for where this
	// gets overridden in subsequent steps of the migration.
	name: string | *"my-app"
	path: "components/\(name)"

	// CUE supports constraints, here we constrain environment to one of three
	// possible values.
	parameters: {
		// For use in components/my-app and components/render-values
		config: json.Marshal(_config)

		// For use in components/componentconfig-gitops.cue
		env: _config.env

		// The output is the reverse of the deployment config filesystem structure,
		// cluster then customer instead of customer then cluster, reflecting the
		// perspective of the platform team compared with the perspective of the app
		// team.  It could be either way.
		outputBaseDir: "clusters/\(_config.cluster)/customers/\(_config.customer)"
	}
	// The app.holos.run/description annotation configures holos render platform
	// log messages.
	annotations: "app.holos.run/description": string | *"\(name) for \(_config.customer) cluster \(_config.cluster)"
	// Selector labels, useful to quickly render combinations of customer,
	// cluster, and app
	labels: customer: _config.customer
	labels: cluster:  _config.cluster
	labels: app:      _config.application
}

// Add one holos component for each config.json file to the
// Platform.spec.components list.
for KEY, CONFIG in deployment.config {
	Platform: Components: (KEY): #MyApp & {
		_config: CONFIG
	}
}

@if(flatten && step2)
package holos

#MyApp: {
	// Renders components/render-values instead of the default components/my-app
	// configured in my-app.cue
	name: "render-values"

	// Bring the config.json values into scope so we can reference them.
	_config: _

	// Adjust the output of holos render platform so it's clear we're flatting
	// values in step 2.
	annotations: "app.holos.run/description": "flattened values for \(_config.customer) \(_config.cluster)"
}

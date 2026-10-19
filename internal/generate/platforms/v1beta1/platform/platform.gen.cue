package holos

// holos represents the platform resource constructed from registered
// components. The platform runtime stays local so init does not vendor the
// Holos CUE SDK.
holos: platform.resource

platform: {
	name: string | *"default"
	components: {[string]: {
		name: string
		path: string
		labels?: {[string]: string}
		annotations?: {[string]: string}
	}} | *{}
	let platformName = name
	let registeredComponents = components
	resource: {
		metadata: name: platformName
		spec: components: [for _, component in registeredComponents {component}]
	}
}

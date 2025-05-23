package holos

platform: spec: components: {
	TaskSet1: {
		name: "TaskSet1"
		path: "components/taskset"
		// Parameters to pass into the component.
		parameters: index: "1"
		// In v1alpha5 "component" was ambiguous.  We disambiguate in v1alpha6 by
		// naming the output of a platform component an "instance"
		labels: "app.holos.run/component.instance": name
		// The component the instance is derived from.
		labels: "app.holos.run/component.path": path
	}
}

holos: platform

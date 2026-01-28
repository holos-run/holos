package holos

// This component produces a TaskSet at the top level holos field
// The TaskSet structure is to be defined
holos: {
	metadata: name: "taskset-example"
	
	// TaskSet structure placeholder - to be defined
	// This will represent the new v1alpha6 TaskSet that replaces BuildPlan
	taskSet: {
		// Example structure for discussion
		tasks: {
			// Tasks as a struct instead of list for better composition
			generateManifests: {
				type: "generator"
				// Additional fields to be defined
			}
			transformManifests: {
				type: "transformer"
				dependsOn: ["generateManifests"]
				// Additional fields to be defined
			}
			validateManifests: {
				type: "validator"
				dependsOn: ["transformManifests"]
				// Additional fields to be defined
			}
		}
	}
}
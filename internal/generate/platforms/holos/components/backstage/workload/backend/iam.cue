package holos

_BackstageIAMConfig: {
	groupAdmin: {
		// https://backstage.io/docs/features/software-catalog/descriptor-format#kind-group
		apiVersion: "backstage.io/v1alpha1"
		kind:       "Group"
		metadata: name: "prod-cluster-admin"
		spec: {
			type: "team"
			children: []
		}
	}

	user1: {
		// https://backstage.io/docs/features/software-catalog/descriptor-format#kind-user
		apiVersion: "backstage.io/v1alpha1"
		kind:       "User"
		metadata: name: "jeff"
		spec: {
			profile: email: "jeff@openinfrastructure.co"
			memberOf: ["prod-cluster-admin"]
		}
	}

	user2: {
		// https://backstage.io/docs/features/software-catalog/descriptor-format#kind-user
		apiVersion: "backstage.io/v1alpha1"
		kind:       "User"
		metadata: name: "gary"
		spec: {
			profile: email: "gary@openinfrastructure.co"
			memberOf: ["prod-cluster-admin"]
		}
	}

	user3: {
		// https://backstage.io/docs/features/software-catalog/descriptor-format#kind-user
		apiVersion: "backstage.io/v1alpha1"
		kind:       "User"
		metadata: name: "nate"
		spec: {
			profile: email: "nate@openinfrastructure.co"
			memberOf: ["prod-cluster-admin"]
		}
	}
}

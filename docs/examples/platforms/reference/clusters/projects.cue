package holos

_Projects: #Projects & {
	holos: {
		clusters: {
			k1: _
			k2: _
		}
		environments: {
			prod: stage: "prod"
			dev: stage:  "dev"
			jeff: stage: dev.stage
			gary: stage: dev.stage
			nate: stage: dev.stage
		}
	}

	iam: {
		clusters: {
			core1: _
			core2: _
		}
	}
}

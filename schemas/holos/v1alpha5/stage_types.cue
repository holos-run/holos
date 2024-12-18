package v1alpha5

// #Stage represents a deployment stage, often called an environment.  Examples
// are dev, test, uat, prod-east, prod-west.
#Stage: {
	// name represents the stage name, e.g. "dev" or "prod-us-east"
	name:  #Name
	#Name: string

	// prior represents the prior stage in the promotion process or the special
	// value "direct" if there is no prior stage.
	prior: #Name | *"direct"

	// tier represents the tier of the stage, usually prod or nonprod.
	tier:  #Tier
	#Tier: "prod" | "nonprod"

	// parameters represents stage specific component parameters.
	parameters: #Parameters
	#Parameters: [string]: string
}

// #Stages represents a collection of #Stage values organized by name.
#Stages: [NAME=#Stage.#Name]: #Stage & {name: NAME}

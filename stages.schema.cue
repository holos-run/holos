package holos

#Stage: {
	#Name: string
	name:  #Name

	#Tier: "prod" | "nonprod"
	tier:  #Tier

	// parameters represents stage specific parameters passed to each component.
	#Parameters: [string]: string
	parameters: #Parameters
}

#Stages: {
	[NAME=#Stage.#Name]: #Stage & {name: NAME}
}

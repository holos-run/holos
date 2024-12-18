package holos

#Stage: {
	#Name: string
	name:  #Name

	#Tier: "prod" | "nonprod"
	tier:  #Tier
}

#Stages: {
	[NAME=#Stage.#Name]: #Stage & {
		name: NAME
	}
}

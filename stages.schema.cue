package holos

#Stage: {
	name: string
	tier: "prod" | "nonprod"
}

#Stages: {
	[NAME=string]: #Stage & {
		name: NAME
	}
}

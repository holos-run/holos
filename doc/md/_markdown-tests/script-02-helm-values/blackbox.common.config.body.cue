package holos

// Schema Definition
#Blackbox: {
	// host constrained to a lower case dns label
	host: string & =~"^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$"
	// port constrained to a valid range
	port: int & >0 & <=65535
}

// Concrete values must validate against the schema.
Blackbox: #Blackbox & {
	host: "blackbox"
	port: 9115
}

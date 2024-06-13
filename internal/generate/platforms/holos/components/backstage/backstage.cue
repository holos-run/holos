package holos

// _DBName is the database name used across multiple holos components in this project
_DBName: "backstage"

_Component: {
	metadata: name:      "backstage"
	metadata: namespace: "backstage"
	spec: hostname:      "backstage.admin.\(_ClusterName).\(_Platform.Model.org.domain)"
	spec: port:          7007
}

package v1

_apiVersion: "rbac.authorization.k8s.io/v1"

#Role: {
	apiVersion: _apiVersion
	kind:       "Role"
}

#RoleBinding: {
	apiVersion: _apiVersion
	kind:       "RoleBinding"
}

#ClusterRole: {
	apiVersion: _apiVersion
	kind:       "ClusterRole"
}

#ClusterRoleBinding: {
	apiVersion: _apiVersion
	kind:       "ClusterRoleBinding"
}

package holos

// Platform wide definitions
#BankOfHolos: {
	Frontend: Namespace: "bank-frontend"
	Backend: Namespace:  "bank-backend"
	Security: Namespace: "bank-security"

	// Resources to manage in each of the namespaces.
	Resources: #Resources
}

// Register namespaces
#Namespaces: (#BankOfHolos.Frontend.Namespace): _
#Namespaces: (#BankOfHolos.Backend.Namespace):  _
#Namespaces: (#BankOfHolos.Security.Namespace): _

// Register projects
#AppProjects: "bank-frontend": _
#AppProjects: "bank-backend":  _
#AppProjects: "bank-security": _

// Register HTTPRoutes.
// bank.example.com routes to Service frontend in the bank-frontend namespace.
#HTTPRoutes: bank: _backendRefs: frontend: namespace: #BankOfHolos.Frontend.Namespace

package holos

// Platform wide definitions
#BankOfHolos: {
	Frontend: Namespace: "bank-frontend"
	Ledger: Namespace:   "bank-ledger"
	Identity: Namespace: "bank-identity"
	Tools: Namespace:    "bank-tools"

	// Resources to manage in each of the namespaces.
	Resources: #Resources
}

// Register namespaces
#Namespaces: (#BankOfHolos.Frontend.Namespace): _
#Namespaces: (#BankOfHolos.Ledger.Namespace):   _
#Namespaces: (#BankOfHolos.Identity.Namespace): _
#Namespaces: (#BankOfHolos.Tools.Namespace):    _

// Register projects
#AppProjects: "bank-frontend": _
#AppProjects: "bank-ledger":   _
#AppProjects: "bank-identity": _
#AppProjects: "bank-tools":    _

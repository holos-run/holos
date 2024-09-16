package holos

// Manage on workload clusters only
for Cluster in #Fleets.workload.clusters {
	// Owned by the security team
	#Platform: Components: "\(Cluster.name)/bank-secrets": {
		path:    "projects/bank-of-holos/security/components/bank-secrets"
		cluster: Cluster.name
	}

	// Owned by the frontend team
	#Platform: Components: "\(Cluster.name)/bank-frontend": {
		path:    "projects/bank-of-holos/frontend/components/bank-frontend"
		cluster: Cluster.name
	}

	// Owned by the backend team
	#Platform: Components: "\(Cluster.name)/bank-backend-config": {
		path:    "projects/bank-of-holos/backend/components/bank-backend-config"
		cluster: Cluster.name
	}
	#Platform: Components: "\(Cluster.name)/bank-accounts-db": {
		path:    "projects/bank-of-holos/backend/components/bank-accounts-db"
		cluster: Cluster.name
	}
	#Platform: Components: "\(Cluster.name)/bank-userservice": {
		path:    "projects/bank-of-holos/backend/components/bank-userservice"
		cluster: Cluster.name
	}

	#Platform: Components: "\(Cluster.name)/bank-ledger-db": {
		path:    "projects/bank-of-holos/backend/components/bank-ledger-db"
		cluster: Cluster.name
	}
	#Platform: Components: "\(Cluster.name)/bank-ledger-writer": {
		path:    "projects/bank-of-holos/backend/components/bank-ledger-writer"
		cluster: Cluster.name
	}
	#Platform: Components: "\(Cluster.name)/bank-balance-reader": {
		path:    "projects/bank-of-holos/backend/components/bank-balance-reader"
		cluster: Cluster.name
	}
	#Platform: Components: "\(Cluster.name)/bank-transaction-history": {
		path:    "projects/bank-of-holos/backend/components/bank-transaction-history"
		cluster: Cluster.name
	}
	#Platform: Components: "\(Cluster.name)/bank-contacts": {
		path:    "projects/bank-of-holos/backend/components/bank-contacts"
		cluster: Cluster.name
	}
}

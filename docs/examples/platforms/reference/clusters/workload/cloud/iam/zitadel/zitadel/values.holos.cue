package holos

#Values: {
	// Database credentials
	// Refer to https://access.crunchydata.com/documentation/postgres-operator/5.2.0/architecture/user-management/
	// Refer to https://zitadel.com/docs/self-hosting/manage/database#postgres
	env: [
		// Connection
		{
			name: "ZITADEL_DATABASE_POSTGRES_HOST"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-\(_DBName)"
			valueFrom: secretKeyRef: key:  "host"
		},
		{
			name: "ZITADEL_DATABASE_POSTGRES_PORT"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-\(_DBName)"
			valueFrom: secretKeyRef: key:  "port"
		},
		{
			name: "ZITADEL_DATABASE_POSTGRES_DATABASE"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-\(_DBName)"
			valueFrom: secretKeyRef: key:  "dbname"
		},
		// The <db>-pguser-<db> secret contains creds for the unpriviliged zitadel user
		{
			name: "ZITADEL_DATABASE_POSTGRES_USER_USERNAME"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-\(_DBName)"
			valueFrom: secretKeyRef: key:  "user"
		},
		{
			name: "ZITADEL_DATABASE_POSTGRES_USER_PASSWORD"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-\(_DBName)"
			valueFrom: secretKeyRef: key:  "password"
		},
		// The postgres component configures privileged postgres user creds.
		{
			name: "ZITADEL_DATABASE_POSTGRES_ADMIN_USERNAME"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-postgres"
			valueFrom: secretKeyRef: key:  "user"
		},
		{
			name: "ZITADEL_DATABASE_POSTGRES_ADMIN_PASSWORD"
			valueFrom: secretKeyRef: name: "\(_DBName)-pguser-postgres"
			valueFrom: secretKeyRef: key:  "password"
		},

	]

	// Refer to https://zitadel.com/docs/self-hosting/manage/database
	zitadel: {
		// Zitadel master key
		masterkeySecretName: "zitadel-masterkey"
		// Note the tls configuration is a challenge to use externally issued certs from the provsioner cluster.
		// We intentionally use pgo managed certs and intend to backup the ca key to the provisioner and restore it for
		// cross cluster replication.  The problems seemed to arise from specifying the user and admin tls secrets in
		// addition to the ca cert secret.
		dbSslCaCrtSecret: "\(_DBName)-cluster-cert"

		// All settings: https://zitadel.com/docs/self-hosting/manage/configure#runtime-configuration-file
		// Helm interface: https://github.com/zitadel/zitadel-charts/blob/zitadel-7.4.0/charts/zitadel/values.yaml#L20-L21
		configmapConfig: {
			// NOTE: You can change the ExternalDomain, ExternalPort and ExternalSecure
			// configuration options at any time. However, for ZITADEL to be able to
			// pick up the changes, you need to rerun ZITADELs setup phase.  Do so with
			// kubectl delete job zitadel-setup, then re-apply the new config.
			//
			// https://zitadel.com/docs/self-hosting/manage/custom-domain
			ExternalSecure: true
			ExternalDomain: #ExternalDomain
			ExternalPort:   443
			TLS: Enabled: false

			// Database connection credentials are injected via environment variables from the db-pguser-db secret.
			Database: postgres: {
				MaxOpenConns:    25
				MaxIdleConns:    10
				MaxConnLifetime: "1h"
				MaxConnIdleTime: "5m"
				// verify-full verifies the host name matches cert dns names in addition to root ca signature
				User: SSL: Mode:  "verify-full"
				Admin: SSL: Mode: "verify-full"
			}
		}
	}
}

package holos

#Values: {

	// https://raw.githubusercontent.com/zitadel/zitadel-charts/main/examples/4-cockroach-secure/zitadel-values.yaml
	zitadel: {
		masterkeySecretName: "zitadel-masterkey"
		// https://github.com/zitadel/zitadel-charts/blob/zitadel-7.4.0/charts/zitadel/templates/configmap.yaml#L13
		configmapConfig: {
			// NOTE: You can change the ExternalDomain, ExternalPort and ExternalSecure
			// configuration options at any time. However, for ZITADEL to be able to
			// pick up the changes, you need to rerun ZITADELs setup phase.  Do so with
			// kubectl delete job zitadel-setup, then re-apply the new config.
			//
			// https://zitadel.com/docs/self-hosting/manage/custom-domain
			ExternalDomain: #ExternalDomain
			ExternalPort:   443
			ExternalSecure: true
			TLS: Enabled: false
			Database: Cockroach: {
				Host: #DB.Host
				User: SSL: Mode:  "verify-full"
				Admin: SSL: Mode: "verify-full"
			}
		}

		// Managed by crdb component
		dbSslCaCrtSecret:    "cockroach-ca"
		dbSslAdminCrtSecret: "cockroachdb-root"
		// Managed by this component
		dbSslUserCrtSecret: "cockroachdb-zitadel"
	}

}

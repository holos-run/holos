package holos

// Imported from https://github.com/holos-run/portal/blob/d5127715fb4710b9b272768e6a1ce2ff122e693e/app-config.production.yaml

_BackstageProductionConfig: {
	app: {
		// Should be the same as backend.baseUrl when using the `app-backend` plugin.
		baseUrl: "${BASE_URL}"
	}

	backend: {
		// Note that the baseUrl should be the URL that the browser and other clients
		// should use when communicating with the backend, i.e. it needs to be
		// reachable not just from within the backend host, but from all of your
		// callers. When its value is "http://localhost:7007", it's strictly private
		// and can't be reached by others.
		baseUrl: "${BASE_URL}"
		// The listener can also be expressed as a single <host>:<port> string. In this case we bind to
		// all interfaces, the most permissive setting. The right value depends on your specific deployment.
		listen: ":7007"

		// config options: https://node-postgres.com/api/client
		database: {
			client: "pg"
			connection: {
				host:     "${POSTGRES_HOST}"
				port:     "${POSTGRES_PORT}"
				user:     "${POSTGRES_USER}"
				password: "${POSTGRES_PASSWORD}"
				ssl: ca: "${PGBOUNCER_CA_ROOT}"
			}
		}

		reading: allow: [{
			host: "holos.run"
		}, {
			host: "*.holos.run"
		}, {
			host: "openinfrastructure.co"
		}, {
			host: "*.openinfrastructure.co"
		}]
	}

	auth: {
		environment: "production"
		providers: {
			guest: null
			holosProxy: {
				issuer:            "https://login.holos.run"
				audience:          "269746002573969304"
				oidcIdTokenHeader: "x-oidc-id-token"
				signIn: resolvers: [{
					resolver: "emailMatchingUserEntityProfileEmail"
				}, {
					resolver: "signInWithoutCatalogUser"
				}]
			}
		}
	}

	catalog: {
		// Overrides the default list locations from app-config.yaml
		// Refer to https://backstage.io/docs/features/software-catalog/#adding-components-to-the-catalog
		//
		// NOTE: In production, CWD is /app  This is different than development where
		// CWD is ./packages/backend, As a result, entries cannot be copied verbatim
		// from app-config.yaml otherwise they will not resolve to the correct location.
		locations: [{
			// Initial iam User and Group data is expected to be provided in the Secret.
			// Local example data, file locations are relative to the backend process, typically `packages/backend`
			type:   "file"
			target: "/config/iam.yaml"
			rules: [{
				allow: ["User", "Group"]}]
		}]

		// GitHub Discovery
		// Refer to https://backstage.io/docs/integrations/github/discovery/#configuration
		providers: {
			github: {
				primaryOrg: {
					organization: string & _Platform.Model.github.primaryOrg
					catalogPath:  "/catalog-info.yaml"
					filters: {
						branch:     "main"
						repository: ".*" // Regex
					}
					// same options as in TaskScheduleDefinition
					schedule: {
						// supports cron, ISO duration, "human duration" as used in code
						frequency: minutes: 30
						// supports ISO duration, "human duration" as used in code
						timeout: minutes: 3
					}
				}
			}
		}
	}

	// Refers to ExternalSecret github-app-credentials.  See the readme.md for how this secret is produced.
	integrations: {
		github: [{
			host: "github.com"
			apps: [{"$include": "/secrets/github-app-credentials/github-app-credentials.yaml"}]
		}]
	}
}

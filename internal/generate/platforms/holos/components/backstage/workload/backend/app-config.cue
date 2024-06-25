package holos

// TODO This entire config should be removed for clarity.  Only the production
// config should be referenced by the Deployment.

_BackstageAppConfig: {
	app: {
		title:   "Holos Portal"
		baseUrl: "${BASE_URL}"
	}

	organization: name: "My Company"

	backend: {
		// Used for enabling authentication, secret is shared by all backend plugins
		// See https://backstage.io/docs/auth/service-to-service-auth for
		// information on the format
		// auth:
		//   keys:
		//     - secret: ${BACKEND_SECRET}
		baseUrl: "${BASE_URL}"
		listen: port: 7007
		// Uncomment the following host directive to bind to specific interfaces
		// host: 127.0.0.1
		csp: {
			"connect-src": ["'self'", "http:", "https:"]
		}
		// Content-Security-Policy directives follow the Helmet format: https://helmetjs.github.io/#reference
		// Default Helmet Content-Security-Policy values can be removed by setting the key to false
		cors: {
			origin: "${BASE_URL}"
			methods: ["GET", "HEAD", "PATCH", "POST", "PUT", "DELETE"]
			credentials: true
		}
	}
	// workingDirectory: /tmp # Use this to configure a working directory for the scaffolder, defaults to the OS temp-dir

	integrations: {
		github: [{
			host: "github.com"
			// This is a Personal Access Token or PAT from GitHub. You can find out how to generate this token, and more information
			// about setting up the GitHub integration here: https://backstage.io/docs/integrations/github/locations#configuration
			token: "${GITHUB_TOKEN}"
		}]
	}
	//## Example for how to add your GitHub Enterprise instance using the API:
	// - host: ghe.example.net
	//   apiBaseUrl: https://ghe.example.net/api/v3
	//   token: ${GHE_TOKEN}

	proxy: null
	//## Example for how to add a proxy endpoint for the frontend.
	//## A typical reason to do this is to handle HTTPS and CORS for internal services.
	// endpoints:
	//   '/test':
	//     target: 'https://example.com'
	//     changeOrigin: true
	// Reference documentation http://backstage.io/docs/features/techdocs/configuration
	// Note: After experimenting with basic setup, use CI/CD to generate docs
	// and an external cloud storage when deploying TechDocs for production use-case.
	// https://backstage.io/docs/features/techdocs/how-to-guides#how-to-migrate-from-techdocs-basic-to-recommended-deployment-approach
	techdocs: {
		builder: "local" // Alternatives - 'external'
		generator: {
			runIn: "docker"
		} // Alternatives - 'local'
		publisher: {
			type: "local"
		}
	} // Alternatives - 'googleGcs' or 'awsS3'. Read documentation for using alternatives.

	auth: {
		environment: "development"
		// see https://backstage.io/docs/auth/ to learn about auth providers
		providers: {
			// See https://backstage.io/docs/auth/guest/provider
			guest: {}
		}
	}

	scaffolder: null
	// see https://backstage.io/docs/features/software-templates/configuration for software template options

	catalog: {
		import: {
			entityFilename:        "catalog-info.yaml"
			pullRequestBranchName: "backstage-integration"
		}
		rules: [{allow: ["Component", "System", "API", "Resource", "Location"]}]
		locations: [
			{
				// Local example data, file locations are relative to the backend process, typically `packages/backend`
				type:   "file"
				target: "../../examples/entities.yaml"
			},
			{
				// Local example template
				type:   "file"
				target: "../../examples/template/template.yaml"
				rules: [{
					allow: ["Template"]}]
			},
			{
				// Local organizational data
				type:   "file"
				target: "../../org.yaml"
				rules: [{
					allow: ["User", "Group"]}]
			},
		]
	}
}

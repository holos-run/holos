@extern(embed)
@if(step1 || step2)
package holos

valueFiles: _ @embed(glob=values/*/*-values.yaml)

component: ValueFiles: [
	{
		name:   "common-values.yaml"
		values: valueFiles["values/11-common/common-values.yaml"]
	},
	{
		name:   "location-values.yaml"
		values: valueFiles["values/10-locations/\(config.location)-values.yaml"]
	},
	{
		name:   "region-values.yaml"
		values: valueFiles["values/09-regions/\(config.region)-values.yaml"]
	},
	{
		name:   "zone-values.yaml"
		values: valueFiles["values/08-zones/\(config.zone)-values.yaml"]
	},
	{
		name:   "scope-values.yaml"
		values: valueFiles["values/07-scopes/\(config.scope)-values.yaml"]
	},
	{
		name:   "tier-values.yaml"
		values: valueFiles["values/06-tiers/\(config.tier)-values.yaml"]
	},
	{
		name:   "env-values.yaml"
		values: valueFiles["values/05-environments/\(config.env)-values.yaml"]
	},
	{
		name:   "cluster-values.yaml"
		values: valueFiles["values/04-clusters/\(config.cluster)-values.yaml"]
	},
	{
		name:   "application-values.yaml"
		values: valueFiles["values/03-applications/\(config.application)-values.yaml"]
	},
	{
		name:   "namespace-values.yaml"
		values: valueFiles["values/02-namespaces/\(config.namespace)-values.yaml"]
	},
	{
		name:   "customer-values.yaml"
		values: valueFiles["values/01-customers/\(config.customer)-values.yaml"]
	},
]

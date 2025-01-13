@extern(embed)
package environments

// We use cue embed functionality as an equivalent replacement for
// ApplicationSet generators.
config: _ @embed(glob=*/config.json)
config: _ @embed(glob=staging/*/config.json)
config: _ @embed(glob=prod/*/config.json)
config: _ @embed(glob=integration/*/config.json)

// With CUE we can constrain the data with a schema.
config: [FILEPATH=string]: #Config

// #Config defines the schema of each config.json file.
#Config: {
	env:     "qa" | "integration-gpu" | "integration-non-gpu" | "staging-us" | "staging-eu" | "prod-us" | "prod-eu"
	region:  "us" | "eu"
	type:    "prod" | "non-prod"
	version: "qa" | "staging" | "prod"
	chart:   =~"^[0-9]+\\.[0-9]+\\.[0-9]+$"
}

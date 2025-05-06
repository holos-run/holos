package holos

// WriteToDefault represents the default value of the --write-to flag which is
// the base directory where holos render component writes rendered configuration
// into.
const WriteToDefault string = "deploy"

// WriteToEnvVar represents the environment variable used to look up the default
// value of the --write-to flag.
const WriteToEnvVar string = "HOLOS_WRITE_TO"

// TypeMetaFile represents the file holos uses to discriminate the api version
// of a component BuildPlan.
const TypeMetaFile string = "typemeta.yaml"

// BuildPlanRequest represents the kind value of a BuildPlanRequest
const BuildPlanRequest string = "BuildPlanRequest"

// BuildPlanResponse represents the kind value of a BuildPlanResponse
const BuildPlanResponse string = "BuildPlanResponse"

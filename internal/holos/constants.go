package holos

// WriteToDefault represents the default value of the --write-to flag which is
// the base directory where holos render component writes rendered configuration
// into.
const WriteToDefault = "deploy"

// WriteToEnvVar represents the environment variable used to look up the default
// value of the --write-to flag.
const WriteToEnvVar string = "HOLOS_WRITE_TO"

// TypeMetaFile represents the file holos uses to discriminate the api version
// of a component BuildPlan.
const TypeMetaFile = "typemeta.yaml"

package v1alpha2

const (
	APIVersion    = "v1alpha2"
	BuildPlanKind = "BuildPlan"
	HelmChartKind = "HelmChart"
	// ChartDir is the directory name created in the holos component directory to cache a chart.
	ChartDir = "vendor"
	// ResourcesFile is the file name used to store component output when post-processing with kustomize.
	ResourcesFile = "resources.yaml"
)

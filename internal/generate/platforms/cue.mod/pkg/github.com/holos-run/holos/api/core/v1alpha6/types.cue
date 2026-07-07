package core

import "strings"

// A chart repository url is required unless the chart name is an oci://
// reference.  OCI charts pull directly from the registry and specify a
// repository only to configure authentication with the auth field.
#Chart: Chart={
	repository?: {
		if !strings.HasPrefix(Chart.name, "oci://") {
			url!: string
		}
	}
}

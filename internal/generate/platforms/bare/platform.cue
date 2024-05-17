package holos

import "encoding/json"

import v1 "github.com/holos-run/holos/api/v1alpha1"

// _Platform provides a platform resource to the holos cli for rendering.  The
// field is hidden because most components need to refer to platform data,
// specifically the platform model and the project models.  The platform
// resource itself is output once when rendering the entire platform, see the
// platform/ subdirectory.
_Platform: v1.#Platform & {
	metadata: name: string | *"bare" @tag(platform_name, type=string)

	// spec is the platform specification
	spec: {
		// model represents the web form values provided by the user.
		model: json.Unmarshal(_model)
		// _model is the json representation of model injected into CUE from holos.
		_model: string | *"{}" @tag(platform_model, type=string)
	}
}

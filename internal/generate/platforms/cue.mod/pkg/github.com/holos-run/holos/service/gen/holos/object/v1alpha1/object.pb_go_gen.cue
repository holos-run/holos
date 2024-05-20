package object

// Override the optional fields which result from the omitempty struct tags.
// Make them required so they work with yaml.Marshal

import "google.golang.org/protobuf/types/known/structpb"

// PlatformConfig represents the data passed from the holos cli to CUE when
// rendering configuration.
#PlatformConfig: {
	// Platform UUID.
	platform_id: string

	// Platform Model.
	platform_model: structpb.#Struct
}

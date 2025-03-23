package platform_test

import "embed"

//go:embed all:examples
var f embed.FS

// must align with embed all:examples directory
const examples string = "examples"

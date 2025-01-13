package main

import "holos.example/config/platform"

// Register all stack components with the platform spec.
for STACK in platform.stacks {
	Platform: Components: STACK.components
}

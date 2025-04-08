package policy

import apps "k8s.io/api/apps/v1"

// Organize by kind then name to avoid conflicts.
kind: [KIND=string]: [NAME=string]: {...}

// Block Secret resources. kind will not unify with "Secret"
kind: secret: [NAME=string]: kind: "Forbidden: Use an ExternalSecret instead: secret/\(NAME)"

// Validate Deployment against Kubernetes type definitions.
kind: deployment: [_]: apps.#Deployment

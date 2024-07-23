package holos

// This file contains specific to the k3d platform.

// Holos CLI client id is used as an audience in various AuthorizationPolicy resources
// so that curl -H "x-oidc-id-token:$(holos token)" works.
_HolosCLIClientID: "270319630705329162@holos_platform"

#AuthorizedUserAgent: string | *"anonymous"
_AuthorizedUserAgent: #AuthorizedUserAgent

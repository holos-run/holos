@if(!NoKargo)
package holos

import "example.com/platform/schemas/kargo"

// KargoProjects represents data to configure Kargo Project policies and
// promotion stages.  Similar in concept to the holos Projects structure.
KargoProjects: kargo.#Projects

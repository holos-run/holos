@extern(embed)
package holos

import "github.com/holos-run/holos/api/core/v1beta1:core"

holos: core.#Platform @embed(file=typemeta.yaml)

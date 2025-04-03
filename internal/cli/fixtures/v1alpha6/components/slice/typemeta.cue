@extern(embed)
package holos

import (
  "encoding/json"
  "github.com/holos-run/holos/api/core/v1alpha6:core"
)

_BuildContext: string | *"{}" @tag(holos_build_context, type=string)
BuildContext: core.#BuildContext & json.Unmarshal(_BuildContext)

holos: core.#BuildPlan & {
  buildContext: BuildContext
}

holos: _ @embed(file=typemeta.yaml)

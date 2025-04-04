@extern(embed)
package holos

import (
  "encoding/json"
  "github.com/holos-run/holos/api/core/v1alpha6:core"
)

holos: core.#BuildPlan & {
  _buildContext: string | *"{}" @tag(holos_build_context, type=string)
  buildContext: json.Unmarshal(_buildContext)
}

holos: _ @embed(file=typemeta.yaml)

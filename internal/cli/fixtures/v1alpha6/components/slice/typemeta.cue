@extern(embed)
package holos

import (
  "encoding/json"
  "github.com/holos-run/holos/api/core/v1alpha6:core"
)

holos: core.#BuildPlan & {
  _context: string | *"{}" @tag(context, type=string)
  context: json.Unmarshal(_context)
}

holos: _ @embed(file=typemeta.yaml)

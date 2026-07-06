@extern(embed)
package holos

import "encoding/json"

holos: _ @embed(file=typemeta.yaml)

holos: {
  _buildContext: string | *"{}" @tag(holos_build_context, type=string)
  buildContext: json.Unmarshal(_buildContext)
}

@extern(embed)
package holos

import "encoding/json"

holos: _ @embed(file=typemeta.yaml)

holos: {
	_context: string | *"{}" @tag(context, type=string)
	context:  json.Unmarshal(_context)
}

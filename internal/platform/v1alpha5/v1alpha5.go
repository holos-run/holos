package v1alpha5

import (
	"cuelang.org/go/cue"
	core "github.com/holos-run/holos/api/core/v1alpha5"
	component "github.com/holos-run/holos/internal/component/v1alpha5"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
)

// Platform represents a platform builder.
type Platform struct {
	Platform core.Platform
}

// Load loads from a cue value.
func (p *Platform) Load(v cue.Value) error {
	return errors.Wrap(v.Decode(&p.Platform))
}

func (p *Platform) Export(encoder holos.Encoder) error {
	if err := encoder.Encode(&p.Platform); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (p *Platform) Select(selectors ...holos.Selector) []holos.Component {
	components := make([]holos.Component, 0, len(p.Platform.Spec.Components))
	for _, com := range p.Platform.Spec.Components {
		if holos.IsSelected(com.Labels, selectors...) {
			components = append(components, &component.Component{Component: com})
		}
	}
	return components
}

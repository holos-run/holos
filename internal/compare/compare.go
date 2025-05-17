package compare

import (
	"github.com/holos-run/holos/internal/errors"
)

// Comparer handles comparison operations between BuildPlans
type Comparer struct {
}

// New creates a new Comparer instance
func New() *Comparer {
	return &Comparer{}
}

// BuildPlans compares two BuildPlan files for semantic equivalence
func (c *Comparer) BuildPlans(one, two string) error {
	return errors.NotImplemented()
}
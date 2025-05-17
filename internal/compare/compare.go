package compare

import (
	"bytes"
	"io"
	"os"

	"github.com/holos-run/holos/internal/errors"
	"gopkg.in/yaml.v3"
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
	// Read both files
	file1, err := os.Open(one)
	if err != nil {
		return errors.Format("opening first file: %w", err)
	}
	defer file1.Close()

	file2, err := os.Open(two)
	if err != nil {
		return errors.Format("opening second file: %w", err)
	}
	defer file2.Close()

	// Read all content from both files
	content1, err := io.ReadAll(file1)
	if err != nil {
		return errors.Format("reading first file: %w", err)
	}

	content2, err := io.ReadAll(file2)
	if err != nil {
		return errors.Format("reading second file: %w", err)
	}

	// Handle empty files case
	if len(content1) == 0 && len(content2) == 0 {
		return errors.NotImplemented()
	}

	// Parse YAML
	var bp1, bp2 map[string]interface{}
	
	if err := yaml.Unmarshal(content1, &bp1); err != nil {
		return errors.Format("parsing first file: %w", err)
	}
	
	if err := yaml.Unmarshal(content2, &bp2); err != nil {
		return errors.Format("parsing second file: %w", err)
	}

	// Compare the parsed structures
	return c.compareStructures(bp1, bp2)
}

// compareStructures compares two BuildPlan structures for semantic equivalence
func (c *Comparer) compareStructures(bp1, bp2 map[string]interface{}) error {
	// For now, implement basic equality check for the minimal test case
	if equalMaps(bp1, bp2) {
		return nil
	}
	
	return errors.New("BuildPlans are not semantically equivalent")
}

// equalMaps performs a basic deep equality check between two maps
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	
	for key, valA := range a {
		valB, ok := b[key]
		if !ok {
			return false
		}
		
		// Convert both values to comparable format
		yamlA, err := yaml.Marshal(valA)
		if err != nil {
			return false
		}
		
		yamlB, err := yaml.Marshal(valB)
		if err != nil {
			return false
		}
		
		if !bytes.Equal(yamlA, yamlB) {
			return false
		}
	}
	
	return true
}
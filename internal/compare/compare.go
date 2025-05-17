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

	// Parse YAML streams (multiple documents)
	docs1, err := parseYAMLStream(content1)
	if err != nil {
		return errors.Format("parsing first file: %w", err)
	}

	docs2, err := parseYAMLStream(content2)
	if err != nil {
		return errors.Format("parsing second file: %w", err)
	}

	// Compare the document lists
	return c.compareDocumentLists(docs1, docs2)
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

// parseYAMLStream parses a byte array containing one or more YAML documents
func parseYAMLStream(content []byte) ([]map[string]interface{}, error) {
	var documents []map[string]interface{}
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	
	for {
		var doc map[string]interface{}
		err := decoder.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if doc != nil {
			documents = append(documents, doc)
		}
	}
	
	return documents, nil
}

// compareDocumentLists compares two lists of YAML documents
func (c *Comparer) compareDocumentLists(docs1, docs2 []map[string]interface{}) error {
	if len(docs1) != len(docs2) {
		return errors.New("different number of documents")
	}
	
	for i := range docs1 {
		if err := c.compareStructures(docs1[i], docs2[i]); err != nil {
			return errors.Format("document %d: %w", i, err)
		}
	}
	
	return nil
}
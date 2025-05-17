package compare

import (
	"bytes"
	"io"
	"os"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	if cmp.Equal(bp1, bp2, cmpopts.EquateEmpty()) {
		return nil
	}
	
	return errors.New("BuildPlans are not semantically equivalent")
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
	
	// Convert to a sortable format for order-independent comparison
	less := func(a, b map[string]interface{}) bool {
		// Sort by a composite key based on available fields
		aKey := getCompositeKey(a)
		bKey := getCompositeKey(b)
		return aKey < bKey
	}
	
	// Use cmp with sort option for unordered comparison
	opts := []cmp.Option{
		cmpopts.SortSlices(less),
		cmpopts.EquateEmpty(),
	}
	
	if !cmp.Equal(docs1, docs2, opts...) {
		return errors.New("BuildPlans are not semantically equivalent")
	}
	
	return nil
}

// getCompositeKey creates a sortable key from a document
func getCompositeKey(doc map[string]interface{}) string {
	// Create a composite key based on common fields
	version, _ := doc["version"].(string)
	kind, _ := doc["kind"].(string)
	apiVersion, _ := doc["apiVersion"].(string)
	
	// If metadata exists, include name and labels
	name := ""
	labelsKey := ""
	if metadata, ok := doc["metadata"].(map[string]interface{}); ok {
		name, _ = metadata["name"].(string)
		
		// Include labels in the key for uniqueness
		if labels, ok := metadata["labels"].(map[string]interface{}); ok {
			// Sort label keys for consistent ordering
			labelKeys := make([]string, 0, len(labels))
			for k := range labels {
				labelKeys = append(labelKeys, k)
			}
			sort.Strings(labelKeys)
			
			// Build labels string
			for _, k := range labelKeys {
				v, _ := labels[k].(string)
				labelsKey += k + "=" + v + ","
			}
		}
	}
	
	return version + kind + apiVersion + name + labelsKey
}
package compare

import (
	"bytes"
	"io"
	"os"
	"sort"
	"strings"

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
	// Create comparison options for go-cmp
	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.Transformer("sortSlices", func(s []interface{}) []interface{} {
			return c.sortSlice(s)
		}),
	}

	// Deep order-independent comparison
	if cmp.Equal(bp1, bp2, opts...) {
		return nil
	}

	// Get the diff for the error message
	diff := cmp.Diff(bp1, bp2, opts...)

	// Extract specific field differences from the diff
	differences := c.extractFieldDifferences(diff)

	// Return error with the extracted differences
	if len(differences) > 0 {
		return errors.New(differences)
	}

	// Fallback to the full diff if no field differences found
	return errors.New(diff)
}

// sortSlice sorts a slice based on comparable string representation
func (c *Comparer) sortSlice(slice []interface{}) []interface{} {
	sorted := make([]interface{}, len(slice))
	copy(sorted, slice)

	sort.Slice(sorted, func(i, j int) bool {
		iStr := c.toComparableString(sorted[i])
		jStr := c.toComparableString(sorted[j])
		return iStr < jStr
	})

	return sorted
}

// toComparableString converts a value to a comparable string
func (c *Comparer) toComparableString(v interface{}) string {
	switch val := v.(type) {
	case map[string]interface{}:
		// Try to get identifying fields
		if artifact, ok := val["artifact"].(string); ok {
			return artifact
		}
		if name, ok := val["name"].(string); ok {
			return name
		}
		if metadata, ok := val["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				return name
			}
		}
		// Fallback to YAML representation
		yamlBytes, _ := yaml.Marshal(val)
		return string(yamlBytes)

	default:
		// Convert to YAML for comparison
		yamlBytes, _ := yaml.Marshal(v)
		return string(yamlBytes)
	}
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

	// Create a bipartite matching between documents
	used := make([]bool, len(docs2))

	// First pass: try to find exact matches
	for _, doc1 := range docs1 {
		for j, doc2 := range docs2 {
			if used[j] {
				continue
			}

			// Check if documents are exactly equal
			if c.documentsExactlyEqual(doc1, doc2) {
				used[j] = true
				break
			}
		}
	}

	// Second pass: handle unmatched documents
	usedIdx := 0
	for i, doc1 := range docs1 {
		// Find if this document was matched in first pass
		matchFound := false
		for j, doc2 := range docs2 {
			if used[j] && c.documentsExactlyEqual(doc1, doc2) {
				matchFound = true
				break
			}
		}

		if !matchFound {
			// Find the next unused document to compare against
			for usedIdx < len(docs2) && used[usedIdx] {
				usedIdx++
			}

			if usedIdx < len(docs2) {
				// Compare structures
				if err := c.compareStructures(doc1, docs2[usedIdx]); err != nil {
					return errors.Format("document %d: %w", i, err)
				}
				used[usedIdx] = true
			}
		}
	}

	return nil
}

// extractFieldDifferences extracts field-level differences from a go-cmp diff
func (c *Comparer) extractFieldDifferences(diff string) string {
	var differences []string
	lines := strings.Split(diff, "\n")

	for _, line := range lines {
		// Extract lines that show field differences
		if (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "+")) && strings.Contains(line, ":") {
			// Skip lines with formatting markers
			if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
				continue
			}

			// Handle the specific format from go-cmp
			trimmed := strings.TrimSpace(line)
			isRemoval := strings.HasPrefix(trimmed, "-")
			isAddition := strings.HasPrefix(trimmed, "+")

			// Remove the +/- prefix and extra spaces
			content := strings.TrimPrefix(trimmed, "-")
			content = strings.TrimPrefix(content, "+")
			content = strings.TrimSpace(content)

			// Look for key-value patterns
			if strings.Contains(content, ":") {
				// Handle quoted field names
				if strings.Contains(content, "\"") {
					// Extract field name between quotes
					start := strings.Index(content, "\"")
					end := strings.Index(content[start+1:], "\"")
					if start >= 0 && end >= 0 {
						fieldName := content[start+1 : start+1+end]

						// Extract value after colon
						colonIdx := strings.Index(content, ":")
						if colonIdx > 0 && colonIdx < len(content)-1 {
							value := strings.TrimSpace(content[colonIdx+1:])
							value = strings.TrimSuffix(value, ",")

							// Handle different value formats
							if strings.HasPrefix(value, "string(") {
								value = strings.TrimPrefix(value, "string(")
								value = strings.TrimSuffix(value, ")")
								value = strings.Trim(value, "\"")
							} else if strings.HasPrefix(value, "int(") {
								value = strings.TrimPrefix(value, "int(")
								value = strings.TrimSuffix(value, ")")
							} else if strings.HasPrefix(value, "\"") {
								value = strings.Trim(value, "\"")
							}

							// Format the difference
							if isRemoval {
								differences = append(differences, "-    "+fieldName+": "+value)
							} else if isAddition {
								differences = append(differences, "+    "+fieldName+": "+value)
							}
						}
					}
				}
			}
		}
	}

	// Return the differences as a single string
	return strings.Join(differences, "\n")
}

// documentsExactlyEqual checks if two documents are exactly equal
func (c *Comparer) documentsExactlyEqual(doc1, doc2 map[string]interface{}) bool {
	// Create comparison options for go-cmp
	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.Transformer("sortSlices", func(s []interface{}) []interface{} {
			return c.sortSlice(s)
		}),
	}

	return cmp.Equal(doc1, doc2, opts...)
}

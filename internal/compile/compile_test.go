package compile

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/holos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// typemetaCUE mirrors the embedded v1beta1 typemeta scaffolding from
// internal/component/platform/components/v1beta1/typemeta.cue.
const typemetaCUE = `@extern(embed)
package holos

import "encoding/json"

holos: _ @embed(file=typemeta.yaml)

holos: {
	_buildContext: string | *"{}" @tag(holos_build_context, type=string)
	buildContext: json.Unmarshal(_buildContext)
}
`

// componentCUE represents a minimal v1beta1 TaskSet component.
const componentCUE = `package holos

holos: {
	metadata: name: "example"
	spec: tasks: {
		resources: {
			kind:   "Resources"
			output: "example.gen.yaml"
			resources: ConfigMap: example: {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: name: "example"
			}
		}
		deploy: {
			kind: "Artifact"
			inputs: ["example.gen.yaml"]
		}
	}
}
`

// TestCompilerBeta1Envelope is a regression test proving a v1beta1 component
// compiles through the existing v1alpha6 BuildPlanRequest envelope.  The
// envelope discriminates the request version while the component's own version
// is re-discriminated from typemeta.yaml inside the read loop.
func TestCompilerBeta1Envelope(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"cue.mod/module.cue":               "module: \"holos.example\"\nlanguage: {\n\tversion: \"v0.12.0\"\n}\n",
		"components/example/typemeta.yaml": "apiVersion: v1beta1\nkind: TaskSet\n",
		"components/example/typemeta.cue":  typemetaCUE,
		"components/example/component.cue": componentCUE,
	}
	for path, content := range files {
		full := filepath.Join(root, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o777))
		require.NoError(t, os.WriteFile(full, []byte(content), 0o666))
	}

	req := BuildPlanRequest{
		APIVersion: "v1alpha6",
		Kind:       holos.BuildPlanRequest,
		Root:       root,
		Leaf:       "components/example",
		WriteTo:    holos.WriteToDefault,
		TempDir:    "${TMPDIR_PLACEHOLDER}",
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)

	var out bytes.Buffer
	c := New()
	c.R = bytes.NewReader(data)
	c.W = &out

	require.NoError(t, c.Run(t.Context()))

	var tm holos.TypeMeta
	require.NoError(t, json.Unmarshal(out.Bytes(), &tm))
	assert.Equal(t, "v1beta1", tm.APIVersion)
	assert.Equal(t, "TaskSet", tm.Kind)

	// The exported TaskSet must round trip the tasks.
	var taskSet struct {
		Spec struct {
			Tasks map[string]any `json:"tasks"`
		} `json:"spec"`
	}
	require.NoError(t, json.Unmarshal(out.Bytes(), &taskSet))
	assert.Contains(t, taskSet.Spec.Tasks, "resources")
	assert.Contains(t, taskSet.Spec.Tasks, "deploy")
}

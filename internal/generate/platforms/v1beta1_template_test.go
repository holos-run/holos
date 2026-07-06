package platforms

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"cuelang.org/go/cue"
	core "github.com/holos-run/holos/api/core/v1beta1"
	holoscue "github.com/holos-run/holos/internal/cue"
	"github.com/holos-run/holos/internal/generate"
)

// TestV1Beta1PlatformRegistered verifies directory-based registration: adding
// internal/generate/platforms/v1beta1 registers the platform with holos init
// platform with no Go code change.
func TestV1Beta1PlatformRegistered(t *testing.T) {
	if !slices.Contains(generate.Platforms(), "v1beta1") {
		t.Fatalf("want v1beta1 in generate.Platforms(), got %v", generate.Platforms())
	}
}

// TestV1Beta1InitPlatformTemplate generates the v1beta1 platform into a temp
// directory, adds a component written per the documented idiom (Component:
// #Kubernetes & {...}; holos: Component.TaskSet), and verifies the component
// exports a valid concrete core.#TaskSet assembled by the author layer.
func TestV1Beta1InitPlatformTemplate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	root := t.TempDir()
	if err := generate.GeneratePlatform(ctx, root, "v1beta1"); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}

	leaf := filepath.Join("components", "example")
	dir := filepath.Join(root, leaf)
	if err := os.MkdirAll(dir, 0o777); err != nil {
		t.Fatal(err)
	}

	typemetaCue := `@extern(embed)
package holos

import "encoding/json"

holos: _ @embed(file=typemeta.yaml)

holos: {
	_buildContext: string | *"{}" @tag(holos_build_context, type=string)
	buildContext: json.Unmarshal(_buildContext)
}
`
	typemetaYaml := "apiVersion: v1beta1\nkind: TaskSet\n"
	componentCue := `package holos

holos: Component.TaskSet

Component: #Kubernetes & {
	Resources: Namespace: example: metadata: name: "example"
}
`
	files := map[string]string{
		"typemeta.cue":  typemetaCue,
		"typemeta.yaml": typemetaYaml,
		"example.cue":   componentCue,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o666); err != nil {
			t.Fatal(err)
		}
	}

	tags := []string{
		"holos_component_name=example",
		"holos_component_path=" + leaf,
	}
	inst, err := holoscue.BuildInstance(root, leaf, tags)
	if err != nil {
		t.Fatalf("could not build instance: %v", err)
	}
	v, err := inst.HolosValue()
	if err != nil {
		t.Fatalf("could not get holos value: %v", err)
	}
	if err := v.Validate(cue.Concrete(true)); err != nil {
		t.Fatalf("holos value is not a valid concrete TaskSet: %v", err)
	}

	var ts core.TaskSet
	if err := v.Decode(&ts); err != nil {
		t.Fatalf("could not decode TaskSet: %v", err)
	}
	if ts.APIVersion != "v1beta1" {
		t.Errorf("want apiVersion v1beta1, got %s", ts.APIVersion)
	}
	if ts.Kind != "TaskSet" {
		t.Errorf("want kind TaskSet, got %s", ts.Kind)
	}
	if ts.Metadata.Name != "example" {
		t.Errorf("want metadata.name example, got %s", ts.Metadata.Name)
	}

	for _, name := range []string{"resources", "kustomize", "deploy"} {
		if _, ok := ts.Spec.Tasks[name]; !ok {
			t.Errorf("want task %s in spec.tasks, got %v", name, ts.Spec.Tasks)
		}
	}
	if got, want := ts.Spec.Tasks["deploy"].Artifact.Path, core.FileOrDirectoryPath("components/example/example.gen.yaml"); got != want {
		t.Errorf("want deploy artifact path %s, got %s", want, got)
	}
	if got, want := ts.Spec.Tasks["kustomize"].Output, core.FileOrDirectoryPath("example.gen.yaml"); got != want {
		t.Errorf("want kustomize output %s, got %s", want, got)
	}
}

package platforms

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// vetAuthorComponent validates one author component expression against the
// generated definitions in
// cue.mod/gen/github.com/holos-run/holos/api/author/v1beta1 by loading an
// overlay package from this directory, which provides the cue.mod module
// context.  The built value is returned so callers may assert on exported
// fields when validation succeeds.
func vetAuthorComponent(t *testing.T, component string) (cue.Value, error) {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src := fmt.Sprintf(`package authorprobe

import author "github.com/holos-run/holos/api/author/v1beta1:author"

component: %s
`, component)
	cfg := &load.Config{
		Dir: dir,
		Overlay: map[string]load.Source{
			filepath.Join(dir, "authorprobe", "probe.cue"): load.FromString(src),
		},
	}
	instances := load.Instances([]string{"./authorprobe"}, cfg)
	if len(instances) != 1 {
		t.Fatalf("want 1 instance, got %d", len(instances))
	}
	if instances[0].Err != nil {
		return cue.Value{}, instances[0].Err
	}
	value := cuecontext.New().BuildInstance(instances[0])
	if err := value.Err(); err != nil {
		return value, err
	}
	return value, value.Validate(cue.Concrete(true))
}

// TestV1Beta1AuthorDefinitions verifies the generated author v1beta1 CUE
// definitions resolve their imports of the core v1beta1 package and accept
// valid component wrapper values, including tasks mixed in through the
// ComponentConfig Tasks field.
func TestV1Beta1AuthorDefinitions(t *testing.T) {
	kubernetes := `author.#Kubernetes & {
		Name: "example"
		Path: "components/example"
		TaskSet: {
			metadata: name: "example"
			spec: tasks: {}
		}
	}`

	helm := `author.#Helm & {
		Name: "vault"
		Path: "components/vault"
		Chart: {
			name:    "vault"
			version: "1.0.0"
			release: "vault"
		}
		TaskSet: {
			metadata: name: "vault"
			spec: tasks: {}
		}
	}`

	kustomize := `author.#Kustomize & {
		Name: "example"
		Path: "components/example"
		KustomizeConfig: Files: "deployment.yaml": _
		TaskSet: {
			metadata: name: "example"
			spec: tasks: {}
		}
	}`

	// tasksMixin emulates one line of the hand-written assembly CUE provided
	// by the platform scaffolding: Tasks unify into TaskSet spec.tasks.
	tasksMixin := `author.#Kubernetes & {
		Name: "example"
		Path: "components/example"
		Tasks: gitops: {
			kind: "Resources"
			resources: Namespace: example: {
				apiVersion: "v1"
				kind:       "Namespace"
				metadata: name: "example"
			}
			output: "gitops/example.gen.yaml"
		}
		TaskSet: {
			metadata: name: "example"
			spec: tasks: Tasks
		}
	}`

	invalidTaskMixin := `author.#Kubernetes & {
		Name: "example"
		Path: "components/example"
		Tasks: broken: {
			kind: "Helm"
			command: args: ["true"]
			output: "broken.gen.yaml"
		}
		TaskSet: {
			metadata: name: "example"
			spec: tasks: {}
		}
	}`

	platform := `author.#Platform & {
		components: example: {
			name: "example"
			path: "components/example"
		}
		resource: {
			metadata: name: "default"
			spec: components: [for c in component.components {c}]
		}
	}`

	testCases := []struct {
		name      string
		component string
		wantErr   bool
		// assertPath and assertWant check one exported string field when set.
		assertPath string
		assertWant string
	}{
		{
			name:      "kubernetes component",
			component: kubernetes,
			wantErr:   false,
		},
		{
			name:      "helm component",
			component: helm,
			wantErr:   false,
		},
		{
			name:      "kustomize component",
			component: kustomize,
			wantErr:   false,
		},
		{
			name:       "kubernetes component with tasks mixin",
			component:  tasksMixin,
			wantErr:    false,
			assertPath: `component.TaskSet.spec.tasks.gitops.kind`,
			assertWant: "Resources",
		},
		{
			name:      "tasks mixin rejects mismatched kind config",
			component: invalidTaskMixin,
			wantErr:   true,
		},
		{
			name:       "platform with one component",
			component:  platform,
			wantErr:    false,
			assertPath: `component.resource.spec.components[0].name`,
			assertWant: "example",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := vetAuthorComponent(t, tc.component)
			if tc.wantErr && err == nil {
				t.Fatalf("want error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("want no error, got: %v", err)
			}
			if tc.assertPath != "" {
				got, err := value.LookupPath(cue.ParsePath(tc.assertPath)).String()
				if err != nil {
					t.Fatalf("lookup %s: %v", tc.assertPath, err)
				}
				if got != tc.assertWant {
					t.Fatalf("want %s at %s, got %s", tc.assertWant, tc.assertPath, got)
				}
			}
		})
	}
}

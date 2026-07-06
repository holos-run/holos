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

// vetTask validates one task expression against the hand authored per-kind
// guards in cue.mod/gen/github.com/holos-run/holos/api/core/v1beta1 by
// loading an overlay package from this directory, which provides the cue.mod
// module context.
func vetTask(t *testing.T, task string) error {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src := fmt.Sprintf(`package vetprobe

import core "github.com/holos-run/holos/api/core/v1beta1:core"

task: core.#Task & %s
`, task)
	cfg := &load.Config{
		Dir: dir,
		Overlay: map[string]load.Source{
			filepath.Join(dir, "vetprobe", "probe.cue"): load.FromString(src),
		},
	}
	instances := load.Instances([]string{"./vetprobe"}, cfg)
	if len(instances) != 1 {
		t.Fatalf("want 1 instance, got %d", len(instances))
	}
	if instances[0].Err != nil {
		return instances[0].Err
	}
	value := cuecontext.New().BuildInstance(instances[0])
	if err := value.Err(); err != nil {
		return err
	}
	return value.Validate(cue.Concrete(true))
}

func TestV1Beta1TaskConstraints(t *testing.T) {
	validHelm := `{
		kind: "Helm"
		helm: {
			chart: {name: "vault", version: "1.0.0", release: "vault"}
			values: {}
		}
		output: "vault.gen.yaml"
	}`

	testCases := []struct {
		name    string
		task    string
		wantErr bool
	}{
		{
			name:    "helm with config and output",
			task:    validHelm,
			wantErr: false,
		},
		{
			name: "resources with config and output",
			task: `{
				kind: "Resources"
				resources: Namespace: example: {
					apiVersion: "v1"
					kind:       "Namespace"
					metadata: name: "example"
				}
				output: "resources.gen.yaml"
			}`,
			wantErr: false,
		},
		{
			name: "file with config and output",
			task: `{
				kind: "File"
				file: source: "deployment.yaml"
				output: "deployment.gen.yaml"
			}`,
			wantErr: false,
		},
		{
			name: "kustomize with inputs and output",
			task: `{
				kind: "Kustomize"
				kustomize: kustomization: resources: ["a.yaml"]
				inputs: ["a.yaml"]
				output: "kustomized.gen.yaml"
			}`,
			wantErr: false,
		},
		{
			name: "join with inputs and output",
			task: `{
				kind: "Join"
				join: {}
				inputs: ["a.yaml", "b.yaml"]
				output: "joined.gen.yaml"
			}`,
			wantErr: false,
		},
		{
			name: "command validator with stdin input",
			task: `{
				kind: "Command"
				inputs: ["vault.gen.yaml"]
				command: {
					stdin: "vault.gen.yaml"
					args: ["holos", "cue", "vet", "-"]
				}
			}`,
			wantErr: false,
		},
		{
			name: "command generator capturing stdout",
			task: `{
				kind: "Command"
				command: {
					args: ["./read-thru-cache", "v0.16.0"]
					isStdoutOutput: true
				}
				output: "crds-bundle.yaml"
			}`,
			wantErr: false,
		},
		{
			name: "artifact with single input",
			task: `{
				kind: "Artifact"
				inputs: ["vault.gen.yaml"]
				artifact: path: "components/vault/vault.gen.yaml"
			}`,
			wantErr: false,
		},
		{
			name: "helm without helm config",
			task: `{
				kind:   "Helm"
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "helm with command config",
			task: `{
				kind: "Helm"
				helm: {
					chart: {name: "a", version: "1", release: "a"}
					values: {}
				}
				command: args: ["true"]
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "helm with inputs",
			task: `{
				kind: "Helm"
				helm: {
					chart: {name: "a", version: "1", release: "a"}
					values: {}
				}
				inputs: ["a.yaml"]
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "helm without output",
			task: `{
				kind: "Helm"
				helm: {
					chart: {name: "a", version: "1", release: "a"}
					values: {}
				}
			}`,
			wantErr: true,
		},
		{
			name: "kustomize without inputs",
			task: `{
				kind: "Kustomize"
				kustomize: kustomization: {}
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "join without inputs",
			task: `{
				kind: "Join"
				join: {}
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "command capturing stdout without output",
			task: `{
				kind: "Command"
				command: {
					args: ["./generate"]
					isStdoutOutput: true
				}
			}`,
			wantErr: true,
		},
		{
			name: "command with stdin not in inputs",
			task: `{
				kind: "Command"
				inputs: ["in.yaml"]
				command: {
					stdin: "other.yaml"
					args: ["cat"]
				}
			}`,
			wantErr: true,
		},
		{
			name: "command with stdin and no inputs",
			task: `{
				kind: "Command"
				command: {
					stdin: "in.yaml"
					args: ["cat"]
				}
			}`,
			wantErr: true,
		},
		{
			name: "command without args",
			task: `{
				kind: "Command"
				command: {}
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "command with empty args",
			task: `{
				kind: "Command"
				command: args: []
				output: "x.gen.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "artifact with two inputs",
			task: `{
				kind: "Artifact"
				inputs: ["a.yaml", "b.yaml"]
				artifact: path: "c"
			}`,
			wantErr: true,
		},
		{
			name: "artifact with output",
			task: `{
				kind: "Artifact"
				inputs: ["a.yaml"]
				artifact: path: "c"
				output: "d.yaml"
			}`,
			wantErr: true,
		},
		{
			name: "closed struct rejects dependson typo",
			task: `{
				kind: "Helm"
				helm: {
					chart: {name: "a", version: "1", release: "a"}
					values: {}
				}
				output: "x.gen.yaml"
				dependson: y: {}
			}`,
			wantErr: true,
		},
		{
			name: "dependsOn composes ordering edges",
			task: `{
				kind: "Artifact"
				dependsOn: validate: {}
				inputs: ["vault.gen.yaml"]
				artifact: path: "components/vault/vault.gen.yaml"
			}`,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := vetTask(t, tc.task)
			if tc.wantErr && err == nil {
				t.Errorf("want error, got nil for task: %s", tc.task)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("want nil, got error: %v for task: %s", err, tc.task)
			}
		})
	}
}

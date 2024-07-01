package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/holos-run/holos/api/core/v1alpha2"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/util"
)

// NewResult returns a new Result with the given holos component.
func NewResult(component v1alpha2.HolosComponent) *Result {
	return &Result{
		Kind:              "Result",
		APIVersion:        "v1alpha2",
		Component:         component,
		accumulatedOutput: "",
	}
}

// Result is the build result for display or writing.  Holos components Render
// the Result as a data pipeline.
type Result struct {
	// Kind is a string value representing the resource this object represents.
	Kind string `json:"kind" yaml:"kind" cue:"string | *\"Result\""`
	// APIVersion represents the versioned schema of this representation of an object.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1alpha2\""`

	// Component represents the common fields of all holos component kinds.
	Component v1alpha2.HolosComponent

	// accumulatedOutput accumulates rendered api objects.
	accumulatedOutput string
}

func (r *Result) GetAPIVersion() string {
	if r == nil {
		return ""
	}
	return r.APIVersion
}

func (r *Result) GetKind() string {
	if r == nil {
		return ""
	}
	return r.Kind
}

// Continue returns true if the result should be skipped over.
func (r *Result) Continue() bool {
	// Skip over a nil result
	if r == nil {
		return true
	}
	return r.Component.Skip
}

// Name returns the name of the component from the Metadata field.
func (r *Result) Name() string {
	if r == nil {
		return ""
	}
	return r.Component.Metadata.Name
}

// Filename returns the filename representing the rendered api objects of the Result.
func (r *Result) Filename(writeTo string, cluster string) string {
	name := r.Name()
	return filepath.Join(writeTo, "clusters", cluster, "components", name, name+".gen.yaml")
}

// KustomizationFilename returns the Flux Kustomization file path.
//
// Deprecated: Use DeployFiles instead.
func (r *Result) KustomizationFilename(writeTo string, cluster string) string {
	return filepath.Join(writeTo, "clusters", cluster, "holos", "components", r.Name()+"-kustomization.gen.yaml")
}

// AccumulatedOutput returns the accumulated rendered output.
func (r *Result) AccumulatedOutput() string {
	if r == nil {
		return ""
	}
	return r.accumulatedOutput
}

// addObjectMap renders the provided APIObjectMap into the accumulated output.
func (r *Result) addObjectMap(ctx context.Context, objectMap v1alpha2.APIObjectMap) {
	if r == nil {
		return
	}
	log := logger.FromContext(ctx)
	b := []byte(r.AccumulatedOutput())
	kinds := make([]v1alpha2.Kind, 0, len(objectMap))
	// Sort the keys
	for kind := range objectMap {
		kinds = append(kinds, kind)
	}
	slices.Sort(kinds)

	for _, kind := range kinds {
		v := objectMap[kind]
		// Sort the keys
		names := make([]v1alpha2.Label, 0, len(v))
		for name := range v {
			names = append(names, name)
		}
		slices.Sort(names)

		for _, name := range names {
			yamlString := v[name]
			log.Debug(fmt.Sprintf("%s/%s", kind, name), "kind", kind, "name", name)
			b = util.EnsureNewline(b)
			header := fmt.Sprintf("---\n# Source: CUE apiObjects.%s.%s\n", kind, name)
			b = append(b, []byte(header+yamlString)...)
			b = util.EnsureNewline(b)
		}
	}
	r.accumulatedOutput = string(b)
}

// kustomize replaces the accumulated output with the output of kustomize build
func (r *Result) kustomize(ctx context.Context) error {
	if r == nil {
		return nil
	}
	log := logger.FromContext(ctx)
	if r.Component.Kustomize.ResourcesFile == "" {
		log.DebugContext(ctx, "skipping kustomize: no resourcesFile")
		return nil
	}
	if len(r.Component.Kustomize.KustomizeFiles) < 1 {
		log.DebugContext(ctx, "skipping kustomize: no kustomizeFiles")
		return nil
	}
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)

	// Write the main api object resources file for kustomize.
	target := filepath.Join(tempDir, r.Component.Kustomize.ResourcesFile)
	b := []byte(r.AccumulatedOutput())
	b = util.EnsureNewline(b)
	if err := os.WriteFile(target, b, 0644); err != nil {
		return errors.Wrap(fmt.Errorf("could not write resources: %w", err))
	}
	log.DebugContext(ctx, "wrote: "+target, "op", "write", "path", target, "bytes", len(b))

	// Write the kustomization tree, kustomization.yaml must be in this map for kustomize to work.
	for file, content := range r.Component.Kustomize.KustomizeFiles {
		target := filepath.Join(tempDir, string(file))
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return errors.Wrap(err)
		}
		b := []byte(content)
		b = util.EnsureNewline(b)
		if err := os.WriteFile(target, b, 0644); err != nil {
			return errors.Wrap(fmt.Errorf("could not write: %w", err))
		}
		log.DebugContext(ctx, "wrote: "+target, "op", "write", "path", target, "bytes", len(b))
	}

	// Run kustomize.
	kOut, err := util.RunCmd(ctx, "kubectl", "kustomize", tempDir)
	if err != nil {
		log.ErrorContext(ctx, kOut.Stderr.String())
		return errors.Wrap(err)
	}
	// Replace the accumulated output
	r.accumulatedOutput = kOut.Stdout.String()
	return nil
}

func (r *Result) WriteDeployFiles(ctx context.Context, path string) error {
	if r == nil {
		return nil
	}
	log := logger.FromContext(ctx)
	if len(r.Component.DeployFiles) == 0 {
		return nil
	}
	for k, content := range r.Component.DeployFiles {
		path := filepath.Join(path, string(k))
		if err := r.Save(ctx, path, string(content)); err != nil {
			return errors.Wrap(err)
		}
		log.InfoContext(ctx, "wrote deploy file", "path", path, "bytes", len(content))
	}
	return nil
}

// Save writes the content to the filesystem for git ops.
func (r *Result) Save(ctx context.Context, path string, content string) error {
	log := logger.FromContext(ctx)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.FileMode(0775)); err != nil {
		log.WarnContext(ctx, "could not mkdir", "path", dir, "err", err)
		return errors.Wrap(err)
	}
	// Write the file content
	if err := os.WriteFile(path, []byte(content), os.FileMode(0644)); err != nil {
		log.WarnContext(ctx, "could not write", "path", path, "err", err)
		return errors.Wrap(err)
	}
	log.DebugContext(ctx, "out: wrote "+path, "action", "write", "path", path, "status", "ok")
	return nil
}

// SkipWriteAccumulatedOutput returns true if writing the accumulated output of
// k8s api objects should be skipped.  Useful for results which only write
// deployment files, like Flux or ArgoCD GitOps resources.
func (r *Result) SkipWriteAccumulatedOutput() bool {
	if r == nil {
		return true
	}
	// This is a hack and should be moved to a HolosComponent field or similar.
	if strings.HasPrefix(r.Component.Metadata.Name, "gitops/") {
		return true
	}
	return false
}

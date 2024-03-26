package v1alpha1

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/util"
	"github.com/holos-run/holos/pkg/wrapper"
)

// Result is the build result for display or writing.  Holos components Render the Result as a data pipeline.
type Result struct {
	HolosComponent
	// accumulatedOutput accumulates rendered api objects.
	accumulatedOutput string
}

func (r *Result) Name() string {
	return r.Metadata.Name
}

func (r *Result) Filename(writeTo string, cluster string) string {
	name := r.Metadata.Name
	return filepath.Join(writeTo, "clusters", cluster, "components", name, name+".gen.yaml")
}

func (r *Result) KustomizationFilename(writeTo string, cluster string) string {
	return filepath.Join(writeTo, "clusters", cluster, "holos", "components", r.Metadata.Name+"-kustomization.gen.yaml")
}

// AccumulatedOutput returns the accumulated rendered output.
func (r *Result) AccumulatedOutput() string {
	return r.accumulatedOutput
}

// addObjectMap renders the provided APIObjectMap into the accumulated output.
func (r *Result) addObjectMap(ctx context.Context, objectMap APIObjectMap) {
	log := logger.FromContext(ctx)
	b := []byte(r.AccumulatedOutput())
	kinds := make([]Kind, 0, len(objectMap))
	// Sort the keys
	for kind := range objectMap {
		kinds = append(kinds, kind)
	}
	slices.Sort(kinds)

	for _, kind := range kinds {
		v := objectMap[kind]
		// Sort the keys
		names := make([]Label, 0, len(v))
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
	log := logger.FromContext(ctx)
	if r.ResourcesFile == "" {
		log.DebugContext(ctx, "skipping kustomize: no resourcesFile")
		return nil
	}
	if len(r.KustomizeFiles) < 1 {
		log.DebugContext(ctx, "skipping kustomize: no kustomizeFiles")
		return nil
	}
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return wrapper.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)

	// Write the main api object resources file for kustomize.
	target := filepath.Join(tempDir, r.ResourcesFile)
	b := []byte(r.AccumulatedOutput())
	b = util.EnsureNewline(b)
	if err := os.WriteFile(target, b, 0644); err != nil {
		return wrapper.Wrap(fmt.Errorf("could not write resources: %w", err))
	}
	log.DebugContext(ctx, "wrote: "+target, "op", "write", "path", target, "bytes", len(b))

	// Write the kustomization tree, kustomization.yaml must be in this map for kustomize to work.
	for file, content := range r.KustomizeFiles {
		target := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return wrapper.Wrap(err)
		}
		b := []byte(content)
		b = util.EnsureNewline(b)
		if err := os.WriteFile(target, b, 0644); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not write: %w", err))
		}
		log.DebugContext(ctx, "wrote: "+target, "op", "write", "path", target, "bytes", len(b))
	}

	// Run kustomize.
	kOut, err := util.RunCmd(ctx, "kubectl", "kustomize", tempDir)
	if err != nil {
		log.ErrorContext(ctx, kOut.Stderr.String())
		return wrapper.Wrap(err)
	}
	// Replace the accumulated output
	r.accumulatedOutput = kOut.Stdout.String()
	return nil
}

// Save writes the content to the filesystem for git ops.
func (r *Result) Save(ctx context.Context, path string, content string) error {
	log := logger.FromContext(ctx)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.FileMode(0775)); err != nil {
		log.WarnContext(ctx, "could not mkdir", "path", dir, "err", err)
		return wrapper.Wrap(err)
	}
	// Write the kube api objects
	if err := os.WriteFile(path, []byte(content), os.FileMode(0644)); err != nil {
		log.WarnContext(ctx, "could not write", "path", path, "err", err)
		return wrapper.Wrap(err)
	}
	log.DebugContext(ctx, "out: wrote "+path, "action", "write", "path", path, "status", "ok")
	return nil
}

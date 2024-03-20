package v1alpha1

import (
	"context"
	"fmt"
	"github.com/holos-run/holos"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/util"
	"github.com/holos-run/holos/pkg/wrapper"
	"os"
	"path/filepath"
	"strings"
)

const (
	HelmChartKind = "HelmChart"
	// ChartDir is the directory name created in the holos component directory to cache a chart.
	ChartDir = "vendor"
)

// A HelmChart represents a helm command to provide chart values in order to render kubernetes api objects.
type HelmChart struct {
	HolosComponent `json:",inline" yaml:",inline"`
	// Namespace is the namespace to install into.  TODO: Use metadata.namespace instead.
	Namespace     string `json:"namespace"`
	Chart         Chart  `json:"chart"`
	ValuesContent string `json:"valuesContent"`
	EnableHooks   bool   `json:"enableHooks"`
}

type Chart struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Release    string     `json:"release"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (hc *HelmChart) Render(ctx context.Context, path holos.PathComponent) (*Result, error) {
	result := Result{
		TypeMeta:      hc.TypeMeta,
		Metadata:      hc.Metadata,
		Kustomization: hc.Kustomization,
	}
	if err := hc.helm(ctx, &result, path); err != nil {
		return nil, err
	}
	result.addObjectMap(ctx, hc.APIObjectMap)
	if err := result.kustomize(ctx); err != nil {
		return nil, wrapper.Wrap(fmt.Errorf("could not kustomize: %w", err))
	}
	return &result, nil
}

// runHelm provides the values produced by CUE to helm template and returns
// the rendered kubernetes api objects in the result.
func (hc *HelmChart) helm(ctx context.Context, r *Result, path holos.PathComponent) error {
	log := logger.FromContext(ctx).With("chart", hc.Chart.Name)
	if hc.Chart.Name == "" {
		log.WarnContext(ctx, "skipping helm: no chart name specified, use a different component type")
		return nil
	}

	cachedChartPath := filepath.Join(string(path), ChartDir, filepath.Base(hc.Chart.Name))
	if isNotExist(cachedChartPath) {
		// Add repositories
		repo := hc.Chart.Repository
		if repo.URL != "" {
			out, err := util.RunCmd(ctx, "helm", "repo", "add", repo.Name, repo.URL)
			if err != nil {
				log.ErrorContext(ctx, "could not run helm", "stderr", out.Stderr.String(), "stdout", out.Stdout.String())
				return wrapper.Wrap(fmt.Errorf("could not run helm repo add: %w", err))
			}
			// Update repository
			out, err = util.RunCmd(ctx, "helm", "repo", "update", repo.Name)
			if err != nil {
				log.ErrorContext(ctx, "could not run helm", "stderr", out.Stderr.String(), "stdout", out.Stdout.String())
				return wrapper.Wrap(fmt.Errorf("could not run helm repo update: %w", err))
			}
		} else {
			log.DebugContext(ctx, "no chart repository url proceeding assuming oci chart")
		}

		// Cache the chart
		if err := cacheChart(ctx, path, ChartDir, hc.Chart); err != nil {
			return fmt.Errorf("could not cache chart: %w", err)
		}
	}

	// Write values file
	tempDir, err := os.MkdirTemp("", "holos")
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not make temp dir: %w", err))
	}
	defer util.Remove(ctx, tempDir)

	valuesPath := filepath.Join(tempDir, "values.yaml")
	if err := os.WriteFile(valuesPath, []byte(hc.ValuesContent), 0644); err != nil {
		return wrapper.Wrap(fmt.Errorf("could not write values: %w", err))
	}
	log.DebugContext(ctx, "helm: wrote values", "path", valuesPath, "bytes", len(hc.ValuesContent))

	// Run charts
	chart := hc.Chart
	args := []string{"template"}
	if !hc.EnableHooks {
		args = append(args, "--no-hooks")
	}
	namespace := hc.Namespace
	args = append(args, "--include-crds", "--values", valuesPath, "--namespace", namespace, "--kubeconfig", "/dev/null", "--version", chart.Version, chart.Release, cachedChartPath)
	helmOut, err := util.RunCmd(ctx, "helm", args...)
	if err != nil {
		stderr := helmOut.Stderr.String()
		lines := strings.Split(stderr, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Error:") {
				err = fmt.Errorf("%s: %w", line, err)
			}
		}
		return wrapper.Wrap(fmt.Errorf("could not run helm template: %w", err))
	}

	r.accumulatedOutput = helmOut.Stdout.String()

	return nil
}

// cacheChart stores a cached copy of Chart in the chart subdirectory of path.
func cacheChart(ctx context.Context, path holos.PathComponent, chartDir string, chart Chart) error {
	log := logger.FromContext(ctx)

	cacheTemp, err := os.MkdirTemp(string(path), chartDir)
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not make temp dir: %w", err))
	}
	defer util.Remove(ctx, cacheTemp)

	chartName := chart.Name
	if chart.Repository.Name != "" {
		chartName = fmt.Sprintf("%s/%s", chart.Repository.Name, chart.Name)
	}
	helmOut, err := util.RunCmd(ctx, "helm", "pull", "--destination", cacheTemp, "--untar=true", "--version", chart.Version, chartName)
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not run helm pull: %w", err))
	}
	log.Debug("helm pull", "stdout", helmOut.Stdout, "stderr", helmOut.Stderr)

	cachePath := filepath.Join(string(path), chartDir)
	if err := os.Rename(cacheTemp, cachePath); err != nil {
		return wrapper.Wrap(fmt.Errorf("could not rename: %w", err))
	}
	log.InfoContext(ctx, "cached", "chart", chart.Name, "version", chart.Version, "path", cachePath)

	return nil
}
func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

package render

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/holos-run/holos"
	core "github.com/holos-run/holos/api/core/v1alpha3"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/util"
)

type HelmChart struct {
	Component core.HelmChart `json:"component"`
}

func (hc *HelmChart) Render(ctx context.Context, path holos.InstancePath) (*Result, error) {
	if hc == nil {
		return nil, nil
	}
	result := NewResult(hc.Component.Component)
	if err := hc.helm(ctx, result, path); err != nil {
		return nil, err
	}
	result.addObjectMap(ctx, hc.Component.APIObjectMap)
	if err := result.kustomize(ctx); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not kustomize: %w", err))
	}
	return result, nil
}

// runHelm provides the values produced by CUE to helm template and returns
// the rendered kubernetes api objects in the result.
func (hc *HelmChart) helm(ctx context.Context, r *Result, path holos.InstancePath) error {
	log := logger.FromContext(ctx).With("chart", hc.Component.Chart.Name)
	if hc.Component.Chart.Name == "" {
		log.WarnContext(ctx, "skipping helm: no chart name specified, use a different component type")
		return nil
	}

	cachedChartPath := filepath.Join(string(path), core.ChartDir, filepath.Base(hc.Component.Chart.Name))
	if isNotExist(cachedChartPath) {
		// Add repositories
		repo := hc.Component.Chart.Repository
		if repo.URL != "" {
			out, err := util.RunCmd(ctx, "helm", "repo", "add", repo.Name, repo.URL)
			if err != nil {
				log.ErrorContext(ctx, "could not run helm", "stderr", out.Stderr.String(), "stdout", out.Stdout.String())
				return errors.Wrap(fmt.Errorf("could not run helm repo add: %w", err))
			}
			// Update repository
			out, err = util.RunCmd(ctx, "helm", "repo", "update", repo.Name)
			if err != nil {
				log.ErrorContext(ctx, "could not run helm", "stderr", out.Stderr.String(), "stdout", out.Stdout.String())
				return errors.Wrap(fmt.Errorf("could not run helm repo update: %w", err))
			}
		} else {
			log.DebugContext(ctx, "no chart repository url proceeding assuming oci chart")
		}

		// Cache the chart
		if err := cacheChart(ctx, path, core.ChartDir, hc.Component.Chart); err != nil {
			return fmt.Errorf("could not cache chart: %w", err)
		}
	}

	// Write values file
	tempDir, err := os.MkdirTemp("", "holos")
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not make temp dir: %w", err))
	}
	defer util.Remove(ctx, tempDir)

	valuesPath := filepath.Join(tempDir, "values.yaml")
	if err := os.WriteFile(valuesPath, []byte(hc.Component.ValuesContent), 0644); err != nil {
		return errors.Wrap(fmt.Errorf("could not write values: %w", err))
	}
	log.DebugContext(ctx, "helm: wrote values", "path", valuesPath, "bytes", len(hc.Component.ValuesContent))

	// Run charts
	chart := hc.Component.Chart
	args := []string{"template"}
	if !hc.Component.EnableHooks {
		args = append(args, "--no-hooks")
	}
	namespace := hc.Component.Metadata.Namespace
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
		return errors.Wrap(fmt.Errorf("could not run helm template: %w", err))
	}

	r.accumulatedOutput = helmOut.Stdout.String()

	return nil
}

// cacheChart stores a cached copy of Chart in the chart subdirectory of path.
//
// It is assumed that the only method responsible for writing to chartDir is
// cacheChart itself.
//
// This relies on the atomicity of moving temporary directories into place on
// the same filesystem via os.Rename. If a syscall.EEXIST error occurs during
// renaming, it indicates that the cached chart already exists, which is an
// expected scenario when this function is called concurrently.
func cacheChart(ctx context.Context, path holos.InstancePath, chartDir string, chart core.Chart) error {
	log := logger.FromContext(ctx)

	cacheTemp, err := os.MkdirTemp(string(path), chartDir)
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not make temp dir: %w", err))
	}
	defer util.Remove(ctx, cacheTemp)

	chartName := chart.Name
	if chart.Repository.Name != "" {
		chartName = fmt.Sprintf("%s/%s", chart.Repository.Name, chart.Name)
	}
	helmOut, err := util.RunCmd(ctx, "helm", "pull", "--destination", cacheTemp, "--untar=true", "--version", chart.Version, chartName)
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not run helm pull: %w", err))
	}
	log.Debug("helm pull", "stdout", helmOut.Stdout, "stderr", helmOut.Stderr)

	cachePath := filepath.Join(string(path), chartDir)

	if err := os.MkdirAll(cachePath, 0777); err != nil {
		return errors.Wrap(fmt.Errorf("could not mkdir: %w", err))
	}

	items, err := os.ReadDir(cacheTemp)
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not read directory: %w", err))
	}

	for _, item := range items {
		src := filepath.Join(cacheTemp, item.Name())
		dst := filepath.Join(cachePath, item.Name())
		log.DebugContext(ctx, "rename", "src", src, "dst", dst)
		if err := os.Rename(src, dst); err != nil {
			var linkErr *os.LinkError
			if errors.As(err, &linkErr) && errors.Is(linkErr.Err, syscall.EEXIST) {
				log.DebugContext(ctx, "cache already exists", "chart", chart.Name, "chart_version", chart.Version, "path", cachePath)
			} else {
				return errors.Wrap(fmt.Errorf("could not rename: %w", err))
			}
		}
	}

	log.InfoContext(ctx, "cached", "chart", chart.Name, "chart_version", chart.Version, "path", cachePath)

	return nil
}
func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

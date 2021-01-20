package helmplugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/artifacthub/hub/internal/hub"
	"github.com/artifacthub/hub/internal/license"
	"github.com/artifacthub/hub/internal/pkg"
	"helm.sh/helm/v3/pkg/plugin"
	"sigs.k8s.io/yaml"
)

var (
	// licenseRE is a regular expression used to locate a license file in the
	// repository.
	licenseRE = regexp.MustCompile(`(?i)license.*`)
)

// TrackerSource is a hub.TrackerSource implementation for Helm plugins
// repositories.
type TrackerSource struct {
	i *hub.TrackerSourceInput
}

// NewTrackerSource creates a new TrackerSource instance.
func NewTrackerSource(i *hub.TrackerSourceInput) *TrackerSource {
	return &TrackerSource{i}
}

// GetPackagesAvailable implements the TrackerSource interface.
func (s *TrackerSource) GetPackagesAvailable() (map[string]*hub.Package, error) {
	packagesAvailable := make(map[string]*hub.Package)

	// Walk the path provided looking for available packages
	err := filepath.Walk(s.i.BasePath, func(pkgPath string, info os.FileInfo, err error) error {
		// Return ASAP if context is cancelled
		select {
		case <-s.i.Svc.Ctx.Done():
			return s.i.Svc.Ctx.Err()
		default:
		}

		// If an error is raised while visiting a path or the path is not a
		// directory, we skip it
		if err != nil || !info.IsDir() {
			return nil
		}

		// Read and parse plugin metadata file
		data, err := ioutil.ReadFile(filepath.Join(pkgPath, plugin.PluginFileName))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				s.warn(fmt.Errorf("error reading plugin metadata file: %w", err))
			}
			return nil
		}
		var md *plugin.Metadata
		if err = yaml.Unmarshal(data, &md); err != nil || md == nil {
			s.warn(fmt.Errorf("error unmarshaling plugin metadata file: %w", err))
			return nil
		}

		// Prepare and store package version
		p := preparePackage(s.i.Repository, md, pkgPath)
		packagesAvailable[pkg.BuildKey(p)] = p

		return nil
	})
	if err != nil {
		return nil, err
	}

	return packagesAvailable, nil
}

// warn is a helper that sends the error provided to the errors collector and
// logs it as a warning.
func (s *TrackerSource) warn(err error) {
	s.i.Svc.Logger.Warn().Err(err).Send()
	s.i.Svc.Ec.Append(s.i.Repository.RepositoryID, err)
}

// preparePackage prepares a package version using the plugin metadata and the
// files in the path provided.
func preparePackage(r *hub.Repository, md *plugin.Metadata, pkgPath string) *hub.Package {
	// Prepare package from metadata
	p := &hub.Package{
		Name:        md.Name,
		Version:     md.Version,
		Description: md.Description,
		Keywords: []string{
			"helm",
			"helm-plugin",
		},
		Links: []*hub.Link{
			{
				Name: "Source",
				URL:  r.URL,
			},
		},
		Repository: r,
	}

	// Include readme file if available
	readme, err := ioutil.ReadFile(filepath.Join(pkgPath, "README.md"))
	if err == nil {
		p.Readme = string(readme)
	}

	// Process and include license if available
	files, err := ioutil.ReadDir(pkgPath)
	if err == nil {
		for _, file := range files {
			if licenseRE.Match([]byte(file.Name())) {
				licenseFile, err := ioutil.ReadFile(filepath.Join(pkgPath, file.Name()))
				if err == nil {
					p.License = license.Detect(licenseFile)
					break
				}
			}
		}
	}

	return p
}

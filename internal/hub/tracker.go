package hub

import (
	"context"
	"net/http"

	"github.com/artifacthub/hub/internal/img"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// ErrorsCollector interface defines the methods that an errors collector
// implementation should provide.
type ErrorsCollector interface {
	Append(repositoryID string, err error)
	Flush()
}

// HTTPClient defines the methods an HTTPClient implementation must provide.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TrackerServices represents a set of services that must be provided to a
// Tracker instance so that it can perform its tasks.
type TrackerServices struct {
	Ctx      context.Context
	Cfg      *viper.Viper
	Rm       RepositoryManager
	Pm       PackageManager
	Rc       RepositoryCloner
	Oe       OLMOCIExporter
	Ec       ErrorsCollector
	Is       img.Store
	GithubRL *rate.Limiter
}

// TrackerSource defines the methods a TrackerSource implementation must
// provide.
type TrackerSource interface {
	// GetPackagesAvailable represents a function that returns a list of
	// available packages in a given repository. Each repository kind will
	// require using a specific TrackerSource implementation that will know
	// best how to get the available packages in the repository. The key used
	// in the returned map is expected to be built using the BuildKey helper
	// function in the pkg package.
	GetPackagesAvailable() (map[string]*Package, error)
}

// TrackerSourceInput represents the input provided to a TrackerSource to get
// the packages available in a repository when tracking it.
type TrackerSourceInput struct {
	Ctx                context.Context
	Cfg                *viper.Viper
	Repository         *Repository
	PackagesRegistered map[string]string
	BasePath           string
	Logger             zerolog.Logger
	Ec                 ErrorsCollector
}

package tracker

import (
	"context"
	"fmt"

	"github.com/artifacthub/hub/internal/hub"
	"github.com/spf13/viper"
)

// GetRepositories gets the repositories the tracker will process based on the
// configuration provided:
//
// - If a list of repositories names, those will be the repositories returned
//   provided they are found.
// - If a list of repositories kinds is provided, all repositories of those
//   kinds will be returned.
// - Otherwise, all the repositories will be returned.
//
// NOTE: disabled repositories will be filtered out.
func GetRepositories(
	ctx context.Context,
	cfg *viper.Viper,
	rm hub.RepositoryManager,
) ([]*hub.Repository, error) {
	reposNames := cfg.GetStringSlice("tracker.repositoriesNames")
	reposKinds := cfg.GetStringSlice("tracker.repositoriesKinds")

	var repos []*hub.Repository
	switch {
	case len(reposNames) > 0:
		for _, name := range reposNames {
			repo, err := rm.GetByName(ctx, name, true)
			if err != nil {
				return nil, fmt.Errorf("error getting repository %s: %w", name, err)
			}
			repos = append(repos, repo)
		}
	case len(reposKinds) > 0:
		for _, kindName := range reposKinds {
			kind, err := hub.GetKindFromName(kindName)
			if err != nil {
				return nil, fmt.Errorf("invalid repository kind found in config: %s", kindName)
			}
			kindRepos, err := rm.GetByKind(ctx, kind, true)
			if err != nil {
				return nil, fmt.Errorf("error getting repositories by kind (%s): %w", kindName, err)
			}
			repos = append(repos, kindRepos...)
		}
	default:
		var err error
		repos, err = rm.GetAll(ctx, true)
		if err != nil {
			return nil, fmt.Errorf("error getting all repositories: %w", err)
		}
	}

	// Filter out disabled repositories
	var reposFiltered []*hub.Repository
	for _, repo := range repos {
		if !repo.Disabled {
			reposFiltered = append(reposFiltered, repo)
		}
	}

	return reposFiltered, nil
}

package tracker

import (
	"context"
	"errors"
	"testing"

	"github.com/artifacthub/hub/internal/hub"
	"github.com/artifacthub/hub/internal/repo"
	"github.com/artifacthub/hub/internal/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetRepositories(t *testing.T) {
	ctx := context.Background()
	repo1 := &hub.Repository{
		Name: "repo1",
		Kind: hub.Helm,
	}
	repo2 := &hub.Repository{
		Name: "repo2",
		Kind: hub.OLM,
	}
	repo3 := &hub.Repository{
		Name:     "repo3",
		Kind:     hub.OPA,
		Disabled: true,
	}

	t.Run("error getting repository by name", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetByName", ctx, "repo1", true).Return(nil, tests.ErrFake)

		// Run test and check expectations
		cfg := viper.New()
		cfg.Set("tracker.repositoriesNames", []string{"repo1"})
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.True(t, errors.Is(err, tests.ErrFake))
		assert.Nil(t, repos)
		rm.AssertExpectations(t)
	})

	t.Run("get repositories by name", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetByName", ctx, "repo1", true).Return(repo1, nil)
		rm.On("GetByName", ctx, "repo2", true).Return(repo2, nil)

		// Run test and check expectations
		cfg := viper.New()
		cfg.Set("tracker.repositoriesNames", []string{"repo1", "repo2"})
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []*hub.Repository{repo1, repo2}, repos)
		rm.AssertExpectations(t)
	})

	t.Run("error getting kind from name", func(t *testing.T) {
		t.Parallel()
		cfg := viper.New()
		cfg.Set("tracker.repositoriesKinds", []string{"invalid"})
		repos, err := GetRepositories(ctx, cfg, nil)
		assert.Error(t, err)
		assert.Nil(t, repos)
	})

	t.Run("error getting repository by kind", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetByKind", ctx, hub.Helm, true).Return(nil, tests.ErrFake)

		// Run test and check expectations
		cfg := viper.New()
		cfg.Set("tracker.repositoriesKinds", []string{"helm"})
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.True(t, errors.Is(err, tests.ErrFake))
		assert.Nil(t, repos)
		rm.AssertExpectations(t)
	})

	t.Run("get repositories by kind", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetByKind", ctx, hub.Helm, true).Return([]*hub.Repository{repo1}, nil)
		rm.On("GetByKind", ctx, hub.OLM, true).Return([]*hub.Repository{repo2}, nil)

		// Run test and check expectations
		cfg := viper.New()
		cfg.Set("tracker.repositoriesKinds", []string{"helm", "olm"})
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []*hub.Repository{repo1, repo2}, repos)
		rm.AssertExpectations(t)
	})

	t.Run("names have preference over kinds when both are provided", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetByName", ctx, "repo1", true).Return(repo1, nil)
		rm.On("GetByName", ctx, "repo2", true).Return(repo2, nil)

		// Run test and check expectations
		cfg := viper.New()
		cfg.Set("tracker.repositoriesNames", []string{"repo1", "repo2"})
		cfg.Set("tracker.repositoriesKinds", []string{"helm", "olm"})
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []*hub.Repository{repo1, repo2}, repos)
		rm.AssertExpectations(t)
	})

	t.Run("error getting all repositories", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetAll", ctx, true).Return(nil, tests.ErrFake)

		// Run test and check expectations
		cfg := viper.New()
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.True(t, errors.Is(err, tests.ErrFake))
		assert.Nil(t, repos)
		rm.AssertExpectations(t)
	})

	t.Run("get all repositories", func(t *testing.T) {
		t.Parallel()

		// Setup expectations
		rm := &repo.ManagerMock{}
		rm.On("GetAll", ctx, true).Return([]*hub.Repository{repo1, repo2, repo3}, nil)

		// Run test and check expectations
		cfg := viper.New()
		repos, err := GetRepositories(ctx, cfg, rm)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []*hub.Repository{repo1, repo2}, repos)
		rm.AssertExpectations(t)
	})
}

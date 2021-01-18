package repo

import (
	"context"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/artifacthub/hub/internal/hub"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// OCITagsGetter provides a mechanism to get all the tags available for a given
// repository in a OCI registry.
type OCITagsGetter struct{}

// Tags returns a list with the tags available for the provided repository.
func (tg *OCITagsGetter) Tags(ctx context.Context, r *hub.Repository) ([]string, error) {
	u := strings.TrimPrefix(r.URL, hub.RepositoryOCIPrefix)
	ociRepo, err := name.NewRepository(u)
	if err != nil {
		return nil, err
	}
	var options []remote.Option
	if r.AuthUser != "" || r.AuthPass != "" {
		options = []remote.Option{
			remote.WithAuth(&authn.Basic{
				Username: r.AuthUser,
				Password: r.AuthPass,
			}),
		}
	}
	tags, err := remote.ListWithContext(ctx, ociRepo, options...)
	if err != nil {
		return nil, err
	}
	sort.Slice(tags, func(i, j int) bool {
		vi, _ := semver.NewVersion(tags[i])
		vj, _ := semver.NewVersion(tags[j])
		return vj.LessThan(vi)
	})
	return tags, nil
}

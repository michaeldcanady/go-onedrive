package infra

import (
	"os"
	"path/filepath"

	domaincache "github.com/michaeldcanady/go-onedrive/internal/cache/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal/profile/infra"
)

type Repository struct {
	path       string
	serializer domaincache.SerializerDeserializer[domainstate.State]
}

func NewRepository(
	path string,
	serializer domaincache.SerializerDeserializer[domainstate.State],
) *Repository {
	return &Repository{
		path:       path,
		serializer: serializer,
	}
}

func (r *Repository) ensureRoot() error {
	dir := filepath.Dir(r.path)
	return os.MkdirAll(dir, 0o700)
}

func (r *Repository) Load() (domainstate.State, error) {
	if err := r.ensureRoot(); err != nil {
		return domainstate.State{}, err
	}

	b, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return r.defaultState(), nil
	}
	if err != nil {
		return domainstate.State{}, err
	}

	return r.serializer.Deserialize(b)
}

func (r *Repository) Save(state domainstate.State) error {
	if err := r.ensureRoot(); err != nil {
		return err
	}

	b, err := r.serializer.Serialize(state)
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, b, 0o600)
}

func (r *Repository) defaultState() domainstate.State {

	return domainstate.State{
		CurrentProfile: infraprofile.DefaultProfileName,
		CurrentDrive:   "",
	}
}

package state

import (
	"context"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
)

type Repository struct {
	path         string
	serializer   abstractions.SerializerDeserializer[state.State]
	driveService drive.DriveService
}

func NewRepository(
	path string,
	serializer abstractions.SerializerDeserializer[state.State],
	driveService drive.DriveService,
) *Repository {
	return &Repository{
		path:         path,
		serializer:   serializer,
		driveService: driveService,
	}
}

func (r *Repository) ensureRoot() error {
	dir := filepath.Dir(r.path)
	return os.MkdirAll(dir, 0o700)
}

func (r *Repository) Load() (state.State, error) {
	if err := r.ensureRoot(); err != nil {
		return state.State{}, err
	}

	b, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return r.defaultState(), nil
	}
	if err != nil {
		return state.State{}, err
	}

	return r.serializer.Deserialize(b)
}

func (r *Repository) Save(state state.State) error {
	if err := r.ensureRoot(); err != nil {
		return err
	}

	b, err := r.serializer.Serialize(state)
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, b, 0o600)
}

func (r *Repository) defaultState() state.State {

	drive, _ := r.driveService.ResolvePersonalDrive(context.Background())
	driveID := ""
	if drive != nil {
		driveID = drive.ID
	}

	return state.State{
		CurrentProfile: profile.DefaultProfileName,
		CurrentDrive:   driveID,
	}
}

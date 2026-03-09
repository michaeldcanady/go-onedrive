package state

import (
	domainstate "github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
)

var _ domainstate.Service = (*Service)(nil)

type Service struct {
	repo            domainstate.Repository
	sessionOverride string
	hasSession      bool
}

func NewService(repo domainstate.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) getState() (domainstate.State, error) {
	return s.repo.Load()
}

func (s *Service) SetSessionProfile(name string) {
	s.sessionOverride = name
	s.hasSession = true
}

func (s *Service) GetCurrentProfile() (string, error) {
	// Session override takes precedence
	if s.hasSession {
		return s.sessionOverride, nil
	}

	// Otherwise use persistent state
	st, err := s.getState()
	if err != nil {
		return "", err
	}

	return st.CurrentProfile, nil
}

func (s *Service) SetCurrentProfile(name string) error {
	st, err := s.getState()
	if err != nil {
		return err
	}

	st.CurrentProfile = name
	return s.repo.Save(st)
}

func (s *Service) ClearCurrentProfile() error {
	st, err := s.getState()
	if err != nil {
		return err
	}

	st.CurrentProfile = profile.DefaultProfileName
	return s.repo.Save(st)
}

func (s *Service) GetCurrentDrive() (string, error) {
	// TODO: Add support for session state

	st, err := s.getState()
	if err != nil {
		return "", err
	}

	return st.CurrentDrive, nil
}

func (s *Service) SetCurrentDrive(name string) error {
	st, err := s.getState()
	if err != nil {
		return err
	}

	st.CurrentDrive = name
	return s.repo.Save(st)
}

func (s *Service) ClearCurrentDrive() error {
	st, err := s.getState()
	if err != nil {
		return err
	}

	// TODO: Figure out how to default to current user's personal drive.
	st.CurrentDrive = ""
	return s.repo.Save(st)
}

func (s *Service) GetDriveAlias(alias string) (string, error) {
	st, err := s.getState()
	if err != nil {
		return "", err
	}

	if st.DriveAliases == nil {
		return "", nil
	}

	return st.DriveAliases[alias], nil
}

func (s *Service) SetDriveAlias(alias, driveID string) error {
	st, err := s.getState()
	if err != nil {
		return err
	}

	if st.DriveAliases == nil {
		st.DriveAliases = make(map[string]string)
	}

	st.DriveAliases[alias] = driveID
	return s.repo.Save(st)
}

func (s *Service) RemoveDriveAlias(alias string) error {
	st, err := s.getState()
	if err != nil {
		return err
	}

	if st.DriveAliases == nil {
		return nil
	}

	delete(st.DriveAliases, alias)
	return s.repo.Save(st)
}

func (s *Service) ListDriveAliases() (map[string]string, error) {
	st, err := s.getState()
	if err != nil {
		return nil, err
	}

	return st.DriveAliases, nil
}

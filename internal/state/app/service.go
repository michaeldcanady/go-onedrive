package app

import (
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal/profile/infra"
)

var _ domainstate.Service = (*Service)(nil)

type Service struct {
	repo                 domainstate.Repository
	sessionProfileOverride string
	hasProfileSession      bool
	sessionDriveOverride   string
	hasDriveSession        bool
}

func NewService(repo domainstate.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) getState() (domainstate.State, error) {
	return s.repo.Load()
}

func (s *Service) GetCurrentProfile() (string, error) {
	if s.hasProfileSession {
		return s.sessionProfileOverride, nil
	}

	st, err := s.getState()
	if err != nil {
		return "", err
	}

	return st.CurrentProfile, nil
}

func (s *Service) SetCurrentProfile(name string, scope domainstate.Scope) error {
	if scope == domainstate.ScopeSession {
		s.sessionProfileOverride = name
		s.hasProfileSession = true
		return nil
	}

	st, err := s.getState()
	if err != nil {
		return err
	}

	st.CurrentProfile = name
	return s.repo.Save(st)
}

func (s *Service) ClearCurrentProfile() error {
	s.hasProfileSession = false
	
	st, err := s.getState()
	if err != nil {
		return err
	}

	st.CurrentProfile = infraprofile.DefaultProfileName
	return s.repo.Save(st)
}

func (s *Service) GetCurrentDrive() (string, error) {
	if s.hasDriveSession {
		return s.sessionDriveOverride, nil
	}

	st, err := s.getState()
	if err != nil {
		return "", err
	}

	return st.CurrentDrive, nil
}

func (s *Service) SetCurrentDrive(id string, scope domainstate.Scope) error {
	if scope == domainstate.ScopeSession {
		s.sessionDriveOverride = id
		s.hasDriveSession = true
		return nil
	}

	st, err := s.getState()
	if err != nil {
		return err
	}

	st.CurrentDrive = id
	return s.repo.Save(st)
}

func (s *Service) ClearCurrentDrive() error {
	s.hasDriveSession = false

	st, err := s.getState()
	if err != nil {
		return err
	}

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

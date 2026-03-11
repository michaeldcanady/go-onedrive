package app

import (
	"fmt"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal/profile/infra"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

var _ domainstate.Service = (*Service)(nil)

type Service struct {
	repo                   domainstate.Repository
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

func (s *Service) Get(key domainstate.Key) (string, error) {
	switch key {
	case domainstate.KeyProfile:
		if s.hasProfileSession {
			return s.sessionProfileOverride, nil
		}

		st, err := s.getState()
		if err != nil {
			return "", err
		}

		return st.CurrentProfile, nil
	case domainstate.KeyDrive:
		if s.hasDriveSession {
			return s.sessionDriveOverride, nil
		}

		st, err := s.getState()
		if err != nil {
			return "", err
		}

		return st.CurrentDrive, nil
	default:
		return "", fmt.Errorf("unsupported state key: %v", key)
	}
}

func (s *Service) Set(key domainstate.Key, value string, scope domainstate.Scope) error {
	switch key {
	case domainstate.KeyProfile:
		if scope == domainstate.ScopeSession {
			s.sessionProfileOverride = value
			s.hasProfileSession = true
			return nil
		}

		st, err := s.getState()
		if err != nil {
			return err
		}

		st.CurrentProfile = value
		return s.repo.Save(st)
	case domainstate.KeyDrive:
		if scope == domainstate.ScopeSession {
			s.sessionDriveOverride = value
			s.hasDriveSession = true
			return nil
		}

		st, err := s.getState()
		if err != nil {
			return err
		}

		st.CurrentDrive = value
		return s.repo.Save(st)
	default:
		return fmt.Errorf("unsupported state key: %v", key)
	}
}

func (s *Service) Clear(key domainstate.Key) error {
	switch key {
	case domainstate.KeyProfile:
		s.hasProfileSession = false

		st, err := s.getState()
		if err != nil {
			return err
		}

		st.CurrentProfile = infraprofile.DefaultProfileName
		return s.repo.Save(st)
	case domainstate.KeyDrive:
		s.hasDriveSession = false

		st, err := s.getState()
		if err != nil {
			return err
		}

		st.CurrentDrive = ""
		return s.repo.Save(st)
	default:
		return fmt.Errorf("unsupported state key: %v", key)
	}
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

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

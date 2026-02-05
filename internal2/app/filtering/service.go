package filtering

import (
	"errors"
	"fmt"
	"reflect"

	domainfiltering "github.com/michaeldcanady/go-onedrive/internal2/domain/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

var _ domainfiltering.Service = (*Service)(nil)

type Service struct {
	registry map[string]domainfiltering.Filter
	byType   map[string]map[reflect.Type]domainfiltering.Filter
	logger   logging.Logger
}

func NewService(logger logging.Logger) *Service {
	return &Service{
		registry: map[string]domainfiltering.Filter{},
		byType:   map[string]map[reflect.Type]domainfiltering.Filter{},
		logger:   logger,
	}
}

// resolve chooses the best filter for the given filter name + value
func (s *Service) resolve(filter string, v any) (domainfiltering.Filter, error) {
	t := reflect.TypeOf(v)
	origType := t

	// If it's a slice, match on element type
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	s.logger.Debug("resolving filter",
		logging.String("filter", filter),
		logging.String("valueType", origType.String()),
		logging.String("matchType", t.String()),
	)

	// First try type-specific
	if typeMap, ok := s.byType[filter]; ok {
		if f, ok := typeMap[t]; ok {
			s.logger.Debug("resolved type-specific filter",
				logging.String("filter", filter),
				logging.String("type", t.String()),
			)
			return f, nil
		}

		s.logger.Debug("no type-specific filter match",
			logging.String("filter", filter),
			logging.String("type", t.String()),
		)
	}

	// Fallback to generic
	if f, ok := s.registry[filter]; ok {
		s.logger.Debug("resolved generic filter",
			logging.String("filter", filter),
		)
		return f, nil
	}

	s.logger.Error("failed to resolve filter",
		logging.String("filter", filter),
		logging.String("valueType", origType.String()),
	)

	return nil, fmt.Errorf("no filter registered for filter=%q and type=%T", filter, v)
}

// Filter implements [filtering.Service].
func (s *Service) Filter(filter string, v any) error {
	s.logger.Debug("filtering started",
		logging.String("filter", filter),
		logging.Any("valueType", reflect.TypeOf(v).String()),
	)

	filterer, err := s.resolve(filter, v)
	if err != nil {
		s.logger.Error("filter resolution failed",
			logging.String("filter", filter),
			logging.Error(err),
		)
		return err
	}

	if err := filterer.Filter(v); err != nil {
		s.logger.Error("filter execution failed",
			logging.String("filter", filter),
			logging.Error(err),
		)
		return err
	}

	s.logger.Info("filter applied successfully",
		logging.String("filter", filter),
		logging.Any("valueType", reflect.TypeOf(v).String()),
	)

	return nil
}

// RegisterWithType registers a filter for a specific Go type under a filter name
func (s *Service) RegisterWithType(filter string, typ reflect.Type, f domainfiltering.Filter) error {
	if s.byType == nil {
		return errors.New("registry is nil")
	}

	origType := typ
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	s.logger.Info("registering type-specific filter",
		logging.String("filter", filter),
		logging.String("type", typ.String()),
		logging.String("originalType", origType.String()),
	)

	if _, ok := s.byType[filter]; !ok {
		s.byType[filter] = make(map[reflect.Type]domainfiltering.Filter)
	}
	s.byType[filter][typ] = f

	return nil
}

// Register implements [filtering.Service].
func (s *Service) Register(filter string, f domainfiltering.Filter) error {
	if s.registry == nil {
		return errors.New("registry is nil")
	}

	s.logger.Info("registering generic filter",
		logging.String("filter", filter),
	)

	s.registry[filter] = f
	return nil
}

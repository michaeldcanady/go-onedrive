package filtering

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	domainfiltering "github.com/michaeldcanady/go-onedrive/internal2/domain/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
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
	cid := util.CorrelationIDFromContext(context.Background()) // fallback if no ctx provided

	t := reflect.TypeOf(v)
	origType := t

	// If it's a slice, match on element type
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	s.logger.Debug("resolving filter",
		logging.String("event", "filter_resolve_start"),
		logging.String("filter", filter),
		logging.String("value_type", origType.String()),
		logging.String("match_type", t.String()),
		logging.String("correlation_id", cid),
	)

	// First try type-specific
	if typeMap, ok := s.byType[filter]; ok {
		if f, ok := typeMap[t]; ok {
			s.logger.Debug("resolved type-specific filter",
				logging.String("event", "filter_resolve_type_match"),
				logging.String("filter", filter),
				logging.String("type", t.String()),
				logging.String("correlation_id", cid),
			)
			return f, nil
		}

		s.logger.Debug("no type-specific filter match",
			logging.String("event", "filter_resolve_type_miss"),
			logging.String("filter", filter),
			logging.String("type", t.String()),
			logging.String("correlation_id", cid),
		)
	}

	// Fallback to generic
	if f, ok := s.registry[filter]; ok {
		s.logger.Debug("resolved generic filter",
			logging.String("event", "filter_resolve_generic_match"),
			logging.String("filter", filter),
			logging.String("correlation_id", cid),
		)
		return f, nil
	}

	s.logger.Error("failed to resolve filter",
		logging.String("event", "filter_resolve_error"),
		logging.String("filter", filter),
		logging.String("value_type", origType.String()),
		logging.String("correlation_id", cid),
	)

	return nil, fmt.Errorf("no filter registered for filter=%q and type=%T", filter, v)
}

// Filter implements [filtering.Service].
func (s *Service) Filter(filter string, v any) error {
	cid := util.CorrelationIDFromContext(context.Background())

	s.logger.Debug("filtering started",
		logging.String("event", "filter_apply_start"),
		logging.String("filter", filter),
		logging.String("value_type", reflect.TypeOf(v).String()),
		logging.String("correlation_id", cid),
	)

	filterer, err := s.resolve(filter, v)
	if err != nil {
		s.logger.Error("filter resolution failed",
			logging.String("event", "filter_apply_resolve_error"),
			logging.String("filter", filter),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if err := filterer.Filter(v); err != nil {
		s.logger.Error("filter execution failed",
			logging.String("event", "filter_apply_exec_error"),
			logging.String("filter", filter),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("filter applied successfully",
		logging.String("event", "filter_apply_success"),
		logging.String("filter", filter),
		logging.String("value_type", reflect.TypeOf(v).String()),
		logging.String("correlation_id", cid),
	)

	return nil
}

// RegisterWithType registers a filter for a specific Go type under a filter name
func (s *Service) RegisterWithType(filter string, typ reflect.Type, f domainfiltering.Filter) error {
	cid := util.CorrelationIDFromContext(context.Background())

	if s.byType == nil {
		s.logger.Error("type registry is nil",
			logging.String("event", "filter_register_type_error"),
			logging.String("filter", filter),
			logging.String("correlation_id", cid),
		)
		return errors.New("registry is nil")
	}

	origType := typ
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	s.logger.Info("registering type-specific filter",
		logging.String("event", "filter_register_type"),
		logging.String("filter", filter),
		logging.String("type", typ.String()),
		logging.String("original_type", origType.String()),
		logging.String("correlation_id", cid),
	)

	if _, ok := s.byType[filter]; !ok {
		s.byType[filter] = make(map[reflect.Type]domainfiltering.Filter)
	}
	s.byType[filter][typ] = f

	return nil
}

// Register implements [filtering.Service].
func (s *Service) Register(filter string, f domainfiltering.Filter) error {
	cid := util.CorrelationIDFromContext(context.Background())

	if s.registry == nil {
		s.logger.Error("generic registry is nil",
			logging.String("event", "filter_register_generic_error"),
			logging.String("filter", filter),
			logging.String("correlation_id", cid),
		)
		return errors.New("registry is nil")
	}

	s.logger.Info("registering generic filter",
		logging.String("event", "filter_register_generic"),
		logging.String("filter", filter),
		logging.String("correlation_id", cid),
	)

	s.registry[filter] = f
	return nil
}

package formatting

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Service struct {
	registry map[string]formatting.Formatter
	byType   map[string]map[reflect.Type]formatting.Formatter
	logger   logging.Logger
}

func NewService(logger logging.Logger) *Service {
	return &Service{
		registry: map[string]formatting.Formatter{},
		byType:   make(map[string]map[reflect.Type]formatting.Formatter),
		logger:   logger,
	}
}

// Register registers a formatter for a specific format name
func (s *Service) Register(format string, formatter formatting.Formatter) error {
	cid := util.CorrelationIDFromContext(context.Background())

	if s.registry == nil {
		s.logger.Error("generic formatter registry is nil",
			logging.String("event", "formatter_register_generic_error"),
			logging.String("format", format),
			logging.String("correlation_id", cid),
		)
		return errors.New("registry is nil")
	}

	s.logger.Info("registering generic formatter",
		logging.String("event", "formatter_register_generic"),
		logging.String("format", format),
		logging.String("correlation_id", cid),
	)

	s.registry[format] = formatter
	return nil
}

// RegisterWithType registers a formatter for a specific Go type under a format name
func (s *Service) RegisterWithType(format string, typ reflect.Type, f formatting.Formatter) error {
	cid := util.CorrelationIDFromContext(context.Background())

	if s.byType == nil {
		s.logger.Error("type-specific formatter registry is nil",
			logging.String("event", "formatter_register_type_error"),
			logging.String("format", format),
			logging.String("correlation_id", cid),
		)
		return errors.New("registry is nil")
	}

	origType := typ
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	s.logger.Info("registering type-specific formatter",
		logging.String("event", "formatter_register_type"),
		logging.String("format", format),
		logging.String("type", typ.String()),
		logging.String("original_type", origType.String()),
		logging.String("correlation_id", cid),
	)

	if _, ok := s.byType[format]; !ok {
		s.byType[format] = make(map[reflect.Type]formatting.Formatter)
	}
	s.byType[format][typ] = f

	return nil
}

// resolve chooses the best formatter for the given format + value
func (s *Service) resolve(format string, v any) (formatting.Formatter, error) {
	cid := util.CorrelationIDFromContext(context.Background())

	t := reflect.TypeOf(v)
	origType := t

	// If it's a slice, match on element type
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	s.logger.Debug("resolving formatter",
		logging.String("event", "formatter_resolve_start"),
		logging.String("format", format),
		logging.String("value_type", origType.String()),
		logging.String("match_type", t.String()),
		logging.String("correlation_id", cid),
	)

	// First try type-specific
	if typeMap, ok := s.byType[format]; ok {
		if f, ok := typeMap[t]; ok {
			s.logger.Debug("resolved type-specific formatter",
				logging.String("event", "formatter_resolve_type_match"),
				logging.String("format", format),
				logging.String("type", t.String()),
				logging.String("correlation_id", cid),
			)
			return f, nil
		}

		s.logger.Debug("no type-specific formatter match",
			logging.String("event", "formatter_resolve_type_miss"),
			logging.String("format", format),
			logging.String("type", t.String()),
			logging.String("correlation_id", cid),
		)
	}

	// Fallback to generic
	if f, ok := s.registry[format]; ok {
		s.logger.Debug("resolved generic formatter",
			logging.String("event", "formatter_resolve_generic_match"),
			logging.String("format", format),
			logging.String("correlation_id", cid),
		)
		return f, nil
	}

	s.logger.Error("failed to resolve formatter",
		logging.String("event", "formatter_resolve_error"),
		logging.String("format", format),
		logging.String("value_type", origType.String()),
		logging.String("correlation_id", cid),
	)

	return nil, fmt.Errorf("no formatter registered for format=%q and type=%T", format, v)
}

func (s *Service) Format(w io.Writer, format string, v any) error {
	cid := util.CorrelationIDFromContext(context.Background())

	s.logger.Debug("formatting started",
		logging.String("event", "formatter_apply_start"),
		logging.String("format", format),
		logging.String("value_type", reflect.TypeOf(v).String()),
		logging.String("correlation_id", cid),
	)

	formatter, err := s.resolve(format, v)
	if err != nil {
		s.logger.Error("formatter resolution failed",
			logging.String("event", "formatter_apply_resolve_error"),
			logging.String("format", format),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if err := formatter.Format(w, v); err != nil {
		s.logger.Error("formatter execution failed",
			logging.String("event", "formatter_apply_exec_error"),
			logging.String("format", format),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("formatting completed successfully",
		logging.String("event", "formatter_apply_success"),
		logging.String("format", format),
		logging.String("value_type", reflect.TypeOf(v).String()),
		logging.String("correlation_id", cid),
	)

	return nil
}

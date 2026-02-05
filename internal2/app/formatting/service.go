package formatting

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
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
	if s.registry == nil {
		return errors.New("registry is nil")
	}

	s.logger.Info("registering generic formatter",
		logging.String("format", format),
	)

	s.registry[format] = formatter
	return nil
}

// RegisterWithType registers a formatter for a specific Go type under a format name
func (s *Service) RegisterWithType(format string, typ reflect.Type, f formatting.Formatter) error {
	if s.byType == nil {
		return errors.New("registry is nil")
	}

	origType := typ
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	s.logger.Info("registering type-specific formatter",
		logging.String("format", format),
		logging.String("type", typ.String()),
		logging.String("originalType", origType.String()),
	)

	if _, ok := s.byType[format]; !ok {
		s.byType[format] = make(map[reflect.Type]formatting.Formatter)
	}
	s.byType[format][typ] = f

	return nil
}

// resolve chooses the best formatter for the given format + value
func (s *Service) resolve(format string, v any) (formatting.Formatter, error) {
	t := reflect.TypeOf(v)
	origType := t

	// If it's a slice, match on element type
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	s.logger.Debug("resolving formatter",
		logging.String("format", format),
		logging.String("valueType", origType.String()),
		logging.String("matchType", t.String()),
	)

	// First try type-specific
	if typeMap, ok := s.byType[format]; ok {
		if f, ok := typeMap[t]; ok {
			s.logger.Debug("resolved type-specific formatter",
				logging.String("format", format),
				logging.String("type", t.String()),
			)
			return f, nil
		}

		s.logger.Debug("no type-specific formatter match",
			logging.String("format", format),
			logging.String("type", t.String()),
		)
	}

	// Fallback to generic
	if f, ok := s.registry[format]; ok {
		s.logger.Debug("resolved generic formatter",
			logging.String("format", format),
		)
		return f, nil
	}

	s.logger.Error("failed to resolve formatter",
		logging.String("format", format),
		logging.String("valueType", origType.String()),
	)

	return nil, fmt.Errorf("no formatter registered for format=%q and type=%T", format, v)
}

func (s *Service) Format(w io.Writer, format string, v any) error {
	s.logger.Debug("formatting started",
		logging.String("format", format),
		logging.Any("valueType", reflect.TypeOf(v).String()),
	)

	formatter, err := s.resolve(format, v)
	if err != nil {
		s.logger.Error("formatter resolution failed",
			logging.String("format", format),
			logging.Error(err),
		)
		return err
	}

	if err := formatter.Format(w, v); err != nil {
		s.logger.Error("formatter execution failed",
			logging.String("format", format),
			logging.Error(err),
		)
		return err
	}

	s.logger.Info("formatting completed successfully",
		logging.String("format", format),
		logging.Any("valueType", reflect.TypeOf(v).String()),
	)

	return nil
}

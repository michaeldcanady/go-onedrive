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
		logger:   logger,
		byType:   make(map[string]map[reflect.Type]formatting.Formatter),
	}
}

// Register registers a formatter for a specific format name
func (s *Service) Register(format string, formatter formatting.Formatter) error {
	if s.registry == nil {
		return errors.New("registry is nil")
	}

	s.registry[format] = formatter
	return nil
}

// RegisterWithType registers a formatter for a specific Go type under a format name
func (s *Service) RegisterWithType(format string, typ reflect.Type, f formatting.Formatter) error {
	if s.byType == nil {
		return errors.New("registry is nil")
	}

	// If it's a slice, match on element type
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	if _, ok := s.byType[format]; !ok {
		s.byType[format] = make(map[reflect.Type]formatting.Formatter)
	}
	s.byType[format][typ] = f

	return nil
}

// resolve chooses the best formatter for the given format + value
func (s *Service) resolve(format string, v any) (formatting.Formatter, error) {
	// First try type-specific
	if typeMap, ok := s.byType[format]; ok {
		t := reflect.TypeOf(v)

		// If it's a slice, match on element type
		if t.Kind() == reflect.Slice {
			t = t.Elem()
		}

		if f, ok := typeMap[t]; ok {
			return f, nil
		}
	}

	// Fallback to generic
	if f, ok := s.registry[format]; ok {
		return f, nil
	}

	return nil, fmt.Errorf("no formatter registered for format=%q and type=%T", format, v)
}

func (s *Service) Format(w io.Writer, format string, v any) error {
	formatter, err := s.resolve(format, v)
	if err != nil {
		return err
	}

	return formatter.Format(w, v)
}

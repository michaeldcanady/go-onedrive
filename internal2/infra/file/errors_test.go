package file

import (
	"errors"
	"net/http"
	"testing"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
)

type mockKiotaError struct {
	statusCode int
}

func (e *mockKiotaError) Error() string {
	return "kiota error"
}

func (e *mockKiotaError) StatusCode() int {
	return e.statusCode
}

func TestMapGraphError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       error
		expectKind  error
		expectPanic bool
	}{
		{
			name:       "nil error",
			input:      nil,
			expectKind: nil,
		},
		{
			name: "ODataError: itemNotFound",
			input: func() error {
				e := odataerrors.NewODataError()
				mainErr := odataerrors.NewMainError()
				code := "itemNotFound"
				mainErr.SetCode(&code)
				e.SetErrorEscaped(mainErr)
				return e
			}(),
			expectKind: ErrNotFound,
		},
		{
			name: "ODataError: ErrorItemNotFound",
			input: func() error {
				e := odataerrors.NewODataError()
				mainErr := odataerrors.NewMainError()
				code := "ErrorItemNotFound"
				mainErr.SetCode(&code)
				e.SetErrorEscaped(mainErr)
				return e
			}(),
			expectKind: ErrNotFound,
		},
		{
			name: "ODataError: accessDenied",
			input: func() error {
				e := odataerrors.NewODataError()
				mainErr := odataerrors.NewMainError()
				code := "accessDenied"
				mainErr.SetCode(&code)
				e.SetErrorEscaped(mainErr)
				return e
			}(),
			expectKind: ErrForbidden,
		},
		{
			name:       "Kiota Error: 404",
			input:      &mockKiotaError{statusCode: http.StatusNotFound},
			expectKind: ErrNotFound,
		},
		{
			name:       "Kiota Error: 401",
			input:      &mockKiotaError{statusCode: http.StatusUnauthorized},
			expectKind: ErrUnauthorized,
		},
		{
			name:       "Generic Error",
			input:      errors.New("generic error"),
			expectKind: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mapGraphError(tt.input)

			if tt.input == nil {
				assert.NoError(t, got)
				return
			}

			assert.IsType(t, &DomainError{}, got)
			assert.Equal(t, tt.expectKind, got.(*DomainError).Kind)
			assert.Equal(t, tt.input, got.(*DomainError).Err)
		})
	}
}

func TestMapGraphError2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      error
		expectKind error
	}{
		{
			name:       "nil error",
			input:      nil,
			expectKind: nil,
		},
		{
			name: "ODataError: itemNotFound",
			input: func() error {
				e := odataerrors.NewODataError()
				mainErr := odataerrors.NewMainError()
				code := "itemNotFound"
				mainErr.SetCode(&code)
				e.SetErrorEscaped(mainErr)
				return e
			}(),
			expectKind: ErrNotFound,
		},
		{
			name:       "Kiota Error: 404",
			input:      &mockKiotaError{statusCode: http.StatusNotFound},
			expectKind: ErrNotFound,
		},
		{
			name:       "Generic Error",
			input:      errors.New("generic error"),
			expectKind: ErrInternal, // Actually returns *DomainError with Kind=ErrInternal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mapGraphError2(tt.input)

			if tt.input == nil {
				assert.NoError(t, got)
				return
			}

			if tt.name == "Generic Error" {
				var domainErr *DomainError
				if errors.As(got, &domainErr) {
					assert.Equal(t, ErrInternal, domainErr.Kind)
				} else {
					t.Fatalf("expected DomainError, got %T", got)
				}
			} else {
				assert.Equal(t, tt.expectKind, got)
			}
		})
	}
}

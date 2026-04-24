package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	pass := PolicyFunc[string](func(s string) error { return nil })
	fail1 := PolicyFunc[string](func(s string) error { return errors.New("error 1") })
	fail2 := PolicyFunc[string](func(s string) error { return errors.New("error 2") })

	tests := []struct {
		name     string
		policies []Policy[string]
		wantErr  bool
		errMsgs  []string
	}{
		{
			name:     "all pass",
			policies: []Policy[string]{pass, pass},
			wantErr:  false,
		},
		{
			name:     "one fails",
			policies: []Policy[string]{pass, fail1},
			wantErr:  true,
			errMsgs:  []string{"error 1"},
		},
		{
			name:     "multiple fail",
			policies: []Policy[string]{fail1, fail2},
			wantErr:  true,
			errMsgs:  []string{"error 1", "error 2"},
		},
		{
			name:     "with nil policy",
			policies: []Policy[string]{pass, nil, fail1},
			wantErr:  true,
			errMsgs:  []string{"error 1"},
		},
		{
			name:     "empty policies",
			policies: []Policy[string]{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := All(tt.policies...)
			err := p.Evaluate("test")
			if tt.wantErr {
				assert.Error(t, err)
				for _, msg := range tt.errMsgs {
					assert.Contains(t, err.Error(), msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAny(t *testing.T) {
	pass := PolicyFunc[string](func(s string) error { return nil })
	fail1 := PolicyFunc[string](func(s string) error { return errors.New("error 1") })
	fail2 := PolicyFunc[string](func(s string) error { return errors.New("error 2") })

	tests := []struct {
		name     string
		policies []Policy[string]
		wantErr  bool
		errMsgs  []string
	}{
		{
			name:     "all pass",
			policies: []Policy[string]{pass, pass},
			wantErr:  false,
		},
		{
			name:     "one passes",
			policies: []Policy[string]{fail1, pass},
			wantErr:  false,
		},
		{
			name:     "all fail",
			policies: []Policy[string]{fail1, fail2},
			wantErr:  true,
			errMsgs:  []string{"error 1", "error 2"},
		},
		{
			name:     "with nil policy",
			policies: []Policy[string]{nil, fail1},
			wantErr:  true,
			errMsgs:  []string{"error 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Any(tt.policies...)
			err := p.Evaluate("test")
			if tt.wantErr {
				assert.Error(t, err)
				for _, msg := range tt.errMsgs {
					assert.Contains(t, err.Error(), msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEach(t *testing.T) {
	pass := PolicyFunc[int](func(i int) error { return nil })
	fail := PolicyFunc[int](func(i int) error {
		if i > 10 {
			return errors.New("too big")
		}
		return nil
	})

	type container struct {
		Nums []int
	}

	tests := []struct {
		name      string
		candidate container
		policy    Policy[int]
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "all pass",
			candidate: container{Nums: []int{1, 2, 3}},
			policy:    pass,
			wantErr:   false,
		},
		{
			name:      "one fails",
			candidate: container{Nums: []int{1, 15, 3}},
			policy:    fail,
			wantErr:   true,
			errMsg:    "too big",
		},
		{
			name:      "empty slice",
			candidate: container{Nums: []int{}},
			policy:    fail,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Each(func(c container) []int { return c.Nums }, tt.policy)
			err := p.Evaluate(tt.candidate)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInList(t *testing.T) {
	allowed := []string{"apple", "banana"}

	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid value",
			value:     "apple",
			fieldName: "fruit",
			wantErr:   false,
		},
		{
			name:      "invalid value",
			value:     "cherry",
			fieldName: "fruit",
			wantErr:   true,
			errMsg:    "invalid fruit 'cherry'; please use one of the following valid options: apple, banana",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InList(tt.value, allowed, tt.fieldName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	type user struct {
		Name string
	}

	tests := []struct {
		name      string
		candidate user
		fieldName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "provided",
			candidate: user{Name: "Alice"},
			fieldName: "name",
			wantErr:   false,
		},
		{
			name:      "empty",
			candidate: user{Name: ""},
			fieldName: "name",
			wantErr:   true,
			errMsg:    "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Required(func(u user) string { return u.Name }, tt.fieldName)
			err := p.Evaluate(tt.candidate)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInListFunc(t *testing.T) {
	type config struct {
		Mode string
	}
	allowed := []string{"auto", "manual"}

	tests := []struct {
		name      string
		candidate config
		wantErr   bool
	}{
		{
			name:      "valid",
			candidate: config{Mode: "auto"},
			wantErr:   false,
		},
		{
			name:      "invalid",
			candidate: config{Mode: "other"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := InListFunc(allowed, "mode", func(c config) string { return c.Mode })
			err := p.Evaluate(tt.candidate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	passPolicy := PolicyFunc[string](func(s string) error { return nil })
	failPolicy := PolicyFunc[string](func(s string) error { return errors.New("fail") })

	t.Run("all pass", func(t *testing.T) {
		p := All(passPolicy, passPolicy)
		assert.NoError(t, p.Evaluate("test"))
	})

	t.Run("one fails", func(t *testing.T) {
		p := All(passPolicy, failPolicy)
		assert.Error(t, p.Evaluate("test"))
	})

	t.Run("all fail", func(t *testing.T) {
		p := All(failPolicy, failPolicy)
		err := p.Evaluate("test")
		assert.Error(t, err)
		// Should join errors
		assert.Contains(t, err.Error(), "fail\nfail")
	})
}

func TestAny(t *testing.T) {
	passPolicy := PolicyFunc[string](func(s string) error { return nil })
	failPolicy := PolicyFunc[string](func(s string) error { return errors.New("fail") })

	t.Run("all pass", func(t *testing.T) {
		p := Any(passPolicy, passPolicy)
		assert.NoError(t, p.Evaluate("test"))
	})

	t.Run("one passes", func(t *testing.T) {
		p := Any(passPolicy, failPolicy)
		assert.NoError(t, p.Evaluate("test"))
	})

	t.Run("all fail", func(t *testing.T) {
		p := Any(failPolicy, failPolicy)
		assert.Error(t, p.Evaluate("test"))
	})
}

func TestInList(t *testing.T) {
	allowed := []string{"a", "b", "c"}

	t.Run("is in list", func(t *testing.T) {
		assert.NoError(t, InList("a", allowed, "field"))
	})

	t.Run("not in list", func(t *testing.T) {
		err := InList("d", allowed, "field")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid field 'd'")
		assert.Contains(t, err.Error(), "a, b, c")
	})
}

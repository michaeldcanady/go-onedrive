package editor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEnvSvc struct {
	isWindows bool
	isLinux   bool
	isMac     bool
	shell     string
	editor    string
	visual    string
}

func (m *mockEnvSvc) HomeDir() (string, error)    { return "", nil }
func (m *mockEnvSvc) ConfigDir() (string, error)  { return "", nil }
func (m *mockEnvSvc) DataDir() (string, error)    { return "", nil }
func (m *mockEnvSvc) StateDir() (string, error)    { return "", nil }
func (m *mockEnvSvc) CacheDir() (string, error)    { return "", nil }
func (m *mockEnvSvc) TempDir() (string, error)     { return os.TempDir(), nil }
func (m *mockEnvSvc) RuntimeDir() (string, error)  { return "", nil }
func (m *mockEnvSvc) LogDir() (string, error)      { return "", nil }
func (m *mockEnvSvc) InstallDir() (string, error)  { return "", nil }
func (m *mockEnvSvc) EnsureAll() error             { return nil }
func (m *mockEnvSvc) IsWindows() bool              { return m.isWindows }
func (m *mockEnvSvc) IsLinux() bool                { return m.isLinux }
func (m *mockEnvSvc) IsMac() bool                  { return m.isMac }
func (m *mockEnvSvc) Shell() (string, error)       { return m.shell, nil }
func (m *mockEnvSvc) Editor() (string, error)      { return m.editor, nil }
func (m *mockEnvSvc) Visual() (string, error)      { return m.visual, nil }
func (m *mockEnvSvc) LogLevel() string             { return "info" }
func (m *mockEnvSvc) LogFormat() string            { return "json" }
func (m *mockEnvSvc) Name() string                 { return "odc" }
func (m *mockEnvSvc) OS() string                   { return "linux" }
func (m *mockEnvSvc) LogOutput() string            { return "stderr" }

func TestGetEditorCmd(t *testing.T) {
	env := &mockEnvSvc{}

	t.Run("explicit editor", func(t *testing.T) {
		s := NewDefaultService(env, nil, WithEditor("my-editor"))
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "my-editor", cmd)
	})

	t.Run("visual variable", func(t *testing.T) {
		env.visual = "visual-editor"
		s := NewDefaultService(env, nil)
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "visual-editor", cmd)
		env.visual = ""
	})

	t.Run("editor variable", func(t *testing.T) {
		env.editor = "editor-cmd"
		s := NewDefaultService(env, nil)
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "editor-cmd", cmd)
		env.editor = ""
	})

	t.Run("windows fallback", func(t *testing.T) {
		env.isWindows = true
		s := NewDefaultService(env, nil)
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "notepad.exe", cmd)
		env.isWindows = false
	})
}

func TestGetEditorParts(t *testing.T) {
	env := &mockEnvSvc{}

	t.Run("simple command", func(t *testing.T) {
		s := NewDefaultService(env, nil, WithEditor("vim"))
		parts, err := s.getEditorParts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"vim"}, parts)
	})

	t.Run("command with args", func(t *testing.T) {
		s := NewDefaultService(env, nil, WithEditor("code --wait"))
		parts, err := s.getEditorParts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"code", "--wait"}, parts)
	})

	t.Run("quoted args", func(t *testing.T) {
		s := NewDefaultService(env, nil, WithEditor(`"my editor" --args`))
		parts, err := s.getEditorParts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"my editor", "--args"}, parts)
	})
}

func TestWithOptions(t *testing.T) {
	env := &mockEnvSvc{}
	s := NewDefaultService(env, nil)

	s2 := s.WithOptions(WithEditor("new-editor"))
	assert.NotEqual(t, Service(s), s2)

	ds2, ok := s2.(*DefaultService)
	assert.True(t, ok)
	assert.Equal(t, "new-editor", ds2.editorCmd)

	// Original should not be modified
	assert.Equal(t, "", s.editorCmd)
}

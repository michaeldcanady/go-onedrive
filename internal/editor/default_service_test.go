package editor

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
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
func (m *mockEnvSvc) StateDir() (string, error)   { return "", nil }
func (m *mockEnvSvc) CacheDir() (string, error)   { return "", nil }
func (m *mockEnvSvc) TempDir() (string, error)    { return os.TempDir(), nil }
func (m *mockEnvSvc) RuntimeDir() (string, error) { return "", nil }
func (m *mockEnvSvc) LogDir() (string, error)     { return "", nil }
func (m *mockEnvSvc) InstallDir() (string, error) { return "", nil }
func (m *mockEnvSvc) EnsureAll() error            { return nil }
func (m *mockEnvSvc) IsWindows() bool             { return m.isWindows }
func (m *mockEnvSvc) IsLinux() bool               { return m.isLinux }
func (m *mockEnvSvc) IsMac() bool                 { return m.isMac }
func (m *mockEnvSvc) Shell() (string, error)      { return m.shell, nil }
func (m *mockEnvSvc) Editor() (string, error)     { return m.editor, nil }
func (m *mockEnvSvc) Visual() (string, error)     { return m.visual, nil }
func (m *mockEnvSvc) LogLevel() string            { return "info" }
func (m *mockEnvSvc) LogFormat() string           { return "json" }
func (m *mockEnvSvc) Name() string                { return "odc" }
func (m *mockEnvSvc) OS() string                  { return "linux" }
func (m *mockEnvSvc) LogOutput() string           { return "stderr" }

type mockConfigSvc struct {
	cfg config.Config
}

func (m *mockConfigSvc) GetEditorCommand(ctx context.Context) (string, error) {
	return m.cfg.Editor.Command, nil
}
func (m *mockConfigSvc) GetConfig(ctx context.Context) (config.Config, error) {
	return m.cfg, nil
}
func (m *mockConfigSvc) GetPath(ctx context.Context) (string, bool)              { return "", false }
func (m *mockConfigSvc) SaveConfig(ctx context.Context, cfg config.Config) error { return nil }
func (m *mockConfigSvc) SetOverride(ctx context.Context, path string) error      { return nil }

type mockLogger struct{}

func (m *mockLogger) Info(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Warn(msg string, kv ...logger.Field)           {}
func (m *mockLogger) Error(msg string, kv ...logger.Field)          {}
func (m *mockLogger) Debug(msg string, kv ...logger.Field)          {}
func (m *mockLogger) SetLevel(level logger.Level)                   {}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger     { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger { return m }

func TestGetEditorCmd(t *testing.T) {
	env := &mockEnvSvc{}
	l := &mockLogger{}

	t.Run("explicit editor", func(t *testing.T) {
		s := NewDefaultService(env, nil, l, WithEditor("my-editor"))
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "my-editor", cmd)
	})

	t.Run("config editor", func(t *testing.T) {
		cfgSvc := &mockConfigSvc{
			cfg: config.Config{
				Editor: config.EditorConfig{Command: "config-editor"},
			},
		}
		s := NewDefaultService(env, nil, l, WithConfig(cfgSvc))
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "config-editor", cmd)
	})

	t.Run("visual variable", func(t *testing.T) {
		env.visual = "visual-editor"
		s := NewDefaultService(env, nil, l)
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "visual-editor", cmd)
		env.visual = ""
	})

	t.Run("editor variable", func(t *testing.T) {
		env.editor = "editor-cmd"
		s := NewDefaultService(env, nil, l)
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "editor-cmd", cmd)
		env.editor = ""
	})

	t.Run("windows fallback", func(t *testing.T) {
		env.isWindows = true
		s := NewDefaultService(env, nil, l)
		cmd, err := s.getEditorCmd()
		assert.NoError(t, err)
		assert.Equal(t, "notepad.exe", cmd)
		env.isWindows = false
	})
}

func TestGetEditorParts(t *testing.T) {
	env := &mockEnvSvc{}
	l := &mockLogger{}

	t.Run("simple command", func(t *testing.T) {
		s := NewDefaultService(env, nil, l, WithEditor("vim"))
		parts, err := s.getEditorParts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"vim"}, parts)
	})

	t.Run("command with args", func(t *testing.T) {
		s := NewDefaultService(env, nil, l, WithEditor("code --wait"))
		parts, err := s.getEditorParts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"code", "--wait"}, parts)
	})

	t.Run("quoted args", func(t *testing.T) {
		s := NewDefaultService(env, nil, l, WithEditor(`"my editor" --args`))
		parts, err := s.getEditorParts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"my editor", "--args"}, parts)
	})
}

func TestWithOptions(t *testing.T) {
	env := &mockEnvSvc{}
	l := &mockLogger{}
	s := NewDefaultService(env, nil, l)

	s2 := s.WithOptions(WithEditor("new-editor"))
	assert.NotEqual(t, Service(s), s2)

	ds2, ok := s2.(*DefaultService)
	assert.True(t, ok)
	assert.Equal(t, "new-editor", ds2.editorCmd)

	// Original should not be modified
	assert.Equal(t, "", s.editorCmd)
}

func TestCreateSession(t *testing.T) {
	env := &mockEnvSvc{}
	l := &mockLogger{}
	factory := fs.NewURIFactory(nil)
	s := NewDefaultService(env, factory, l)

	t.Run("creates valid session", func(t *testing.T) {
		remoteURI := &fs.URI{Provider: "/onedrive", Path: "/test.txt"}
		content := "hello world"

		session, err := s.CreateSession(context.Background(), remoteURI, strings.NewReader(content))
		assert.NoError(t, err)
		assert.NotNil(t, session)
		defer s.Cleanup(context.Background(), session)

		assert.Equal(t, remoteURI, session.RemoteURI)
		assert.Contains(t, session.LocalURI.Path, "odc-edit-")
		assert.Equal(t, "/local", session.LocalURI.Provider)
		assert.Equal(t, StateCreated, session.State())
	})
}

func TestSessionStateTransitions(t *testing.T) {
	env := &mockEnvSvc{}
	l := &mockLogger{}
	factory := fs.NewURIFactory(nil)
	s := NewDefaultService(env, factory, l, WithEditor("true")) // Use 'true' as a command that always succeeds immediately

	t.Run("lifecycle transitions", func(t *testing.T) {
		remoteURI := &fs.URI{Provider: "/onedrive", Path: "/test.txt"}
		session, err := s.CreateSession(context.Background(), remoteURI, strings.NewReader("content"))
		assert.NoError(t, err)
		assert.Equal(t, StateCreated, session.State())

		err = s.Open(context.Background(), session)
		assert.NoError(t, err)
		assert.Equal(t, StateCompleted, session.State())

		modified, err := s.Modified(session)
		assert.NoError(t, err)
		assert.False(t, modified)

		err = s.Cleanup(context.Background(), session)
		assert.NoError(t, err)
		assert.Equal(t, StateClosed, session.State())
	})

	t.Run("invalid transitions", func(t *testing.T) {
		remoteURI := &fs.URI{Provider: "/onedrive", Path: "/test.txt"}
		session, err := s.CreateSession(context.Background(), remoteURI, strings.NewReader("content"))
		assert.NoError(t, err)

		// Cannot check modified before opening
		_, err = s.Modified(session)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot check modifications")

		// Cannot get content before opening
		_, err = s.NewContent(session)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get content")

		s.Cleanup(context.Background(), session)

		// Cannot open after cleanup
		err = s.Open(context.Background(), session)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transition")
	})
}

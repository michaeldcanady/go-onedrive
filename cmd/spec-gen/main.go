// spec-gen is a development tool that generates Cobra command boilerplate from
// Markdown-based specifications, ensuring consistency across feature slices.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
)

type Dependency struct {
	Method string
	Type   string
	Import string
	Alias  string
}

var dependencyRegistry = map[string]Dependency{
	"Logger":        {"Logger", "logger.Service", "github.com/michaeldcanady/go-onedrive/internal/core/logger", ""},
	"DB":            {"DB", "*persistence.DB", "github.com/michaeldcanady/go-onedrive/pkg/persistence", ""},
	"PluginManager": {"PluginManager", "plugins.Manager", "github.com/michaeldcanady/go-onedrive/internal/core/plugins", ""},
	"VFS":           {"VFS", "vfs.VFS", "github.com/michaeldcanady/go-onedrive/internal/features/vfs", ""},
	"Formatter":     {"Formatter", "format.Factory", "github.com/michaeldcanady/go-onedrive/pkg/format", ""},
	"Config":        {"Config", "config.Service", "github.com/michaeldcanady/go-onedrive/internal/features/config", ""},
	"Drive":         {"Drive", "drive.Service", "github.com/michaeldcanady/go-onedrive/internal/features/drive", ""},
	"Profile":       {"Profile", "profile.Service", "github.com/michaeldcanady/go-onedrive/internal/features/profile", ""},
	"Identity":      {"Identity", "identity.Service", "github.com/michaeldcanady/go-onedrive/internal/features/identity", ""},
	"Token":         {"Token", "identity.TokenService", "github.com/michaeldcanady/go-onedrive/internal/features/identity", ""},
	"Mounts":        {"Mounts", "mount.Service", "github.com/michaeldcanady/go-onedrive/internal/features/mount", ""},
	"Mount":         {"Mounts", "mount.Service", "github.com/michaeldcanady/go-onedrive/internal/features/mount", ""},
	"FS":            {"VFS", "vfs.VFS", "github.com/michaeldcanady/go-onedrive/internal/features/vfs", ""},
	"Editor":        {"Editor", "editor.Service", "github.com/michaeldcanady/go-onedrive/internal/features/editor", ""},
	"Resolver":      {"Resolver", "resolver.Service", "github.com/michaeldcanady/go-onedrive/internal/core/resolver", ""},
}

type Spec struct {
	Name         string   `yaml:"name"`
	Parent       string   `yaml:"parent"`
	Slice        string   `yaml:"slice"`
	Short        string   `yaml:"short"`
	Long         string   `yaml:"long"`
	Usage        string   `yaml:"usage"`
	Args         []Arg    `yaml:"args"`
	Flags        []Flag   `yaml:"flags"`
	Dependencies []string `yaml:"dependencies"`
}

type Arg struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
	Resolve     string `yaml:"resolve"`
}

type Flag struct {
	Name        string      `yaml:"name"`
	Shorthand   string      `yaml:"shorthand"`
	Type        string      `yaml:"type"`
	Default     interface{} `yaml:"default"`
	Description string      `yaml:"description"`
	Resolve     string      `yaml:"resolve"`
}

func pascal(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	parts := strings.Split(s, "-")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][0:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func camel(s string) string {
	s = pascal(s)
	if len(s) > 0 {
		return strings.ToLower(s[0:1]) + s[1:]
	}
	return s
}

func getDependency(name string) Dependency {
	if d, ok := dependencyRegistry[name]; ok {
		return d
	}
	return Dependency{Method: name, Type: "interface{}", Import: "", Alias: ""}
}

func getImports(spec Spec) []string {
	imports := make(map[string]struct{})
	for _, depName := range spec.Dependencies {
		dep := getDependency(depName)
		if dep.Import != "" {
			if dep.Alias != "" {
				imports[fmt.Sprintf("%s \"%s\"", dep.Alias, dep.Import)] = struct{}{}
			} else {
				imports[fmt.Sprintf("\"%s\"", dep.Import)] = struct{}{}
			}
		}
	}
	// Always include the internal logger interface
	imports["\"github.com/michaeldcanady/go-onedrive/internal/core/logger\""] = struct{}{}

	var res []string
	for imp := range imports {
		res = append(res, imp)
	}
	return res
}

func baseUsage(usage string) string {
	parts := strings.Split(usage, " ")
	for _, part := range parts {
		if part == "odc" {
			// Find the command name which follows 'odc' (and potentially parent commands)
			// Actually, we can just skip everything until the command name itself
			// But it's easier to just find the index of the command name
		}
	}
	// Simple heuristic: find the first part that isn't 'odc' or a parent command
	// For now, let's just return the part after the last command name
	return usage // Placeholder, I'll implement a better one
}

func getBaseUsage(spec Spec) string {
	usage := spec.Usage
	if spec.Parent != "" {
		usage = strings.Replace(usage, "odc "+spec.Parent+" ", "", 1)
	} else {
		usage = strings.Replace(usage, "odc ", "", 1)
	}
	return usage
}

func needsFmt(spec Spec) bool {
	for _, arg := range spec.Args {
		if arg.Resolve != "" {
			return true
		}
	}
	for _, flag := range spec.Flags {
		if flag.Resolve != "" {
			return true
		}
	}
	return false
}

var funcMap = template.FuncMap{
	"title":         strings.Title,
	"pascal":        pascal,
	"camel":         camel,
	"getDependency": getDependency,
	"getImports":    getImports,
	"getBaseUsage":  getBaseUsage,
	"needsFmt":      needsFmt,
}

func main() {
	specsDir := "specs/commands"
	files, err := os.ReadDir(specsDir)
	if err != nil {
		fmt.Printf("Error reading specs directory: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		path := filepath.Join(specsDir, file.Name())
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening spec %s: %v\n", path, err)
			continue
		}

		var spec Spec
		_, err = frontmatter.Parse(f, &spec)
		f.Close()
		if err != nil {
			fmt.Printf("Error parsing frontmatter in %s: %v\n", path, err)
			continue
		}

		if spec.Name == "" || spec.Slice == "" {
			fmt.Printf("Skipping %s: missing name or slice\n", path)
			continue
		}

		err = generateCommand(spec)
		if err != nil {
			fmt.Printf("Error generating command for %s: %v\n", spec.Name, err)
		}
	}
}

func generateCommand(spec Spec) error {
	outputDir := filepath.Join("internal/features", spec.Slice, "cmd", spec.Name)
	if spec.Parent != "" {
		outputDir = filepath.Join("internal/features", spec.Slice, "cmd", spec.Parent, spec.Name)
	}

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	templates := []struct {
		Name     string
		Target   string
		Optional bool
	}{
		{"command.go.tmpl", "command.go", false},
		{"options.go.tmpl", "options.go", false},
		{"handler_gen.go.tmpl", "handler_gen.go", false},
		{"handler.go.tmpl", "handler.go", true},
	}

	for _, t := range templates {
		outputPath := filepath.Join(outputDir, t.Target)
		if t.Optional {
			if _, err := os.Stat(outputPath); err == nil {
				continue // Skip optional file if it already exists
			}
		}

		tmplPath := filepath.Join("cmd/spec-gen/templates", t.Name)
		tmpl, err := template.New(t.Name).Funcs(funcMap).ParseFiles(tmplPath)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %w", t.Name, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, spec)
		if err != nil {
			return fmt.Errorf("error executing template %s: %w", t.Name, err)
		}

		err = os.WriteFile(outputPath, buf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("error writing to %s: %w", outputPath, err)
		}
		fmt.Printf("Generated %s for %s\n", t.Target, spec.Name)
	}

	return nil
}

---
name: e2e-test-writer
description: Generates End-to-End (E2E) tests for the CLI. These tests execute the compiled `odc` binary against a controlled environment (local files or mock server) to validate the entire application stack from CLI arguments to output.
---

# E2E Test Writer

## Overview
E2E tests treat the application as a black box. They verify that the compiled binary behaves correctly when executed with various flags and arguments. This ensures that the entry point, DI container, and all underlying services are correctly wired and functional.

## Workflow
1. **Build the Binary**: Ensure the `odc` binary is compiled (usually via `just build`).
2. **Setup Environment**:
    - Create a temporary directory for local file tests.
    - Set up environment variables (e.g., `ODC_CONFIG_PATH`).
3. **Define Scenarios**:
    - Command execution with various flag combinations.
    - Input redirection (stdin).
    - Output verification (stdout/stderr).
    - Exit code verification.
4. **Implement Test**: Use Go's `os/exec` or a library like `testscript`.
5. **Verify**: Run `go test ./tests/e2e/...`

## Patterns

### Command Execution Template
```go
func TestE2E_Ls(t *testing.T) {
    binary := "../../odc" // Path to compiled binary
    
    tests := []struct {
        name     string
        args     []string
        wantOut  string
        wantCode int
    }{
        {
            name:     "version command",
            args:     []string{"version"},
            wantOut:  "odc version",
            wantCode: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command(binary, tt.args...)
            output, err := cmd.CombinedOutput()
            
            assert.Equal(t, tt.wantCode, cmd.ProcessState.ExitCode())
            assert.Contains(t, string(output), tt.wantOut)
        })
    }
}
```

## Guidelines
- **Black Box**: Do not use internal project code or types in E2E tests.
- **Dependency on Build**: E2E tests require a built binary. Include instructions or `just` targets to ensure this happens.
- **Cleanup**: Always clean up temporary directories and environment changes.

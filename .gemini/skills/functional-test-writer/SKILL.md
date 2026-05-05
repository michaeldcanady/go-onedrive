---
name: functional-test-writer
description: Generates functional tests for CLI features. These tests verify complete business flows within the application boundary, mocking external I/O like the Microsoft Graph API but exercising multiple internal services together.
---

# Functional Test Writer

## Overview
Functional tests verify that a feature works as expected from an end-user perspective but within the Go test environment. They exercise multiple services (e.g., FS, Profile, Config) working together, typically mocking only the outermost boundaries (like network calls to OneDrive).

## Workflow
1. **Identify the Feature**: Determine the user-facing feature to test (e.g., "Profile Login Flow").
2. **Define Test Cases**:
    - **Happy Path**: Complete successful flow.
    - **Component Interaction**: Ensure services pass data correctly to each other.
    - **Error Propagation**: Verify that errors from deep services are surfaced correctly to the UI layer.
3. **Setup Test Environment**:
    - Initialize required services with real implementations where possible.
    - Mock external dependencies (e.g., `pkg/fs` implementations that talk to APIs).
    - Use a temporary directory for local file operations.
4. **Implement Test**:
    - Use the standard `go test` framework with `stretchr/testify`.
    - Follow the table-driven test pattern for multiple scenarios.
5. **Verify**: Run `go test -v ./internal/features/...`

## Patterns

### Functional Test Template
```go
func TestFeature_Functional(t *testing.T) {
    // Setup temporary environment
    tmpDir := t.TempDir()
    
    tests := []struct {
        name    string
        setup   func(ctx context.Context, deps *dependencies)
        action  func(ctx context.Context, deps *dependencies) error
        wantErr bool
    }{
        {
            name: "complete flow success",
            setup: func(ctx context.Context, deps *dependencies) {
                // Prepare state across multiple services
            },
            action: func(ctx context.Context, deps *dependencies) error {
                return deps.service.DoFeature(ctx)
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            deps := setupDependencies(t, tmpDir)
            if tt.setup != nil {
                tt.setup(ctx, deps)
            }
            
            err := tt.action(ctx, deps)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                // Verify side effects across services
            }
        })
    }
}
```

## Guidelines
- **Real Services**: Prefer using real service implementations over mocks for internal project services to test their integration.
- **External Mocks**: Strictly mock external I/O (API calls, OAuth flows).
- **Isolation**: Each test case should start with a clean state (e.g., fresh `bbolt` DB in a temp dir).

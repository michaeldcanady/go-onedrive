---
name: integration-test-writer
description: Generates integration tests for CLI commands, focusing on Cobra command wiring, flag parsing, and Dependency Injection (DI) integration. Use when testing how components work together within the CLI lifecycle.
---

# Integration Test Writer

## Overview
This skill provides a standardized workflow for writing CLI integration tests. Unlike unit tests, these tests verify the "wiring"—how flags map to options, how the DI container is used, and how the command lifecycle (PreRun, RunE, etc.) is orchestrated.

## Workflow

1. **Identify the Command**: Select the `cobra.Command` creator function (e.g., `CreateLsCmd`).
2. **Mock the DI Container**: Create a mock implementation of `di.Container` to provide mocked services to the command.
3. **Set Up Captured Output**: Use `cmd.SetOut` and `cmd.SetErr` with `bytes.Buffer` to capture and verify output.
4. **Define Test Cases**:
    - **Flag Parsing**: Verify that flags (e.g., `--recursive`, `-o json`) correctly populate the internal `Options` struct.
    - **Argument Handling**: Verify that positional arguments are correctly captured.
    - **Wiring Verification**: Ensure that the command calls the expected handler methods with the right parameters.
    - **Error Handling**: Verify that invalid flags or failing handlers result in appropriate CLI errors.
5. **Execute the Command**: Use `cmd.SetArgs(...)` followed by `cmd.ExecuteContext(ctx)`.
6. **Verify Results**: Assert on the captured stdout/stderr and the calls made to the mocked services.

## Implementation Patterns

### CLI Integration Template
```go
func TestCommand_Integration(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        setup   func(m *mockContainer)
        wantErr bool
        wantOut string
    }{
        {
            name: "list with recursive flag",
            args: []string{"ls", "od:/test", "--recursive"},
            setup: func(m *mockContainer) {
                // Setup expectations on mocked services
                m.fs.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]pkgfs.Item{}, nil)
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mContainer := new(mockContainer)
            if tt.setup != nil {
                tt.setup(mContainer)
            }

            cmd := CreateLsCmd(mContainer)
            buf := new(bytes.Buffer)
            cmd.SetOut(buf)
            cmd.SetArgs(tt.args[1:]) // Skip the command name in args

            err := cmd.Execute()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Contains(t, buf.String(), tt.wantOut)
            }
            mContainer.AssertExpectations(t)
        })
    }
}
```

## Guidelines
- **Mocking**: Use a single `mockContainer` that aggregates mocks for all services returned by the `di.Container` interface.
- **Independence**: Integration tests should not rely on a real filesystem or network; keep them fast by using service-level mocks.
- **Flag Verification**: Focus on testing that the Cobra command correctly translates user input into the `Options` struct passed to the handler.

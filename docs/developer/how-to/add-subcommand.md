# Adding a New Subcommand

The OneDrive CLI (`odc`) follows a consistent pattern for adding new 
commands. This guide explains how to create a new vertical slice for a 
subcommand.

## Overview of a Vertical Slice

A typical subcommand slice resides in 
`internal/<feature>/ui/cli/<command>/` and consists of three files:

1.  **`command.go`:** Defines the Cobra command, its flags, and validation 
    logic.
2.  **`command_cmd.go`:** Contains the `Command` which implements the 
    business logic for the command.
3.  **`options.go`:** Defines the options struct for the command and its 
    validation.

## Step-by-Step Guide

### 1. Create the Directory Structure

Create a new directory for your command. For example, if you're adding 
a `search` command to the `fs` feature:

```bash
mkdir -p internal/fs/ui/cli/search
```

### 2. Define the Options

In `options.go`, define the parameters your command needs.

```go
package search

import "io"

type Options struct {
    Query  string
    Stdout io.Writer
}

func (o Options) Validate() error {
    if o.Query == "" {
        return fmt.Errorf("search query is required")
    }
    return nil
}
```

### 3. Create the Command

In `command_cmd.go`, implement the `Command` and its `Handle` method.

```go
package search

import (
    "context"
    "github.com/michaeldcanady/go-onedrive/internal/fs"
    "github.com/michaeldcanady/go-onedrive/internal/features/logger"
)

type Command struct {
    fs     fs.Service
    logger logger.Logger
}

func NewCommand(fs fs.Service, logger logger.Logger) *Command {
    return &Command{fs: fs, logger: logger}
}

func (c *Command) Handle(ctx context.Context, opts Options) error {
    // Implement your command's logic here
    return nil
}
```

### 4. Define the Cobra Command

In `command.go`, use the `Create<Name>Cmd` pattern.

```go
package search

import (
    "github.com/michaeldcanady/go-onedrive/internal/di"
    "github.com/spf13/cobra"
)

func CreateSearchCmd(container di.Container) *cobra.Command {
    var opts Options

    cmd := &cobra.Command{
        Use:   "search <query>",
        Short: "Search for items",
        RunE: func(cmd *cobra.Command, args []string) error {
            l, _ := container.Logger().CreateLogger("search")
            handler := NewCommand(container.FS(), l)
            return handler.Handle(cmd.Context(), opts)
        },
    }

    // Add flags to opts here
    return cmd
}
```

### 5. Register the Command

Finally, register your new command in the root command in 
`internal/root/root.go`.

```go
rootCmd.AddCommand(
    // ... other commands
    search.CreateSearchCmd(container),
)
```

## Next steps

- **[Architecture Overview](../explanation/architecture.md)**
- **[Testing Your Changes](testing.md)**

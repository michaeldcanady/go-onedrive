# Developer setup

This guide provides step-by-step instructions for setting up your local 
development environment for the OneDrive CLI (`odc`)

## Prerequisites

Before you begin, confirm you have the following installed on your machine:

- **Go:** Version 1.25 or higher
- **Just:** A handy command runner (optional, but recommended)
- **Git:** For version control
- **An IDE:** Users recommend **VS Code** with the **Go** extension

## Getting the source code

1.  **Fork the repository** on GitHub
2.  **Clone your fork** to your local machine:

```bash
git clone https://github.com/<your-username>/go-onedrive.git
cd go-onedrive
```

## Installing dependencies

Install the required Go dependencies using the following command:

```bash
go get ./..
```

## Building the application

You can build the `odc` binary using `just` or the standard `go build` 
command

```bash
# Using just
just build

# Using go directly
go build -o odc ./cmd/odc/main.go
```

The resulting binary will be named `odc` in the root of the project

## Running tests

Confirm your environment is set up correctly by running the project's tests

```bash
# Using just
just test

# Using go directly
go test ./..
```

## Developing with devcontainers

The easiest way to get started is by using the provided **VS Code 
DevContainer**. This environment comes pre-configured with all the 
necessary tools and extensions

1.  Open the project in VS Code
2.  When prompted, click **Reopen in Container**

## Next steps

- **[Architecture Overview](../explanation/architecture.md)**
- **[Adding a New Command](../how-to/add-subcommand.md)**

# OneDrive CLI (odc) - Plugin Architecture

This document provides an overview of the plugin-based architecture for the `odc` project.

## Overview
`odc` (OneDrive CLI) leverages an extensible gRPC-based plugin system to support multiple storage backends and identity providers.

## Architecture
The system is divided into a core binary and independent plugin binaries. Communication occurs via gRPC sockets managed by `hashicorp/go-plugin`.

## Directory Structure
- `bin/plugins/`: Compiled plugin binaries.
- `internal/features/storage/proto/`: Storage service Protobuf definitions.
- `internal/features/identity/proto/`: Identity service Protobuf definitions.
- `cmd/storage-plugin-onedrive/`: OneDrive storage backend plugin implementation.
- `cmd/storage-plugin-local/`: Local filesystem storage backend plugin implementation.
- `cmd/identity-plugin-azure/`: Azure identity provider plugin implementation.

## Developing Plugins
To create a new plugin, implement the relevant Protobuf interface, create a plugin binary in `cmd/`, and register it in the plugin directory.

## Running odc with Plugins
1. Build all plugins using `just build`.
2. Ensure plugins are located in `./bin/plugins/`.
3. Configure mounts in `odc` as required.
4. Run `odc [command]`.

# 5. External Plugin Architecture using gRPC and go-plugin

Date: 2026-04-22

## Status

Status: Accepted

## Context

The `odc` project aims to be a flexible CLI tool capable of interacting with various storage backends beyond OneDrive, such as local filesystems, Google Drive, or S3. To maintain a lean core binary and enable third-party developers to contribute new storage providers without modifying the main codebase or requiring a full recompile, a robust plugin architecture is necessary.

## Decision

We have decided to implement an external plugin architecture using HashiCorp's `go-plugin` library. 
- **Protocol:** Plugins will communicate with the host process via gRPC over local Unix domain sockets or TCP.
- **Interfaces:** A set of standard interfaces for storage operations (e.g., `Read`, `Write`, `List`, `Delete`) will be defined using Protocol Buffers.
- **Discovery:** Plugins will be discovered as standalone binaries in a designated plugin directory.
- **Lifecycle:** The host process will manage the lifecycle (starting/stopping) of the plugin processes.

## Consequences

### Benefits
- **Decoupling:** Storage implementations are isolated from the core logic.
- **Extensibility:** New providers can be added easily by dropping a binary into the plugin folder.
- **Language Agnostic:** While primarily targeting Go, the use of gRPC allows plugins to be written in any language supported by gRPC.
- **Resilience:** A crash in a plugin process does not necessarily crash the main CLI tool.

### Trade-offs
- **Complexity:** Introduces the overhead of managing gRPC definitions and the `go-plugin` lifecycle.
- **Performance:** Inter-process communication (IPC) adds latency compared to in-process function calls.
- **Versioning:** Requires careful management of the Protobuf interface to maintain backward compatibility.

## Links

- [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)
- [gRPC](https://grpc.io/)

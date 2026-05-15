# Plugin Architecture Specification

## Overview
The system utilizes an external plugin architecture to offload storage and identity operations to standalone binary providers. This ensures the core application remains lightweight and extensible while providing isolation between the core and third-party integrations.

## IPC Mechanism
- **Transport:** gRPC over local Unix domain sockets or TCP.
- **Protocol:** A specialized host process spawns plugin binaries and manages their lifecycle.

## Handshake Protocol
To ensure compatibility, the host and plugins must agree on a magic cookie:
- **Key:** `ODC_PLUGIN`
- **Value:** `odc`
- **Version:** Defined by the host to ensure protocol compatibility.

## Plugin Lifecycle
1. **Discovery:** The host scans a designated directory for executable binaries.
2. **Handshake:** The host executes the binary and performs the handshake protocol.
3. **Registration:** The plugin registers its implementation of defined gRPC services.
4. **Active Session:** The host calls RPC methods, passing authentication tokens and path identifiers.
5. **Termination:** The host sends a shutdown signal to the plugin process when it is no longer needed.

## Statelessness Requirement
To ensure consistency and reliability, plugins (especially Identity Plugins) MUST be stateless:
- **No Local Persistence:** Plugins must not store tokens, credentials, or session data in local files or databases.
- **Host-Managed Lifecycle:** The host is responsible for managing the lifecycle of authentication tokens and session state.
- **On-Demand Execution:** Plugins are treated as on-demand adapters that perform specific operations and return results to the host.

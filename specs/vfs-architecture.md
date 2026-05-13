# VFS Architecture Specification

## Overview
The Virtual FileSystem (VFS) is the central orchestration layer. It presents a unified, hierarchical view of multiple disparate storage backends (e.g., cloud storage, local files) by mapping path prefixes to specific backend providers via "mount points."

## Key Components

### 1. VFS Orchestrator
- **Mount Management:** Maintains a registry of mount points mapping prefixes to backend instances.
- **Path Resolution:** Translates absolute VFS paths into a backend provider and a relative path.
- **Token Injection:** Proactively retrieves valid authentication tokens from the `TokenService` (orchestrating host-side refreshes if necessary) before delegating operations to backends.

### 2. URI System
The VFS uses a structured identifier to locate resources:
- `Provider`: The mount point name.
- `DriveID`: (Optional) The specific storage container identifier (e.g., a specific cloud drive).
- `Path`: The relative path within that provider.

### 3. Backend Interface
All storage providers must implement a standard interface, including:
- Metadata retrieval (Stat).
- Directory listing (List).
- Binary I/O (Open/Create).
- Structure management (Mkdir/Remove).
- Capability reporting (Native Copy/Move support).

## Operational Lifecycle

### Path Resolution Algorithm
1. Receive an absolute path.
2. Find the longest matching mount point prefix.
3. Strip the prefix to derive the relative path.
4. If no prefix matches, the operation fails.

### Request Delegation
1. Resolve the path to a backend.
2. Request a valid token from the `TokenService` for the backend's associated identity.
3. The `TokenService` performs a pre-flight check, refreshing the token via the plugin if expired.
4. Pass the verified token and relative path to the backend.
5. Handle streaming data or structured metadata responses.

## Resilience & Performance
- **Streaming:** Data transfers use streaming protocols to minimize memory footprint.
- **Lazy Auth:** Tokens are only requested at the moment of invocation.
- **Cross-Backend Operations:** For operations between different backends, the VFS orchestrates the transfer by reading from the source and writing to the destination.

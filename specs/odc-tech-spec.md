# Technical Specification - OneDrive CLI (odc)

## Overview
The `odc` (OneDrive CLI) is a terminal-native tool designed to treat cloud storage (specifically OneDrive) as a mountable, manageable filesystem. It leverages a modular architecture and an extensible plugin system to support diverse storage backends and identity providers.

## Core Pillars

### 1. Virtual FileSystem (VFS)
The VFS layer acts as a proxy, mapping user-facing paths to pluggable backends. It handles cross-backend transfers, path resolution, and transparent authentication token injection.

### 2. Extensible Plugin System
- **IPC:** gRPC-based communication between the host and external binaries.
- **Isolation:** Each plugin runs in its own process, ensuring the core remains stable and permitting language-agnostic plugin development.

### 3. Vertical Slice Architecture
Functionality is grouped by domain (e.g., storage, identity, mounting). To maintain slice isolation, inter-slice communication is governed by narrow, consumer-defined interfaces. The identity slice is further divided into:
- **`IdentityService`**: Manages user profiles, account metadata, and identity discovery.
- **`TokenService`**: Orchestrates the token lifecycle, performing host-side caching and refresh orchestration.

### 4. Dependency Injection & Service Communication
A centralized container manages the lifecycle and wiring of services. 

**Interface Locality Principle:**
To ensure loose coupling and maintain architectural integrity, `odc` follows the principle: *Ideally, packages should depend on interfaces they define locally that external packages happen to match, rather than depending on large, globally defined interfaces.*

- **Consumer-Defined Interfaces:** Consuming packages (e.g., a specific command handler or a cross-slice service) define a minimal interface containing only the methods they actually use.
- **Structural Typing:** Concrete services injected by the DI container satisfy these local interfaces without requiring the consumer to import the service implementation package.
- **Abstraction Strength:** We adhere to the Go proverb: "The bigger the interface, the weaker the abstraction." Small, local interfaces are easier to mock, test, and evolve.

## Data & Persistence
- **State Store:** Persistent storage is used for profiles, mount configurations, and cached credentials.
- **Credential Cache:** Authentication tokens are stored in a secure, host-managed cache (bbolt), keyed by a composite `provider:identity_id`.
- **Configuration:** User-defined settings are managed through a dedicated configuration layer, with support for profile-specific overrides.
- **Unified Node Model:** A consistent data structure represents files and directories across all integrated backends.

## Operational Flow
1. **Command Execution:** The CLI parses user input and invokes a specific domain handler.
2. **Dependency Wiring:** The handler is initialized with necessary services from the DI container.
3. **Logic Delegation:** The handler calls high-level domain services.
4. **VFS Orchestration:** The VFS resolves paths to specific backends and handles authentication requirements.
5. **Output Formatting:** Results are transformed into human or machine-readable formats for display.

## Security & Reliability
- **Host-Side Token Management:** Authentication tokens are stored and managed solely by the CLI host. Plugins are stateless adapters that perform login and refresh flows on demand.
- **Credential Protection:** Authentication tokens are stored securely in a restricted-access database and encrypted where platform support is available.
- **Streaming I/O:** Large data transfers are streamed to prevent memory exhaustion.
- **Error Handling:** Structured error classification ensures consistent and actionable feedback for the user.

## Logging Service
The logging service provides structured logging across the application. It is designed to be generic and configurable.

### Interface
The `logger.Service` interface includes:
- Standard logging levels: `Debug`, `Info`, `Warn`, `Error`, `Fatal`.
- Contextual logging: `With(keysAndValues ...any) Service`.
- Management tasks:
    - `Sync() error`: Flushes any buffered log entries.
    - `SetLevel(level string) error`: Dynamically changes the logging level at runtime.
    - `GetLevel() string`: Returns the current logging level.

### Implementation
- The logging service is responsible for creating and managing logger instances.
- It ensures that multiple logger instances (e.g., for different components) can be kept in sync regarding log levels and flushing.
- The service returns a generic `Service` interface, allowing for different implementations (e.g., Zap, Standard Log).
- The default implementation uses `uber-go/zap`.
- Log output is configurable (file, stdout, or both). Logs are typically directed to `~/.config/odc/logs/app.log`.

## Configuration Management
Configuration in `odc` is managed at the profile level, allowing for environment-specific settings.

### Default Configuration Values
The following defaults are hardcoded into the `ConfigService` and used as fallbacks when not present in the persistent store:
- `auth.provider`: `microsoft`
- `identity.microsoft.client_id`: `6b1e6ec0-ad93-4175-a0e0-84c02e13f206`
- `identity.microsoft.tenant_id`: `common`
- `identity.microsoft.redirect_uri`: `http://localhost:8400`

### User Discovery Mechanism
Identity plugins are responsible for performing "User Discovery" immediately following a successful login. 
1. The plugin obtains an access token via the configured OAuth2 flow.
2. The plugin uses the token to call the provider's metadata endpoint (e.g., `https://graph.microsoft.com/v1.0/me`).
3. The plugin extracts the unique ID, display name, and email to populate the `Identity` object returned to the host.

### Dynamic Key Parsing
- Configuration keys support a hierarchical dot-notation (e.g., `identity.azure.tenant_id`).
- The configuration service dynamically parses these keys into sections and subsections.
- Example: `odc config set identity.azure.tenant_id "common"`
    - Parses as `identity` -> `azure` -> `tenant_id`.
- This allows for flexible extension of configuration without needing to pre-define every possible key in a flat structure.

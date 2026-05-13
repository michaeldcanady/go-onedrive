# Implementation Plan: odc (OneDrive CLI)

This document outlines the phased implementation strategy for the odc CLI, adhering to the architectural principles defined in the ADRs and Technical Specifications.

## 1. Architectural Foundation

The project follows a Vertical Slice Architecture (ADR 0001) to ensure feature isolation and maintainability. Cross-cutting concerns are managed via a centralized Dependency Injection Container (ADR 0003).

### Key Principles
- Domain-First Access: CLI commands never bypass domain services (ADR 0004).
- Repository Pattern: Standardized data access interfaces (ADR 0008).
- Stateless Plugins: External plugins (gRPC/go-plugin) are stateless adapters; the host manages all persistence (ADR 0005).
- Modern Go: Context-aware, thread-safe, and interface-driven design (ADR 0009).

---

## 2. Phased Implementation Approach

### Phase 1: Foundation (Internal Core)
Establish the base infrastructure and persistence layer.
- [ ] Logging Service: Implement structured logging using Zap, wrapped in a generic interface.
- [ ] Storage Infrastructure:
    - Initialize bbolt database management (ADR 0002).
    - Implement base storage service for managing DB handles.
- [ ] Config Service: Implement dynamic key-value configuration with profile-level overrides.
- [ ] DI Container: Set up the container in internal/di/ to wire base services.
- [ ] Profile Service: Implement profile management (CRUD, current profile state).

### Phase 2: Plugin & Identity Infrastructure
Enable external plugin communication and authentication.
- [ ] gRPC Interfaces: Define Protobuf for IdentityPlugin and StoragePlugin.
- [ ] Plugin Manager: Implement the host-side logic for discovery, handshake, and lifecycle management (go-plugin).
- [ ] Identity Service:
    - Orchestrate login/logout flows.
    - Implement TokenService for host-side caching and refresh orchestration.
- [ ] Azure Identity Plugin: Implement the first plugin using azidentity to handle OneDrive authentication.

### Phase 3: VFS & Mount Management
Create the unified filesystem abstraction.
- [ ] Mount Service: Implement persistence and management of mount points (mapping prefixes to providers).
- [ ] VFS Orchestrator: 
    - Path resolution logic (mapping absolute paths to specific mounts).
    - Token injection (fetching valid tokens before calling storage plugins).
- [ ] Local Storage Plugin: Implement a plugin for local filesystem access to test VFS functionality.

### Phase 4: OneDrive Integration
Full OneDrive support via the storage plugin.
- [ ] OneDrive Storage Plugin: 
    - Implement msgraph-sdk-go integration.
    - Map OneDrive Drive/Item hierarchy to the VFS Node model.
    - Support large file streaming via Kiota.

### Phase 5: CLI Surface (Command Implementation)
Use spec-gen to generate command structures and implement handlers.
- [ ] Config Commands: config-get, config-set.
- [ ] Profile Commands: profile-create, profile-list, profile-use, profile-current, profile-delete.
- [ ] Identity Commands: identity-login, identity-logout, identity-list.
- [ ] Mount Commands: mount-add, mount-list, mount-remove.
- [ ] Drive Commands: drive-list, drive-get.
- [ ] Filesystem Commands: ls, stat, mkdir, rm, touch, cat, cp, mv, upload, download, edit.

### Phase 6: Hardening & Release
- [ ] Error Handling: Standardize domain error wrapping and CLI output.
- [ ] Performance: Implement concurrent transfers for cp/upload/download.
- [ ] CI/CD: Finalize GitHub Actions for linting, testing, and multi-platform releases via GoReleaser (ADR 0007).
- [ ] Documentation: Complete MkDocs content and command references.

---

## 3. Detailed Component Mapping

| Component | Responsibility | Implementation Detail |
| :--- | :--- | :--- |
| internal/core/di | Service Wiring | Singleton container providing validated services. |
| internal/features/storage | Persistence | bbolt repository implementation (ADR 0008). |
| internal/features/vfs | Orchestration | Path resolver and multi-backend dispatcher. |
| internal/features/identity| Auth Management | Token cache (host-side) and plugin coordination. |
| plugins/storage-* | Backend Adapter | gRPC server implementing storage interface. |

---

## 4. Command Checklist & Status

### Configuration & Profiles
- [ ] config-get: Retrieve configuration values.
- [ ] config-set: Update configuration values.
- [ ] profile-create: Create a new user profile.
- [ ] profile-list: List available profiles.
- [ ] profile-use: Set the active profile.
- [ ] profile-current: Show active profile details.
- [ ] profile-delete: Remove a profile.

### Identity & Authentication
- [ ] identity-login: Authenticate with a provider.
- [ ] identity-logout: Clear local session.
- [ ] identity-list: List authorized identities.

### Mount Management
- [ ] mount-add: Link a provider path to a logical VFS path.
- [ ] mount-list: List active mounts.
- [ ] mount-remove: Unlink a mount point.

### Drive Operations
- [ ] drive-list: List available drives/libraries in an account.
- [ ] drive-get: Get details for a specific drive.

### File Operations
- [ ] ls: List directory contents.
- [ ] stat: Display file/directory metadata.
- [ ] mkdir: Create directories.
- [ ] rm: Remove files/directories.
- [ ] touch: Update timestamps or create empty files.
- [ ] cat: Display file content.
- [ ] cp: Copy files/directories (cross-mount support).
- [ ] mv: Move files/directories (cross-mount support).
- [ ] upload: Upload local files to VFS.
- [ ] download: Download files from VFS to local.
- [ ] edit: Open a VFS file in the local $EDITOR.

---

## 5. Technical Debt & Refactoring Targets
- DI Container: Currently manually wired; consider google/wire if complexity grows (Ref Recommendation 03).
- Provider Decomposition: Ensure OneDrive specific logic stays in the plugin, not the VFS (Ref Recommendation 02).
- Error Consistency: Audit all slices for domain error wrapping (Ref Recommendation 05).

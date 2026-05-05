# Dependency injection

The OneDrive CLI (`odc`) uses a **Dependency Injection (DI)** container to 
manage the lifecycle and resolution of core services. This ensures that 
components are  coupled, easy to test, and follow a consistent 
initialization pattern

## The container pattern

Instead of components creating their own dependencies or relying on global 
variables, `odc` employs a central `Container` interface. This container 
acts as a registry for all shared services in the application

The container provides access to services such as:

- **Logger:** For structured logging
- **Config:** For managing user configuration
- **State:** For session and persistent runtime state
- **Identity:** For managing authentication and identity providers
- **Profile:** For profile management
- **FS:** For the filesystem abstraction layer
- **Drive:** For drive-specific operations
- **Mounts:** For managing virtual filesystem mount points

## Why use DI

Using a DI container offers some benefits for `odc`:

- **Decoupling:** Slices and services only know about the interfaces they 
  depend on, not the specific implementations
- **Testability:** You can  inject mock implementations of the 
  container or individual services during unit testing
- **Lifecycle Management:** The container ensures that services are 
  initialized in the correct order and only when needed (lazy loading where 
  appropriate)
- **Consistency:** All parts of the application access shared resources in 
   a uniform way

## How it works

1. **Definition:** The `Container` interface defines in `internal/di/`
2. **Implementation:** A concrete implementation of the container is 
   initialized at the application's entry point (`cmd/odc/main.go`)
3. **Injection:** The container is passed to the root command and 
   subsequently to each vertical slice
4. **Usage:** Slices retrieve the specific services they need from the 
   container

## Implementation details

The core logic for service resolution and wiring resides in:

- **`internal/di/container.go`:** Defines the `Container` interface
- **`internal/di/service.go`:** Contains the concrete implementation and 
  logic for wiring the services together

## Next steps

- **[Architecture Overview](architecture.md)**
- **[Configuration Management](configuration-management.md)**

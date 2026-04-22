# Authentication

The OneDrive CLI (`odc`) provides a flexible and modular system for managing 
authentication and identity. It is designed to support multiple 
authentication methods and is built to be easily extensible for different 
identity providers.

## Core Concepts

The authentication system is built around several key components:

- **Identity Providers:** These represent the cloud services (like 
  Microsoft OneDrive) that `odc` interacts with.
- **Authenticators:** These handle the logic for a specific identity 
  provider. They are responsible for performing the authentication flow and 
  returning a token.
- **Authentication Methods:** `odc` supports multiple ways to 
  authenticate, including:
    - **Interactive:** Opens a web browser for the user to log in.
    - **Device Code:** Displays a code for the user to enter in their 
      browser.
    - **Client Secret:** For service principal authentication.
    - **Environment:** Retrieves credentials from environment variables.

## How Authentication Works

1. **Initialization:** The application initializes the identity provider 
   registry, which contains the supported providers (currently 
   `microsoft`).
2. **Login:** When a user executes the `login` command, `odc` retrieves 
   the appropriate authenticator for the selected provider.
3. **Execution:** The authenticator performs the chosen authentication 
   flow (e.g., interactive login).
4. **Result:** A successful login returns an `AccessToken`, which 
   contains the token itself and its metadata (e.g., expiration).
5. **Persistence:** The token is cached securely for future use by the 
   active profile.

## The Identity Registry

`odc` uses a registry pattern to manage identity providers. This allows the 
application to resolve the correct authenticator based on the configuration 
at runtime.

- **`internal/features/identity/registry/`:** Contains the registry and service for 
  managing identity providers.
- **`internal/features/identity/providers/`:** Contains the implementations for 
  specific providers (e.g., `microsoft`).

## Next steps

- **[Profile Management](../how-to/manage-profiles.md)**
- **[Configuration Management](configuration-management.md)**

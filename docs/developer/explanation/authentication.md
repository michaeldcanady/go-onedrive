# Authentication

The OneDrive CLI (`odc`) provides a flexible and modular system for managing 
authentication and identity. It's designed to support multiple 
authentication methods and builds to be  extendable for different 
identity providers

## Core concepts

The authentication system builds around some key components:

- **Identity Providers:** These represent cloud services (like 
  Microsoft OneDrive) that `odc` interacts with
- **Authenticators:** These handle the logic for a specific identity 
  provider. They're responsible for performing the authentication flow and 
  returning a token
- **Authentication Methods:** `odc` supports multiple ways to 
  authenticate, including:
    - **Interactive:** Opens a web browser for the user to log in
    - **Device Code:** Displays a code for the user to enter in their 
      browser
    - **Client Secret:** For service principal authentication
    - **Environment:** Retrieves credentials from environment variables

## How authentication works

1. **Initialization:** The application initializes the identity provider 
   registry, which contains the supported providers (currently 
   `microsoft`)
2. **Login:** When a user executes the `login` command, `odc` retrieves 
   the appropriate authenticator for the selected provider
3. **Execution:** The authenticator performs the chosen authentication 
   flow (for example, interactive login)
4. **Result:** A successful login returns an `AccessToken`, which 
   contains the token itself and its metadata (for example, expiration)
5. **Persistence:** The token caches securely for future use by the 
   active profile

## The identity registry

`odc` uses a registry pattern to manage identity providers. This lets users the 
application to resolve the correct authenticator based on the configuration 
at runtime

- **`internal/features/identity/registry/`:** Contains the registry and service for 
  managing identity providers
- **`internal/features/identity/providers/`:** Contains the implementations for 
  specific providers (for example, `microsoft`)

## Next steps

- **[Profile Management](../../user/how-to/manage-profiles.md)**
- **[Configuration Management](configuration-management.md)**

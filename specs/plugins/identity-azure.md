# Plugin Specification: `identity-azure`

## Overview
The `identity-azure` plugin provides authentication and identity management using Azure Active Directory (Microsoft Entra ID).

## Capabilities
- **Login Methods:**
    - Interactive (Browser-based)
    - Device Code
    - Client Secret (Service Principal)
    - Environment Variables (Azure standard)
- **Token Management:**
    - Handles OAuth2 flow to obtain Access, Refresh, and ID tokens.
    - Returns tokens to the host for secure storage.
- **Identity Discovery:**
    - Retrieves user profile information (email, display name, tenant ID).

## Behavior
- **On `Login`**: 
    - Executes the configured authentication flow using `azidentity`.
    - Must respect the `redirect_uri` passed in the `options` map for Interactive flows.
    - Performs user discovery by calling `https://graph.microsoft.com/v1.0/me` via HTTP GET with the new access token.
    - Returns the `AccessToken` and discovered `Identity` to the host.
- **On `Refresh`**: Uses a refresh token (provided by the host) to obtain a new access token.
- **On `GetIdentity`**: Returns the user's profile details.

## Configuration Options
The following options can be set via `odc config set identity.microsoft.<key> <value>` (alias `azure` is supported):
- `client_id`: The Azure Application (client) ID.
- `tenant_id`: The Azure Tenant ID (e.g., `common`, `organizations`, or a UUID).
- `method`: The authentication method (`interactive`, `device`).
- `redirect_uri`: The loopback URI for interactive login (defaults to `http://localhost:8400`).

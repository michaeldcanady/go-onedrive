---
name: login
parent: identity
slice: identity
short: Authenticate with an identity provider
long: |
  Authenticate with a cloud provider (Azure, Google) using various methods (Interactive, Device Code, Service Principal).
  You can specify the provider and method via flags or in your configuration.
usage: odc identity login [flags]
flags:
  - name: provider
    type: string
    default: "azure"
    description: The identity provider to use (e.g., azure, google)
  - name: id
    type: string
    default: ""
    description: The specific identity (email) to authenticate
  - name: alias
    type: string
    default: ""
    description: An optional human-friendly name for this identity
  - name: show-token
    type: bool
    default: false
    description: Display the access token after login
  - name: force
    shorthand: f
    type: bool
    default: false
    description: Force re-authentication even if a valid profile exists
  - name: method
    type: string
    default: ""
    description: Authentication method (interactive, device-code, client-secret, environment)
  - name: tenant-id
    type: string
    default: ""
    description: Azure AD tenant ID (Azure only)
  - name: client-id
    type: string
    default: ""
    description: Client ID for the application
  - name: client-secret
    type: string
    default: ""
    description: Client secret for the application
  - name: scopes
    type: string
    default: ""
    description: Comma-separated list of scopes to request
dependencies:
  - Config
  - Identity
  - Token
  - Profile
  - Logger
---
# Command Specification: `identity login`

## Description
Authenticate with an identity provider and discover your identity.

## Behavior
- Invokes the `Login` method on the selected identity plugin.
- The plugin performs the authentication flow (e.g., opens a browser or provides a device code).
- Upon success, the CLI host receives the `AccessToken` and `Identity` metadata.
- The CLI host saves the identity metadata to the `IdentityService`.
- The CLI host saves the tokens to the `TokenService` (host-managed cache).

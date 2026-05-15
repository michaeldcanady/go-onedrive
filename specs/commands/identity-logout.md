---
name: logout
parent: identity
slice: identity
short: Sign out from OneDrive
long: Sign out from OneDrive for the active profile by clearing the cached authentication tokens.
usage: odc identity logout [flags]
flags:
  - name: id
    resolve: identity
    type: string
    default: ""
    description: The specific account to logout (optional)
  - name: force
    shorthand: f
    type: bool
    default: false
    description: Clear all cached credentials for the profile
dependencies:
  - Identity
  - Token
  - Profile
  - Logger
---
# Command Specification: `identity logout`

## Description
Sign out from OneDrive.

## Behavior
- Clears the cached authentication tokens from the `TokenService` for the active profile or specified identity.
- (Optional) Invokes the `Logout` method on the identity plugin to invalidate the remote session.
- Removes any transient identity state from the `IdentityService`.

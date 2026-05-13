---
name: list
parent: identity
slice: identity
short: List all authenticated identities
long: |
  List all authenticated identities managed by the host.
  This includes identities from all configured identity providers.
usage: odc identity list [flags]
flags:
  - name: format
    shorthand: o
    type: string
    default: table
    description: Output format (table, json, yaml)
dependencies:
  - Identity
  - Profile
  - Formatter
  - Logger
---
# Command Specification: `identity list`

## Description
List all authenticated identities managed by the host.

## Behavior
- Retrieves all identities stored in the `IdentityService`.
- Displays the identity ID (email), alias, and the provider used.
- Indicates which identity is currently active in the session profile (if any).
- For table output, it includes columns for:
    - **ACTIVE**: Marked with `*` if the identity is active in the current profile.
    - **ID**: The unique identifier (usually email).
    - **ALIAS**: User-friendly name or alias.
    - **PROVIDER**: The name of the identity provider (e.g., azure).

## Errors
- `failed to list identities`: Returned if the identity service fails to retrieve records.
- `unsupported format`: Returned if the requested output format is not supported.

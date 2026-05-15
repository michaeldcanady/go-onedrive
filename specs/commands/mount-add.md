---
name: add
parent: mount
slice: mount
short: Add a mount point
usage: odc mount add <path> <type> <identity-id> [flags]
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The path where the mount point will be created.
  - name: type
    type: string
    required: true
    description: The type of storage backend (e.g., local, onedrive, googledrive).
  - name: identity-id
    resolve: identity
    type: string
    required: true
    description: The identity to use for the mount point.
flags:
  - name: identity-provider
    type: string
    default: "azure"
    description: The identity provider to use (e.g., azure, google)
  - name: option
    type: stringSlice
    default: []
    description: Provider-specific options in key=value format (repeatable)
dependencies:
  - Mounts
  - Profile
  - Logger
---
# Command Specification: `mount add`

## Description
Add a mount point.

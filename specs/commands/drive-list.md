---
name: list
parent: drive
slice: drive
short: List all available storage drives
long: Retrieve all storage drives associated with your authenticated accounts, distinguishing between mounted and unmounted drives.
usage: odc drive list [flags]
flags:
  - name: id
    resolve: identity
    type: string
    default: ""
    description: The specific identity (email or alias) to list drives for
  - name: format
    shorthand: o
    type: string
    default: table
    description: Output format (table, json, yaml)
  - name: all
    shorthand: a
    type: bool
    default: false
    description: Show all drives, including those not currently mounted
dependencies:
  - Drive
  - Mount
  - Mounts
  - Profile
  - Formatter
  - Logger
---
# Command Specification: `drive list`

## Description
Retrieve and list all drives associated with authenticated plugin accounts.

## Behavior
- Lists all drives discovered across configured identity providers.
- Cross-references discovered drives with active mount points in the `MountService`.
- If a drive is mounted, it shows the mount path.
- By default, it shows drives that are either mounted or available for the current identity.
- Use the `--all` flag to see all drives discovered across all authenticated identities.
- Table output columns:
    - **MOUNTED**: Path where the drive is mounted, or empty if unmounted.
    - **ID**: The drive identifier.
    - **NAME**: The display name of the drive.
    - **IDENTITY**: The identity associated with the drive.
    - **TYPE**: Drive type (e.g., business, personal).

## Errors
- `failed to list drives`: Returned if the drive discovery service encounters an error.
- `unsupported format`: Returned if the requested format is not valid.

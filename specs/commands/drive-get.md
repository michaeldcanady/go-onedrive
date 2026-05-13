---
name: get
parent: drive
slice: drive
short: Display details for a specific drive
long: Retrieve and show the metadata for a OneDrive drive identified by its ID or name.
usage: odc drive get <drive-ref> [flags]
args:
  - name: drive-ref
    resolve: drive
    type: string
    required: true
    description: The identifier for the drive, which can be its ID or name.
flags:
  - name: id
    resolve: identity
    type: string
    default: ""
    description: The specific identity (email or alias) to get the personal drive for
dependencies:
  - Drive
  - Profile
  - Logger
---
# Command Specification: `drive get`


## Description
Retrieve and display the metadata for a specific OneDrive drive.

## Usage
`odc drive get <drive-ref> [flags]`

## Arguments
- `<drive-ref>`: The identifier for the drive, which can be its ID or name.

## Flags
| Flag | Description | Default |
| :--- | :--- | :--- |
| `--id` | The specific identity (email or alias) to get the personal drive for | *None* |

## Behavior
- Resolves the specified drive using the available identity and mount configurations.
- Displays the drive's name, ID, and type.

## Errors
- `failed to resolve drive`: Returned if the drive cannot be found or if there is an error during resolution.
- `missing argument`: Returned if the drive reference is not provided.

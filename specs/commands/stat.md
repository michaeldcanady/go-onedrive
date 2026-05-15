---
name: stat
slice: fs
short: Display file or directory status
usage: odc stat <path>
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The filesystem path to the item to status.
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `stat`

## Description
Display detailed information about a file or directory.

## Usage
`odc stat <path>`

## Arguments
- `<path>`: The filesystem path to the item to status.

## Flags
*None*

## Behavior
- Retrieves metadata for the specified path and displays it in a human-readable format.

## Errors
- `failed to stat`: Returned if the item metadata cannot be retrieved.

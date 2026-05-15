---
name: touch
slice: fs
short: Create a new empty file
usage: odc touch <path>
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The path to the file.
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `touch`


## Description
Create a new empty file or update the timestamp of an existing file.

## Usage
`odc touch <path>`

## Arguments
- `<path>`: The path to the file.

## Flags
*None*

## Behavior
- Creates an empty file if one does not exist, or updates its metadata if it does.

## Errors
- `invalid path`: Returned if the path cannot be resolved.
- `failed to touch`: Returned if the operation fails.

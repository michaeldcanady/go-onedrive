---
name: rm
slice: fs
short: Remove files and directories
usage: odc rm <path>
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The path to the file or directory to remove.
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `rm`


## Description
Remove files or directories.

## Usage
`odc rm <path>`

## Arguments
- `<path>`: The path to the file or directory to remove.

## Flags
*None*

## Behavior
- Deletes the file or directory at the specified path.

## Errors
- `invalid path`: Returned if the path cannot be resolved.
- `failed to remove`: Returned if the operation fails.

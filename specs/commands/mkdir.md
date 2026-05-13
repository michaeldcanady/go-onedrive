---
name: mkdir
slice: fs
short: Create a new directory
usage: odc mkdir <path>
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The path where the new directory will be created.
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `mkdir`


## Description
Create a new directory.

## Usage
`odc mkdir <path>`

## Arguments
- `<path>`: The path where the new directory will be created.

## Flags
*None*

## Behavior
- Creates a new directory at the specified path.

## Errors
- `invalid path`: Returned if the path cannot be resolved.
- `failed to create directory`: Returned if the operation fails (e.g., parent directory missing or permissions).

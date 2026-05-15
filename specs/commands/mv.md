---
name: mv
slice: fs
short: Move files and directories
usage: odc mv <source> <destination>
args:
  - name: source
    resolve: path
    type: string
    required: true
    description: The current path of the item.
  - name: destination
    resolve: path
    type: string
    required: true
    description: The new path of the item.
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `mv`


## Description
Move or rename files and directories.

## Usage
`odc mv <source> <destination>`

## Arguments
- `<source>`: The current path of the item.
- `<destination>`: The new path of the item.

## Flags
*None*

## Behavior
- Moves or renames the source item to the destination path.

## Errors
- `invalid source/destination path`: Returned if paths cannot be resolved.
- `failed to move`: Returned if the operation fails.

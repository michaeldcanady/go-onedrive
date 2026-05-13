---
name: cp
slice: fs
short: Copy files and directories
usage: odc cp <source> <destination> [flags]
args:
  - name: source
    resolve: path
    type: string
    required: true
    description: The path to the item to copy.
  - name: destination
    resolve: path
    type: string
    required: true
    description: The path where the item should be copied.
flags:
  - name: recursive
    shorthand: r
    type: bool
    default: false
    description: Copy directories recursively
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `cp`


## Description
Copy files and directories.

## Usage
`odc cp <source> <destination> [flags]`

## Arguments
- `<source>`: The path to the item to copy.
- `<destination>`: The path where the item should be copied.

## Flags
| Flag | Description | Default |
| :--- | :--- | :--- |
| `-r`, `--recursive` | Copy directories recursively | `false` |

## Behavior
- Copies the source item to the destination path.
- If the source is a directory, the recursive flag must be set.

## Errors
- `invalid source/destination path`: Returned if paths cannot be resolved.
- `failed to copy`: Returned if the copy operation fails.

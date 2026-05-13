---
name: upload
slice: fs
short: Upload files and directories
usage: odc upload <source> <destination> [flags]
args:
  - name: source
    resolve: path
    type: string
    required: true
    description: The local path to the file or directory.
  - name: destination
    resolve: path
    type: string
    required: true
    description: The remote path on OneDrive where the item should be uploaded.
flags:
  - name: recursive
    shorthand: r
    type: bool
    default: false
    description: Upload directories recursively
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `upload`


## Description
Upload files and directories from the local filesystem to OneDrive.

## Usage
`odc upload <source> <destination> [flags]`

## Arguments
- `<source>`: The local path to the file or directory.
- `<destination>`: The remote path on OneDrive where the item should be uploaded.

## Flags
| Flag | Description | Default |
| :--- | :--- | :--- |
| `-r`, `--recursive` | Upload directories recursively | `false` |

## Behavior
- Uploads the local item to the specified destination path.
- Handles both single files and directory trees (with `-r`).

## Errors
- `invalid source/destination path`: Returned if paths cannot be resolved.
- `failed to upload`: Returned if the upload operation fails.

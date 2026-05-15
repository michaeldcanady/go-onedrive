---
name: download
slice: fs
short: Download files and directories
usage: odc download <source> <destination> [flags]
args:
  - name: source
    resolve: path
    type: string
    required: true
    description: The remote path on OneDrive.
  - name: destination
    resolve: path
    type: string
    required: true
    description: The local path where the item should be downloaded.
flags:
  - name: recursive
    shorthand: r
    type: bool
    default: false
    description: Download directories recursively
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `download`


## Description
Download files and directories from OneDrive to the local filesystem.

## Usage
`odc download <source> <destination> [flags]`

## Arguments
- `<source>`: The remote path on OneDrive.
- `<destination>`: The local path where the item should be downloaded.

## Flags
| Flag | Description | Default |
| :--- | :--- | :--- |
| `-r`, `--recursive` | Download directories recursively | `false` |

## Behavior
- Downloads the remote item to the specified local destination path.
- Handles both single files and directory trees (with `-r`).

## Errors
- `invalid source/destination path`: Returned if paths cannot be resolved.
- `failed to download`: Returned if the download operation fails.

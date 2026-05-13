---
name: edit
slice: editor
short: Edit a file
long: Open a file in your default editor. Changes are synced back when the editor closes.
usage: odc edit <path> [flags]
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The path to the file to be edited.
flags:
  - name: editor
    type: string
    default: ""
    description: Editor to use (overrides config and environment)
  - name: force
    type: bool
    default: false
    description: Force upload even if the remote file has changed
dependencies:
  - FS
  - Profile
  - Editor
  - Logger
---
# Command Specification: `edit`


## Description
Open a file in your default editor. Changes are synced back when the editor closes.

## Usage
`odc edit <path> [flags]`

## Arguments
- `<path>`: The path to the file to be edited.

## Flags
| Flag | Description | Default |
| :--- | :--- | :--- |
| `--editor` | Editor to use (overrides config and environment) | *System Default* |
| `--force` | Force upload even if the remote file has changed | `false` |

## Behavior
1. Downloads the remote file to a temporary local location and captures the remote ETag.
2. Launches the configured editor.
3. Detects if the file was modified upon editor exit.
4. If modified, uploads the new content back to the remote location.
5. Performs an optimistic concurrency check using the captured ETag, unless `--force` is specified.
6. Cleans up temporary files.

## Errors
- `invalid path`: Returned if the provided path cannot be resolved.
- `failed to read file`: Returned if the remote file cannot be downloaded.
- `failed to create editor session`: Returned if the temporary environment cannot be prepared.
- `editor session failed`: Returned if the editor execution fails.
- `failed to check for modifications`: Returned if the file's modification status cannot be verified.
- `failed to write changes`: Returned if the updated content cannot be uploaded back or if the remote version has changed.

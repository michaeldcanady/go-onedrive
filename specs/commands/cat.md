---
name: cat
slice: fs
short: Display file contents
usage: odc cat <path>
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The filesystem path to the file to display.
dependencies:
  - FS
  - Profile
  - Logger
---
# Command Specification: `cat`


## Description
Display the contents of a file.

## Usage
`odc cat <path>`

## Arguments
- `<path>`: The filesystem path to the file to display.

## Flags
*None*

## Behavior
- Reads the file content from the specified path and writes it to standard output.

## Errors
- `invalid path`: Returned if the path cannot be resolved.
- `failed to open file`: Returned if the file does not exist or cannot be accessed.

---
name: ls
slice: fs
short: List items in a directory
long: List the items in a specified directory in OneDrive or the local filesystem.
usage: odc ls [path] [flags]
args:
  - name: path
    resolve: path
    type: string
    required: false
    description: The path to the directory to list (defaults to the current working directory).
flags:
  - name: format
    shorthand: o
    type: string
    default: short
    description: Output format (short, long, json, yaml, tree, table)
  - name: recursive
    shorthand: r
    type: bool
    default: false
    description: List items recursively
  - name: all
    shorthand: a
    type: bool
    default: false
    description: Show hidden items
  - name: sort
    type: stringSlice
    default: ["name"]
    description: Sort items by field (name, size, modified)
  - name: desc
    type: bool
    default: false
    description: Sort in descending order
dependencies:
  - FS
  - Profile
  - Formatter
  - Logger
---
# Command Specification: `ls`


## Description

List items in a directory.

## Usage

`odc ls [path] [flags]`

## Arguments

- `path` (optional): The path to the directory to list (defaults to the current working directory).

## Flags

| Flag                | Description                                          | Default    |
| :------------------ | :--------------------------------------------------- | :--------- |
| `-o`, `--format`    | Output format (short, long, json, yaml, tree, table) | `short`    |
| `-r`, `--recursive` | List items recursively                               | `false`    |
| `-a`, `--all`       | Show hidden items                                    | `false`    |
| `--sort`            | Sort items by field (name, size, modified)           | `["name"]` |
| `--desc`            | Sort in descending order                             | `false`    |

## Behavior

- Lists all files and directories at the specified path.
- Supports recursive listing for compatible output formats.

## Errors

- `list failed`: Returned if the directory cannot be listed.
- `unknown output format`: Returned if an unsupported format is provided.
- `recursive mode not supported`: Returned if recursion is used with an incompatible format.

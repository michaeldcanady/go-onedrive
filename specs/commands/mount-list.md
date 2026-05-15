---
name: list
parent: mount
slice: mount
short: List all mount points
usage: odc mount list [flags]
flags:
  - name: format
    shorthand: o
    type: string
    default: table
    description: Output format (table, json, yaml)
dependencies:
  - Mounts
  - Profile
  - Formatter
  - Logger
---
# Command Specification: `mount list`

## Description
List all configured mount points.

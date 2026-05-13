---
name: list
parent: profile
slice: profile
short: List all profiles
usage: odc profile list [flags]
flags:
  - name: format
    shorthand: o
    type: string
    default: table
    description: Output format (table, json, yaml)
dependencies:
  - Profile
  - Formatter
  - Logger
---
# Command Specification: `profile list`

## Description
List all profiles.

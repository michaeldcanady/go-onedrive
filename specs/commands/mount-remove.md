---
name: remove
parent: mount
slice: mount
short: Remove a mount point
usage: odc mount remove <path>
args:
  - name: path
    resolve: path
    type: string
    required: true
    description: The path of the mount point to remove.
dependencies:
  - Mounts
  - Profile
  - Logger
---
# Command Specification: `mount remove`

## Description
Remove an existing mount point.

---
name: use
parent: profile
slice: profile
short: Set the active profile
long: Specify a profile name to be used for subsequent commands.
usage: odc profile use <name>
args:
  - name: name
    type: string
    required: true
    description: The name of the profile to use.
dependencies:
  - Profile
  - Logger
---
# Command Specification: `profile use`

## Description
Set the active profile.

---
name: create
parent: profile
slice: profile
short: Create a new profile
long: Create a new profile with the specified name.
usage: odc profile create <name>
args:
  - name: name
    type: string
    required: true
    description: The name of the profile to create.
dependencies:
  - Profile
  - Logger
---
# Command Specification: `profile create`

## Description
Create a new profile.

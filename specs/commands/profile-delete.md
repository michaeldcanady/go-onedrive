---
name: delete
parent: profile
slice: profile
short: Delete a profile
long: Permanently delete a profile and its associated configuration and authentication tokens.
usage: odc profile delete <name>
args:
  - name: name
    type: string
    required: true
    description: The name of the profile to delete.
dependencies:
  - Profile
  - Logger
---
# Command Specification: `profile delete`

## Description
Delete a profile.

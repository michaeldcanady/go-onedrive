---
name: set
parent: config
slice: config
short: Set configuration
usage: odc config set <key> <value>
args:
  - name: key
    type: string
    required: true
    description: The configuration key to update.
  - name: value
    type: string
    required: true
    description: The new value for the configuration setting.
dependencies:
  - Config
  - Logger
---
# Command Specification: `config set`


## Description
Set the value of a configuration setting.

## Usage
`odc config set <key> <value>`

## Arguments
- `key`: The configuration key to update.
- `value`: The new value for the configuration setting.

## Flags
*None*

## Behavior
- Updates the specified configuration key with the provided value in the active profile's configuration.

## Errors
- `configuration update failed`: Returned if the key is invalid or the value cannot be set.

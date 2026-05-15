---
name: get
parent: config
slice: config
short: Get configuration
usage: odc config get [key] [flags]
args:
  - name: key
    type: string
    required: true
    description: The configuration key to retrieve (e.g., auth.provider, logging.level).
flags:
  - name: format
    shorthand: o
    type: string
    default: value
    description: Output format (value, json, yaml)
dependencies:
  - Config
  - Formatter
  - Logger
---
# Command Specification: `config get`


## Description
Retrieve the value of a specific configuration key.

## Usage
`odc config get [key]`

## Arguments
- `key`: The configuration key to retrieve (e.g., `auth.provider`, `logging.level`).

## Flags
*None*

## Behavior
- Retrieves the configuration setting for the active profile and displays its value.
- If the key is not supported, an error is returned.

## Errors
- `invalid path`: Returned if the configuration path cannot be resolved.
- `configuration key not supported`: Returned if the requested key is not recognized.

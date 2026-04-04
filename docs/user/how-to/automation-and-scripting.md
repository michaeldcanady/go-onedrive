# Automation and Scripting

`odc` (OneDrive CLI) is designed to be a first-class citizen in your terminal
workflow. This guide covers the features that make `odc` powerful for
automation, from machine-readable output to environment-driven configuration.

## Machine-Readable Output

Most `odc` commands support the `-o json` (or `--format json`) flag. This
provides a consistent JSON structure that is perfect for parsing with tools
like `jq`.

### Common jq patterns

- **List filenames and IDs:**
  ```bash
  odc ls /Documents -o json | jq '.[] | {name: .name, id: .id}'
  ```

- **Find the largest file in a directory:**
  ```bash
  odc ls /Videos -o json | jq 'sort_by(.size) | last'
  ```

- **Filter files by extension:**
  ```bash
  odc ls /Photos -o json | jq '.[] | select(.name | endswith(".jpg"))'
  ```

## Integration with Shell Tools

Since `odc` follows Unix philosophy, you can easily combine it with other
powerful shell tools like `xargs`, `find`, and `grep`.

### Bulk operations with xargs
You can use `odc` in combination with `xargs` to perform bulk operations
efficiently.

```bash
# Delete all files in a directory ending in .tmp
odc ls /Temp -o json | jq -r '.[] | select(.name | endswith(".tmp")) | .path' | xargs -I {} odc rm {}
```

### Scripting complex workflows
You can use `odc` within bash or zsh scripts to automate complex data
management tasks.

```bash
#!/bin/bash
# Backup a local folder to OneDrive with a timestamp
TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
BACKUP_DIR="/Backups/Backup_$TIMESTAMP"

echo "Creating backup directory: $BACKUP_DIR"
odc mkdir "$BACKUP_DIR"

echo "Uploading project files..."
odc upload -r ./my_project "$BACKUP_DIR"
```

## Environment Variables

You can configure `odc` using environment variables, which is especially
useful in CI/CD pipelines where interactive configuration isn't possible.

| Variable | Description |
| :--- | :--- |
| `ODC_CONFIG` | Path to the configuration file. |
| `ODC_LOG_LEVEL` | Logging level (debug, info, warn, error). |
| `ODC_PROFILE` | The active profile to use for the execution. |

## Using Correlation IDs for Diagnostics

Every execution of `odc` is assigned a unique **Correlation ID**. This ID is
printed to the log file and is essential for tracing a specific execution from
the CLI through to the Microsoft Graph API.

You can view the correlation ID for a command by enabling debug-level logging.

```bash
odc ls / --level debug
```

Check the `app.log` file in the `~/.local/state/odc/logs/` directory for detailed
information associated with a specific correlation ID.

## Non-Interactive Authentication

For automated systems and CI/CD pipelines, you can use **Client Secret** or
**Environment-based** authentication methods to avoid interactive prompts.

```bash
# Login using a service principal
odc auth login --method client-secret \
  --client-id <id> \
  --client-secret <secret> \
  --tenant-id <tenant>
```

## Next steps

- **[Architecture Overview](../../developer/explanation/architecture.md)**
- **[Configuration Management](manage-config.md)**

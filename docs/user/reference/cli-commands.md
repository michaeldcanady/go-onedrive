# CLI Commands Reference

This page provides a detailed reference for all `odc` commands, their flags, and
the path syntax used throughout the application.

## Path Syntax

`odc` supports multiple storage providers and drive aliases through path
prefixes.

- **OneDrive (Default):** Absolute paths starting with `/` refer to your active
  OneDrive drive (e.g., `/Documents/report.txt`). You can also use the
  `onedrive:` prefix.
- **Local Filesystem:** Paths starting with `local:` refer to your local machine
  (e.g., `local:/home/user/notes.txt`).
- **Drive Aliases:** Use a drive alias as a prefix to target a specific drive
  (e.g., `work-share:/Reports/january.pdf`).

## Standard Filesystem Commands

### `ls` - List files and directories
List the contents of a directory.

- **Usage:** `odc ls [PATH]`
- **Flags:**
    - `-r`, `--recursive`: List items recursively.
    - `-o`, `--format`: Output format (`short`, `long`, `json`, `yaml`, `tree`, `table`).
    - `-a`, `--all`: Show hidden items.
    - `--sort`: Sort items by field (`name`, `size`, `modified`).
    - `--desc`: Sort in descending order.
- **Examples:**
    - `odc ls /` (Lists root of active drive)
    - `odc ls local:./projects` (Lists local directory)
    - `odc ls -r -o tree /Documents` (Recursive tree listing)

### `mkdir` - Create a directory
Create a new folder in OneDrive or local filesystem.

- **Usage:** `odc mkdir [PATH]`
- **Examples:**
    - `odc mkdir /NewFolder`
    - `odc mkdir local:./new_local_folder`

### `touch` - Create a new file
Create a new, empty file or update the timestamp of an existing file.

- **Usage:** `odc touch [PATH]`

### `rm` - Remove files or directories
Delete files or folders.

- **Usage:** `odc rm [PATH]`
- **Flags:**
    - `-r`, `--recursive`: Remove directories recursively.
- **Examples:**
    - `odc rm /file.txt`
    - `odc rm -r /OldFolder`

### `cp` - Copy files
Copy files or directories between providers or within a provider.

- **Usage:** `odc cp [SOURCE] [DESTINATION]`
- **Flags:**
    - `-r`, `--recursive`: Copy directories recursively.
- **Examples:**
    - `odc cp local:file.txt /remote-copy.txt`
    - `odc cp /file1.txt /folder/file1.txt`

### `mv` - Move files
Move or rename a file or directory.

- **Usage:** `odc mv [SOURCE] [DESTINATION]`
- **Examples:**
    - `odc mv /old-name.txt /new-name.txt`
    - `odc mv /file.txt /MyFolder/file.txt`


### `cat` - Display file content
Print the content of a file to your terminal.

- **Usage:** `odc cat [PATH]`

---

## Data Transfer and Editing

### `upload` - Upload local files
Transfer a file or directory from your local machine to OneDrive. This is a
shortcut for `cp local:SOURCE onedrive:DESTINATION`.

- **Usage:** `odc upload [LOCAL_PATH] [REMOTE_PATH]`
- **Flags:**
    - `-r`, `--recursive`: Upload directories recursively.

### `download` - Download remote files
Transfer a file or directory from OneDrive to your local machine. This is a
shortcut for `cp onedrive:SOURCE local:DESTINATION`.

- **Usage:** `odc download [REMOTE_PATH] [LOCAL_PATH]`
- **Flags:**
    - `-r`, `--recursive`: Download directories recursively.

### `edit` - Edit a file in your local editor
Download a OneDrive file to a temporary location, open it with your local
editor, and automatically upload it back when you save and exit.

- **Usage:** `odc edit [REMOTE_PATH]`
- **Flags:**
    - `-f`, `--force`: Overwrite existing items if they exist.

---

## Authentication and Profile Management

### `auth` - Manage authentication
Manage your authentication session for the active profile.

- **Subcommands:**
    - `login`: Authenticate your current profile.
        - **Flags:**
            - `--method`: Auth method (`interactive`, `device-code`, `client-secret`, `environment`).
            - `--client-id`: Azure AD Application (Client) ID.
            - `--tenant-id`: Azure AD Tenant ID.
            - `--client-secret`: Client secret (for Service Principals).
            - `--show-token`: Print the access token to stdout.
    - `logout`: Clear the authentication state.

### `profile` - Manage account profiles
Profiles allow you to switch between multiple OneDrive accounts.

- **Subcommands:**
    - `create [NAME]`: Create a new profile.
    - `list`: List all available profiles.
    - `use [NAME]`: Set the active profile.
    - `delete [NAME]`: Delete a profile.
    - `current`: Display the name of the active profile.

---

## Drive Management

### `drive` - Manage drives
OneDrive accounts can have multiple drives (e.g., personal, shared SharePoint libraries).

- **Subcommands:**
    - `list`: List all drives you have access to.
    - `get [ID|ALIAS]`: Get information about a specific drive.
    - `use [ID|ALIAS]`: Set the default drive for the current session.
    - `alias`: Manage shortcuts for drives.
        - `list`: List all drive aliases.
        - `set [NAME] [ID]`: Create a shortcut for a drive.
        - `remove [NAME]`: Delete a drive alias.

---

## Configuration and Utilities

### `config` - Manage settings
Directly manage configuration keys for the active profile.

- **Subcommands:**
    - `get [KEY]`: View all or specific configuration settings.
    - `set [KEY] [VALUE]`: Update a specific configuration setting.

### `completion` - Generate completion script
Generate shell completion scripts for your environment.

- **Usage:** `odc completion [bash|zsh|fish|powershell]`

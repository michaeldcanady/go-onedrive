# CLI Commands Reference

This page provides a detailed reference for all `odc` commands and their flags.

## Standard Filesystem Commands

### `ls` - List files and directories
List the contents of your OneDrive directory.

- **Usage:** `odc ls [PATH]`
- **Flags:**
    - `-r`, `--recursive`: List items recursively.
    - `-o`, `--format`: Output format (`short`, `long`, `json`, `yaml`, `tree`, `table`).
    - `-a`, `--all`: Show hidden items.
    - `--sort`: Sort items by field (`name`, `size`, `modified`).
    - `--desc`: Sort in descending order.
- **Examples:**
    - `odc ls` (Lists root)
    - `odc ls MyFolder`
    - `odc ls -r --format tree` (Recursive tree listing)
    - `odc ls --sort name --desc` (Sort by name in descending order)

### `mkdir` - Create a directory
Create a new folder in your OneDrive.

- **Usage:** `odc mkdir [PATH]`
- **Examples:**
    - `odc mkdir NewFolder`

### `touch` - Create a new file
Create a new, empty file or update the timestamp of an existing file.

- **Usage:** `odc touch [PATH]`
- **Examples:**
    - `odc touch newfile.txt`

### `rm` - Remove files or directories
Delete files or folders from your OneDrive.

- **Usage:** `odc rm [PATH]`
- **Examples:**
    - `odc rm file.txt`

### `cp` - Copy files
Copy files or directories from one location to another.

- **Usage:** `odc cp [SOURCE] [DESTINATION]`
- **Flags:**
    - `-r`, `--recursive`: Copy directories recursively.
- **Examples:**
    - `odc cp local:file.txt onedrive:copy-of-file.txt`
    - `odc cp onedrive:file1.txt local:file1.txt`
    - `odc cp -r onedrive:Folder1 onedrive:Folder2`

### `mv` - Move files
Move or rename a file in your OneDrive.

- **Usage:** `odc mv [SOURCE] [DESTINATION]`
- **Examples:**
    - `odc mv old-name.txt new-name.txt`
    - `odc mv file.txt MyFolder/file.txt`


### `cat` - Concatenate and display file content
Print the content of a file to your terminal.

- **Usage:** `odc cat [PATH]`
- **Examples:**
    - `odc cat onedrive:notes.txt`

---

## Data Transfer and Editing

### `upload` - Upload local files
Transfer a file or directory from your local machine to OneDrive.

- **Usage:** `odc upload [LOCAL_PATH] [REMOTE_PATH]`
- **Flags:**
    - `-r`, `--recursive`: Upload directories recursively.
- **Examples:**
    - `odc upload local-data.csv RemoteFolder/`
    - `odc upload -r Projects/ RemoteProjects/`

### `download` - Download remote files
Transfer a file or directory from OneDrive to your local machine.

- **Usage:** `odc download [REMOTE_PATH] [LOCAL_PATH]`
- **Flags:**
    - `-r`, `--recursive`: Download directories recursively.
- **Examples:**
    - `odc download RemoteFolder/data.csv .`
    - `odc download -r RemoteProjects/ LocalProjects/`

### `edit` - Edit a file in your local editor
Download a OneDrive file to a temporary location, open it with your local editor, and automatically upload it back when you save and exit.

- **Usage:** `odc edit [REMOTE_PATH]`
- **Flags:**
    - `-f`, `--force`: Overwrite existing items if they exist.
- **Examples:**
    - `odc edit config.json`


---

## Authentication and Profile Management

### `auth` - Authenticate with OneDrive
Manage your authentication session.

- **Subcommands:**
    - `login`: Authenticate your current profile.
    - `logout`: Clear the authentication state.

### `profile` - Manage your account profiles
Profiles allow you to switch between multiple OneDrive accounts.

- **Subcommands:**
    - `create [NAME]`: Create a new profile.
    - `list`: List all available profiles.
    - `use [NAME]`: Set the active profile.
    - `delete [NAME]`: Delete a profile.

---

## Drive Management

### `drive` - Manage and browse different drives
OneDrive accounts can have multiple drives (e.g., your personal drive, shared SharePoint libraries).

- **Subcommands:**
    - `list`: List all drives you have access to.
    - `get`: Get information about a specific drive.
    - `use [DRIVE_ID]`: Set the default drive for the current session.
    - `alias`: Create and manage shortcuts (aliases) for drives.
        - `alias list`: List all drive aliases.
        - `alias set [NAME] [DRIVE_ID]`: Create a shortcut for a drive.
        - `alias remove [NAME]`: Delete a drive alias.

---

## Utility and Help

### `completion` - Generate completion script
Generate shell completion scripts for your environment.

- **Usage:** `odc completion [bash|zsh|fish|powershell]`
- **Examples:**
    - `source <(odc completion bash)`
    - `odc completion zsh > "${fpath[1]}/_odc"`


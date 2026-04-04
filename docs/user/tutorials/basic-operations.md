# Basic Operations

Once you've authenticated, you can start managing your OneDrive files using
standard Unix-style commands. This tutorial covers the most common file and
directory operations.

## Understanding Paths

When you use `odc`, you specify file and directory locations using paths. `odc`
supports multiple "providers" (like OneDrive and your local filesystem) and
uses prefixes to distinguish between them.

- **Default path:** Without a prefix, paths refer to your active OneDrive drive
  (e.g., `/Documents/report.txt`).
- **Explicit OneDrive path:** Use the `onedrive:` prefix (e.g.,
  `onedrive:/Documents/report.txt`).
- **Local path:** Use the `local:` prefix to refer to your local machine (e.g.,
  `local:/home/user/notes.txt`).
- **Drive aliases:** You can use a drive alias as a prefix to target a specific
  drive directly (e.g., `work-share:/Reports/january.pdf`).

> **Note:** Most `odc` commands expect absolute paths starting with `/` when
> referring to OneDrive items. For local paths, you can use absolute or
> relative paths after the `local:` prefix.

## Listing Files and Directories

The `ls` command lets you view the contents of a directory.

```bash
# List files in the root directory
odc ls /

# List files in a specific directory
odc ls /Documents

# View detailed information (long format)
odc ls -o long /Documents
```

### Powerful listing options
`odc` offers powerful flags to help you find exactly what you're looking for:

- **Recursive listing:** View all files and subdirectories.
  ```bash
  odc ls -r /Projects
  ```

- **Tree view:** View your directory structure in a visual tree format.
  ```bash
  odc ls -o tree /Projects
  ```

- **Sorting:** Organize your files by name, size, or modification date.
  ```bash
  odc ls --sort modified --desc /Documents
  ```

## Creating Files and Directories

Use `mkdir` to create directories and `touch` to create empty files or update
timestamps.

```bash
# Create a new directory
odc mkdir /Work

# Create an empty file
odc touch /Work/new_draft.txt
```

## Copying and Moving Items

Use `cp` to copy items and `mv` to move or rename them.

```bash
# Copy a file to another folder
odc cp /Work/new_draft.txt /Backup/old_draft.txt

# Rename a file
odc mv /Work/old_name.txt /Work/new_name.txt

# Move a file to another directory
odc mv /Work/new_name.txt /Documents/
```

## Deleting Items

Use the `rm` command to delete files or directories.

```bash
# Delete a file
odc rm /Work/old_draft.txt

# Delete a directory and all its contents
odc rm -r /OldProject
```

> **Warning:** Be careful when using `rm -r`. This command permanently deletes
> the directory and its contents without a confirmation prompt.

## Displaying File Content

Use `cat` to display the contents of a text file directly in your terminal.
This is useful for quickly checking the contents of small files.

```bash
odc cat /Documents/notes.txt
```

## Next steps

- **[Working with different drives](../how-to/work-with-drives.md)**
- **[Transferring files](../how-to/transfer-files.md)**
- **[Editing files natively](../how-to/edit-files.md)**

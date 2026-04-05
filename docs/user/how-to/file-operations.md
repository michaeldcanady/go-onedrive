# Manage Files and Folders

`odc` provides a set of standard filesystem commands that let you manage files
and directories in your OneDrive as if they were on your local machine. This
guide explores how to use these commands to organize and navigate your cloud
storage effectively.

## Listing Files and Directories

The `ls` command is the primary tool for viewing the contents of your OneDrive.
By default, it lists the root directory of your active drive.

```bash
# List the root directory
odc ls

# List a specific path
odc ls /Documents/Projects
```

### Choosing an output format
You can change how information is displayed using the `--format` (or `-o`) flag.

- **short (default):** A simple list of names.
- **long:** Detailed information including size, type, and modification date.
- **tree:** A hierarchical view of your directory structure.
- **json:** Machine-readable output, perfect for automation.
- **yaml:** Human-friendly structured data.

```bash
# View detailed information
odc ls -o long /Documents

# View a directory tree
odc ls -o tree /Projects
```

### Sorting and filtering
Customize the order and visibility of your files to find what you need quickly.

- **Sort by field:** Use `--sort` to order by `name`, `size`, or `modified`.
- **Reverse order:** Add the `--desc` flag to sort in descending order.
- **Show hidden items:** Use the `--all` (or `-a`) flag to include items that
  are typically hidden.

```bash
# Sort by size, largest first
odc ls --sort size --desc

# Show all items, including hidden ones
odc ls -a
```

### Recursive listing
To see all files in subdirectories, use the `--recursive` (or `-r`) flag.

> **Note:** Recursive mode is only supported with the `tree`, `long`, `json`,
> and `yaml` formats.

```bash
# List everything in a project folder
odc ls -r -o tree /Projects/MyApp
```

## Creating Items

### Create a directory
Use the `mkdir` command to create new folders.

```bash
odc mkdir /Documents/NewProject
```

### Create an empty file
Use the `touch` command to create a new, empty file. This is useful for
initializing log files or placeholders.

```bash
odc touch /Documents/notes.txt
```

## Organizing Items

### Move and rename
The `mv` command lets you rename items or move them to a different folder.

```bash
# Rename a file
odc mv old-name.txt new-name.txt

# Move a file to a folder
odc mv report.pdf /Documents/Archive/
```

### Copy items
Use the `cp` command to create a copy of a file in a new location.

```bash
odc cp template.docx /Projects/Proposal.docx
```

## Removing Items

Use the `rm` command to delete files or directories.

```bash
# Delete a single file
odc rm /Documents/temp.txt

# Delete a directory and all its contents
odc rm -r /OldProjects
```

> **Warning:** Removing an item is permanent and cannot be undone via `odc`.
> Always double-check your path before running `rm -r`.

## Next steps

- **[Transfer files between local and cloud](transfer-files.md)**
- **[Edit cloud files directly](edit-files.md)**
- **[Automate workflows with scripting](automation-and-scripting.md)**

# Perform file operations

`odc` provides a set of standard filesystem commands that allow you to 
manage files and directories in your OneDrive as if they were on your local 
machine.

## List files and directories

To view the contents of a directory, use the `ls` command.

```bash
odc ls [PATH]
```

By default, this lists the root directory of your active OneDrive.

### View a directory tree

To see a hierarchical view of your files and folders, use the `--format 
tree` flag.

```bash
odc ls --format tree
```

### List recursively

To see all files in subdirectories, use the `-r` or `--recursive` flag. 
This flag requires a compatible format such as `tree`, `long`, `json`, or 
`yaml`.

```bash
odc ls -r --format tree
```

### Sort output

You can sort the list by `name`, `size`, or `modified` using the 
`--sort` flag. Combine it with `--desc` for descending order.

```bash
odc ls --sort size --desc
```

## Create items

### Create a directory

To create a new folder, use the `mkdir` command.

```bash
odc mkdir NewFolder
```

### Create an empty file

To create a new, empty file, use the `touch` command.

```bash
odc touch NewFile.txt
```

## View file content

To print the content of a file to your terminal, use the `cat` command.

```bash
odc cat MyFile.txt
```

## Move and rename items

To rename a file or move it to a different folder, use the `mv` command.

```bash
odc mv OldName.txt NewName.txt
odc mv File.txt Folder/
```

## Copy items

To copy a file from one location to another, use the `cp` command.

```bash
odc cp Source.txt Destination.txt
```

## Remove items

To delete a file or directory, use the `rm` command.

```bash
odc rm FileToDelete.txt
```

> **Warning:** Removing an item is permanent and cannot be undone via 
> `odc`.

## Next steps

Now that you can manage your files, learn how to [transfer files](transfer-files.md) 
between your local machine and OneDrive or [edit files](edit-files.md) 
directly.

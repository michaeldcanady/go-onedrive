# Transfer files

Transferring files and directories between your local machine and OneDrive 
is done using the `upload` and `download` commands.

## Upload files and directories

To transfer an item from your local machine to OneDrive, use the 
`upload` command.

```bash
odc upload [LOCAL_PATH] [REMOTE_PATH]
```

### Upload a single file

```bash
odc upload backup.zip Backups/
```

### Upload a directory recursively

To upload a local directory and all of its contents, use the `-r` or 
`--recursive` flag.

```bash
odc upload -r Projects/ ActiveProjects/
```

## Download files and directories

To transfer an item from OneDrive to your local machine, use the 
`download` command.

```bash
odc download [REMOTE_PATH] [LOCAL_PATH]
```

### Download a single file

```bash
odc download "Documents/Meeting Notes.md" .
```

### Download a directory recursively

To download a remote directory and all of its contents, use the `-r` or 
`--recursive` flag.

```bash
odc download -r RemoteFolder/ ./LocalCopy/
```

## Next steps

After transferring your files, you can [manage them](file-operations.md) 
further or [edit them](edit-files.md) directly in your terminal.

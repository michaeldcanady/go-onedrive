# Transfer Files

Transferring files and directories between your local machine and OneDrive
is a core function of `odc`. Whether you're backing up local data or
retrieving work from the cloud, the `upload` and `download` commands
provide a simple and powerful way to manage your data.

## Upload to OneDrive

To move items from your local machine to OneDrive, use the `upload` command.
This is perfect for backing up your projects or sharing local files.

```bash
odc upload [LOCAL_PATH] [REMOTE_PATH]
```

### Upload a single file
To upload a single file, provide its local path and the target path in
OneDrive.

```bash
odc upload backup.zip /Backups/
```

### Upload a directory recursively
To upload an entire directory and all its contents, including subfolders,
use the `-r` (or `--recursive`) flag.

```bash
odc upload -r ./my_project /ActiveProjects/
```

## Download from OneDrive

To retrieve items from OneDrive and save them to your local machine, use the
`download` command.

```bash
odc download [REMOTE_PATH] [LOCAL_PATH]
```

### Download a single file
To download a file to your current local directory, use `.` for the
`LOCAL_PATH`.

```bash
odc download "/Documents/Meeting Notes.md" .
```

### Download a directory recursively
To download a remote directory and all of its contents, use the `-r`
(or `--recursive`) flag.

```bash
odc download -r /RemoteFolder/ ./LocalCopy/
```

## Best Practices for Data Transfer

- **Large transfers:** For very large folders or thousands of small files,
  use the `-r` flag with caution. Ensure you have a stable internet connection.
- **Handling conflicts:** By default, `odc` will prompt or fail if you try to
  upload or download over an existing file. Use the `--force` flag (if
  available) to overwrite without confirmation.
- **Paths with spaces:** If your file or folder names contain spaces, wrap the
  entire path in double quotes (e.g., `"/Documents/My Project/"`).
- **Relative vs. absolute:** For remote paths in OneDrive, it's always safer to
  use an absolute path starting with `/`.

## Next steps

- **[File operations](file-operations.md)**
- **[Automation and scripting](automation-and-scripting.md)**
- **[Native file editing](edit-files.md)**

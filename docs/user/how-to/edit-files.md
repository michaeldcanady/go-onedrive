# Edit files

The `edit` command allows you to modify OneDrive files using your local 
terminal editor (such as `vim`, `nano`, or the editor defined in your 
`$EDITOR` environment variable) without manually downloading and 
uploading.

## Edit a file

To edit a file, provide its path in OneDrive.

```bash
odc edit [REMOTE_PATH]
```

When you run this command:

1.  `odc` downloads the file to a temporary location on your machine.
2.  Your default editor opens the file.
3.  You make your changes and save the file.
4.  Once you exit the editor, `odc` automatically uploads the modified 
    content back to OneDrive.

## Overwrite existing content

If you need to ensure that your changes overwrite the remote file 
even if there are potential conflicts, use the `--force` flag.

```bash
odc edit config.json --force
```

## Configure your editor

`odc` respects the standard `$EDITOR` environment variable. If it's not 
set, `odc` defaults to common editors like `vim` or `nano` depending 
on your operating system.

To change your default editor, add the following to your shell profile 
(e.g., `~/.bashrc` or `~/.zshrc`):

```bash
export EDITOR='code --wait' # To use VS Code
```

## Next steps

After editing your files, you can [manage them](file-operations.md) or 
[transfer them](transfer-files.md) to other locations.

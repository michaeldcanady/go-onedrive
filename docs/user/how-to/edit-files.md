# Native file editing

The `edit` command in `odc` lets users you to modify OneDrive files using your
favorite local terminal editor (like `vim`, `nano`, or `code`). It streamlines
your workflow by handling the download and upload process automatically,
providing a seamless editing experience

## Edit a file

To edit a file, provide its path in OneDrive. `odc` will download it to a
temporary location and open it in your default editor

```bash
odc edit /Documents/config.json
```

### How it works
1.  **Download:** `odc` downloads the file to a secure temporary location
2.  **Edit:** Your terminal editor opens the file
3.  **Save & Exit:** Once you save your changes and exit the editor, `odc`
    automatically uploads the modified content back to OneDrive
4.  **Cleanup:** The temporary file is  removed from your local machine

## Handle conflicts

If you want to confirm that your changes overwrite the remote file even if
there have been changes since you started editing, use the `--force` flag

```bash
odc edit /Documents/config.json --force
```

## Configure your editor

`odc` respects your system's `$EDITOR` environment variable. If it's not set,
`odc` defaults to common terminal editors

To set your preferred editor, add the following line to your shell profile
(for example, `~/.bashrc` or `~/.zshrc`):

```bash
# Example: Use vs code as your default editor for odc
export EDITOR='code --wait'

# Example: Use vim as your default editor
export EDITOR='vim'
```

## Best practices

- **Save often:** While `odc` only uploads when you exit, saving frequently in
  your editor is still a good habit
- **Wait for upload:** Confirm you stay in your terminal until `odc` confirms the
  upload has been completed after you exit your editor
- **Large files:** For  large files, consider using `download` and `upload`
  separately for more control over the transfer process

## Next steps

- **[File operations](file-operations.md)**
- **[Transferring files](transfer-files.md)**
- **[Automation and scripting](automation-and-scripting.md)**

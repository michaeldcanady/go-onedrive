# Work with Drives

In `odc`, a drive represents a specific storage location within your account.
This can be your personal OneDrive, a shared library, or a SharePoint site.
`odc` makes it easy to manage multiple drives, switch between them, and create
aliases for faster access.

## List Available Drives

To see all the drives that your account has access to, use the `drive list`
command. This is often the first step when you're looking for a shared folder
or a SharePoint library.

```bash
odc drive list
```

The output includes each drive's name, ID, and its type (e.g., `personal`,
`business`).

## Get Drive Details

To view detailed information about a specific drive, including its storage
limits and ownership, use the `drive get` command followed by the drive's ID
or an alias.

```bash
odc drive get [DRIVE_ID]
```

This returns metadata such as the total storage quota, available space, and
owner information, which is useful for monitoring your storage usage.

## Set the Active Drive

By default, `odc` uses your primary OneDrive. If you frequently work with a
different drive (like a team SharePoint site), you can set it as the default
for your current session using the `drive use` command.

```bash
# Set the active drive using its ID
odc drive use [DRIVE_ID]
```

This command updates your current profile's configuration. Subsequent commands
like `ls` or `upload` will target this drive unless you switch back or use an
explicit flag.

## Create Drive Aliases

Drive IDs are long, complex, and difficult to remember. Aliases let you assign
a friendly, memorable name to a drive ID for faster access in your commands.

### Set an alias
To create or update an alias, use the `drive alias set` command.

```bash
odc drive alias set work-share [DRIVE_ID]
```

Now you can use `work-share` anywhere you would normally use the long
`DRIVE_ID`.

### List aliases
To view all your current drive aliases and the IDs they point to, use the
`drive alias list` command.

```bash
odc drive alias list
```

### Using Aliases in Paths
Once you've set an alias, you can use it as a prefix in any command that
accepts a path. This allows you to target a specific drive without switching
the active drive for your profile.

```bash
# List files in the 'work-share' drive directly
odc ls work-share:/Projects

# Copy a file from your personal drive to 'work-share'
odc cp /Reports/quarterly.pdf work-share:/Archive/
```

### Remove an alias
To delete an alias that you no longer need, use the `drive alias remove`
command.

```bash
odc drive alias remove work-share
```

## Next steps

- **[Basic operations](../tutorials/basic-operations.md)**
- **[Transferring files](transfer-files.md)**
- **[Managing profiles](manage-profiles.md)**

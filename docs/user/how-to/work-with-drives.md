# Work with drives

In `odc`, a drive represents a storage location. This can be your 
personal OneDrive, a shared library, or a SharePoint site. If you have 
access to multiple drives, you can list them, switch between them, and 
create aliases for easier access.

## List available drives

To see all the drives that your account can access, use the `drive list` 
command.

```bash
odc drive list
```

The output includes each drive's name, ID, and its type (e.g., `personal`, 
`business`).

## Get drive details

To view detailed information about a specific drive, use the `drive get` 
command followed by the drive's ID or an alias.

```bash
odc drive get [DRIVE_ID]
```

This returns metadata such as the total storage quota, available space, 
and owner information.

## Set the default drive

By default, `odc` uses your primary OneDrive. If you frequently work 
with a different drive, you can set it as the default for your profile 
using the `drive use` command.

```bash
odc drive use [DRIVE_ID]
```

This command updates your current profile's configuration. Subsequent 
commands like `ls` or `upload` will target this drive unless 
specified otherwise.

## Create drive aliases

Drive IDs are long and difficult to remember. Aliases let you assign a 
friendly name to a drive ID for faster access.

### Set an alias

To create or update an alias, use the `drive alias set` command.

```bash
odc drive alias set work-share [DRIVE_ID]
```

### List aliases

To view all your current drive aliases, use the `drive alias list` 
command.

```bash
odc drive alias list
```

### Remove an alias

To delete an alias that you no longer need, use the `drive alias remove` 
command.

```bash
odc drive alias remove work-share
```

## Next steps

After configuring your drives, you can begin [performing file 
operations](file-operations.md) or [transferring files](transfer-files.md).

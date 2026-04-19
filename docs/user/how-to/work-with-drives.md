# Work with Drives

In `odc`, a drive represents a specific storage location within your account.
This can be your personal OneDrive, a shared library, or a SharePoint site.
`odc` makes it easy to discover available drives and mount them to your
virtual filesystem for easy access.

## Discover Available Drives

To see all the drives that your account has access to, use the
`backend-discovery list` command. This is often the first step when you're
looking for a shared folder or a SharePoint library.

```bash
odc backend-discovery list
```

The output includes each drive's name and ID. Your primary personal drive is
marked with an asterisk (*).

> **Note:** `drive` is an alias for `backend-discovery`, so you can also use
> `odc drive list`.

## Get Personal Drive Details

To view information about your primary personal OneDrive drive, use the
`backend-discovery get` command.

```bash
odc backend-discovery get
```

This returns the name and ID of your personal drive.

## Manage Mount Points

Mount points allow you to map a specific OneDrive drive or local directory to
a path in your virtual filesystem. This provides a consistent way to access
different storage locations.

### List mount points
To view all your current mount points, use the `mount list` command.

```bash
odc mount list
```

By default, `odc` mounts your local filesystem at `/local` and your personal
OneDrive at `/onedrive`.

### Add a mount point
To add a new mount point, use the `mount add` command. You can mount another
OneDrive drive by specifying its drive ID.

```bash
odc mount add /work --type onedrive --drive-id [DRIVE_ID]
```

Now you can access this drive using the `/work` path or the `work:` prefix.

### Using Mount Points in Paths
Once you've added a mount point, you can use it as a prefix or an absolute
path in any command.

```bash
# List files in the 'work' mount point using absolute path
odc ls /work/Projects

# List files using the mount prefix
odc ls work:/Projects

# Copy a file from your personal drive to 'work'
odc cp /onedrive/Reports/quarterly.pdf /work/Archive/
```

### Remove a mount point
To remove a mount point that you no longer need, use the `mount remove`
command.

```bash
odc mount remove /work
```

## Next steps

- **[Basic operations](../tutorials/basic-operations.md)**
- **[Transferring files](transfer-files.md)**
- **[Managing profiles](manage-profiles.md)**

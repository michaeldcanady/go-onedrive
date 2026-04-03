# Manage profiles

Profiles in `odc` allow you to manage multiple OneDrive accounts or 
configurations on the same machine. Each profile stores its own 
authentication state, drive preferences, and configuration settings.

## List profiles

To see all available profiles, use the `profile list` command.

```bash
odc profile list
```

The output displays the names of all created profiles. The active profile 
is typically highlighted or indicated in the list.

## View the current profile

To check which profile is currently active in your session, use the 
`profile current` command.

```bash
odc profile current
```

This returns the name of the active profile that `odc` uses for all 
subsequent commands unless overridden by the `--profile` flag.

## Create a profile

To add a new account or configuration, use the `profile create` command 
followed by a unique name.

```bash
odc profile create my-work-account
```

After creating a profile, you must authenticate it using the 
`odc auth login` command before you can access OneDrive files.

## Switch between profiles

To change the active profile for your environment, use the `profile use` 
command.

```bash
odc profile use my-work-account
```

This command updates your persistent state so that future `odc` 
commands use the specified profile by default.

## Delete a profile

To remove a profile and its associated state, use the `profile delete` 
command.

> **Warning:** Deleting a profile permanently removes its authentication 
> tokens and configuration. You must re-authenticate if you recreate the 
> profile later.

```bash
odc profile delete old-account
```

## Next steps

Now that you know how to manage profiles, you can [authenticate](authenticate.md) 
your new profiles or start [working with drives](work-with-drives.md).

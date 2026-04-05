# Authenticating with odc

To interact with OneDrive, you need to authenticate your profiles. `odc` supports several authentication methods.

## Standard Login

The most common way to authenticate is with the `auth login` command. By default, it uses the **Interactive Browser** flow.

```bash
odc auth login
```

This will:

1.  Open your default web browser.
2.  Ask you to sign in to your Microsoft account.
3.  Request permissions for `odc` to access your files.
4.  Once completed, the browser window will close, and your session will be authenticated.

## Authentication Methods

You can specify a different authentication method using the `--method` flag:

| Method          | Description                                                                                                                      |
| :-------------- | :------------------------------------------------------------------------------------------------------------------------------- |
| `interactive`   | (Default) Opens a web browser for authentication.                                                                                |
| `device-code`   | Provides a code to enter at [https://microsoft.com/devicelogin](https://microsoft.com/devicelogin). Useful for headless servers. |
| `client-secret` | For use with Azure Service Principals (requires `client-id`, `client-secret`, and `tenant-id`).                                  |
| `environment`   | Reads credentials from standard Azure environment variables.                                                                     |

Example using device code:

```bash
odc auth login --method device-code
```

## Logging Out

To clear the authentication state for your current profile:

```bash
odc auth logout
```

## How Authentication Works

`odc` uses the [Azure Identity](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity) SDK under the hood. Tokens are cached locally in your profile's state, allowing you to run commands without re-authenticating every time.

## Next steps

Once you are authenticated, you can [manage your profiles](manage-profiles.md) or 
begin [working with drives](work-with-drives.md).


# Manage configuration

`odc` allows you to manage its configuration directly from the command 
line. This is useful for updating your Microsoft Graph Client ID, tenant, 
or redirect URIs without manually editing YAML files.

## View configuration

To see the complete configuration for your active profile, use the 
`config get` command.

```bash
odc config get
```

### View a specific setting

You can also retrieve a specific configuration value by providing its 
key.

```bash
odc config get auth.client_id
```

## Update configuration

To update a configuration setting, use the `config set` command 
followed by the key and the new value.

```bash
odc config set auth.redirect_uri "http://localhost:9000"
```

If your profile does not already have a configuration file, `odc` will 
automatically create one in its default configuration directory and 
update your profile metadata.

## Available configuration keys

The following keys are available for the `microsoft` provider:

| Key                  | Description                                      |
| :------------------- | :----------------------------------------------- |
| `auth.provider`      | The identity provider (e.g., `microsoft`).      |
| `auth.client_id`     | Your Azure AD Application (Client) ID.          |
| `auth.tenant_id`     | Your Azure AD Tenant ID (or `common`).          |
| `auth.client_secret` | The client secret for Service Principals.       |
| `auth.method`        | Default auth method (`interactive`, `device-code`).|
| `auth.redirect_uri`  | The URI used for interactive browser login.     |

## Configuration Schema

If you prefer to edit your configuration manually, we provide a JSON schema 
to ensure your `config.yaml` is valid. You can use this schema in editors like 
VS Code to get validation and autocompletion.

To use the schema in VS Code, add the following to your `settings.json`:

```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/michaeldcanady/go-onedrive/main/internal/config/schema.json": "config.yaml"
  }
}
```

## Next steps

After updating your configuration, you may need to [re-authenticate](authenticate.md) 
your profile for the changes to take effect.

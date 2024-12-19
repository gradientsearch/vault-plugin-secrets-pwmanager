# vault-password-manager-plugin
🔐 Password Manager Plugin for Vault: Vault-on-Vault Pattern


## Definitions

- **safe** a safe is a isolated storage location for secrets. This plugin uses a KV secret mount per safe.

## Development Setup

> [!NOTE]
> This will create everything you need including a test user `stephen` and password `hashicorp`
> in the `userpass` auth mount.
>

For now we simply just need to run the following commands:

```
make build & 
make configure
```
## User API

This plugin has a simple user API including a `safe`  endpoint. The `safe` endpoint is exposed to users allowing them to manage their safes.

## Admin API

The `config` endpoint is only available to Vault admins and configures the plugin with necessary credentials to manage safes. The conig consist a `role_id` and a `secret_id` allowing `pwmgr` to retrieve a Vault token to manage `safes` and `policies`.


## Plugin

The plugin `authorizes` users to perform CRUD requests against safes. It does this by managing Vault policies for each user.
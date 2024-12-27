# kv secret mount pwmgr uses as a data store
path "/sys/mounts/vaults/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

path "/sys/policies/acl/pwmanager/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

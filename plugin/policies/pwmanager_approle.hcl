# kv secret mount pwmgr uses as a data store
path "/sys/mounts/bundles/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

path "/sys/policies/acl/pwmanager/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

path "identity/entity/id/+" {
    capabilities = ["read"]
}
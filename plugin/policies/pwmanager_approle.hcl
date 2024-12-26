# kv secret mount pwmgr uses as a data store
path "pwmanager/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}
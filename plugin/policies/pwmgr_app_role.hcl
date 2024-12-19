# kv secret mount pwmgr uses as a data store
path "pwmgr/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}


# AppRole has full access to the policies under /pwmgr
# Users will have a default policy allowing access to /pwmgr/{identity.entity.name}
path "sys/policy/pwmgr/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

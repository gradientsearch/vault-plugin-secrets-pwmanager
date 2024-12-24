# kv secret mount pwManager uses as a data store
path "pwManager/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}


# AppRole has full access to the policies under /pwManager
# Users will have a default policy allowing access to /pwManager/{identity.entity.name}
path "sys/policy/pwManager/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

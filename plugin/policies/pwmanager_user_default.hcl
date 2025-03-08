# <mount>/<entity-id>
path "pwmanager/register" {
    capabilities = ["create"]
}

path "pwmanager/users/{{ identity.entity.id }}" {
    capabilities = ["update", "read"]
}

path "pwmanager/users/{{ identity.entity.id }}" {
    capabilities = ["update", "read"]
}

path "bundles/data/{{ identity.entity.id }}/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

path "bundles/metadata/{{ identity.entity.id }}/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

path "pwmanager/bundles" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

path "pwmanager/bundles/+/+/users/+" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}

// User needs to know what their entity name is. 
path "identity/entity/id/{{ identity.entity.id }}" {
    capabilities = ["read"]
}
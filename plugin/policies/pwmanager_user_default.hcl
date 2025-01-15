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

path "bundles/{{ identity.entity.id }}/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}
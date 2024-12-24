# <mount>/<entity-id>
path "pwmanager/register" {
    capabilities = ["create"]
}

path "pwmanager/users/{{ identity.entity.id }}" {
    capabilities = ["update", "read"]
}
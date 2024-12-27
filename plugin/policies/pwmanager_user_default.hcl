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

// users default private vault
// naming conventions for vaults: vault/{{ identity.entity.id }}/<uuid>
// private vault is the exception to the rule since we create their vault.
// maybe we can create this vault with a uuid on the first secret they create
// in their private vault via the client... This will do for now.
path "vault/{{ identity.entity.id }}/private" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}
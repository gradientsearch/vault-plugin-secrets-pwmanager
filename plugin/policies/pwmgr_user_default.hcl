# <mount>/<entity-id>
path "pwmanager/{{identity.entity.id}}/*" {
    capabilities = ["create", "read", "update", "patch", "delete", "list"]
}
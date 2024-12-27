path "pwmanager/users" {
    capabilities = ["list"]
}

path "pwmanager/users/*" {
    capabilities = ["delete"]
}
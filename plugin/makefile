# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

build: plugin run

plugin:
	go build -o vault/plugins/vault-password-manager-plugin cmd/vault-password-manager-plugin/main.go

run:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

enable:
	vault secrets enable -path=pwmgr vault-password-manager-plugin

appRole:
	vault auth enable approle && \
	vault write auth/approle/role/my-role \
    token_type=batch \
    secret_id_ttl=10m \
    token_ttl=20m \
    token_max_ttl=30m \
    secret_id_num_uses=40

readConfig:
	vault read pwmgr/config

.SILENT:
config:
	export ROLE_ID=$(shell vault read -field=role_id  auth/approle/role/my-role/role-id ) \
	export SECRET_ID=$(shell vault write -f -field=secret_id auth/approle/role/my-role/secret-id) && \
	vault write pwmgr/config role_id=$$ROLE_ID secret_id=$$SECRET_ID url="http://localhost:8200"
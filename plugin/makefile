# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

build: plugin run

debug: plugin  run-debug

plugin:
	go build -o vault/plugins/local/pwmanager cmd/vault-plugin-secrets-pwmanager/main.go

run:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins/local -log-level=debug

run-debug:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level=debug 2>&1 | grep pwManager

configure: enable appRole config secretMount userpass entity pwManager

enable:
	vault secrets enable -path=pwManager pwmanager

appRole:
	vault policy write plugins/pwManager-app-role policies/pwManager_app_role.hcl && \
	vault auth enable approle && \
	vault write auth/approle/role/pwManager \
    token_type=batch \
    secret_id_ttl=10m \
    token_ttl=60m \
    token_max_ttl=120m \
    secret_id_num_uses=0 \
	token_policies="plugins/pwManager-app-role"

readAppRole:
	vault read auth/approle/role/pwManager

readConfig:
	vault read pwManager/config

#.SILENT:
config:
	vault write pwManager/config role_id=$(shell vault read -field=role_id  auth/approle/role/pwManager/role-id ) \
	secret_id=$(shell vault write -f -field=secret_id auth/approle/role/pwManager/secret-id) \
	url="http://localhost:8200"

secretMount:
	vault secrets enable -version=2 -path=pwManagerkv kv

userpass:
	vault auth enable userpass 
	vault write auth/userpass/users/stephen password=hashicorp token_policies="plugins/pwManager-user-default,pwManager/entity/stephen"

login:
	vault write -format json auth/userpass/login/stephen password=hashicorp

entity:
	vault write /identity/entity/id/$(shell vault write -format json auth/userpass/login/stephen password=hashicorp | jq -r '.auth.entity_id' | tr -d "\n") name=stephen
	vault policy write plugins/pwManager-user-default policies/pwManager_user_default.hcl

pwManager:
	vault write pwManager/role/stephen username=dontneed



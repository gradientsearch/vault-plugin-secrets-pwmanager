# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

build: plugin run

debug: plugin  run-debug

plugin:
	go build -o vault/plugins/pwmanager cmd/vault-plugin-secrets-pwmanager/main.go

run:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level=debug

run-debug:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level=debug 2>&1 | grep pwmgr

configure: enable appRole config secretMount userpass entity pwmgr

enable:
	vault secrets enable -path=pwmgr pwmanager

appRole:
	vault policy write plugins/pwmgr-app-role policies/pwmgr_app_role.hcl && \
	vault auth enable approle && \
	vault write auth/approle/role/pwmgr \
    token_type=batch \
    secret_id_ttl=10m \
    token_ttl=60m \
    token_max_ttl=120m \
    secret_id_num_uses=0 \
	token_policies="plugins/pwmgr-app-role"

readAppRole:
	vault read auth/approle/role/pwmgr

readConfig:
	vault read pwmgr/config

#.SILENT:
config:
	vault write pwmgr/config role_id=$(shell vault read -field=role_id  auth/approle/role/pwmgr/role-id ) \
	secret_id=$(shell vault write -f -field=secret_id auth/approle/role/pwmgr/secret-id) \
	url="http://localhost:8200"

secretMount:
	vault secrets enable -version=2 -path=pwmgrkv kv

userpass:
	vault auth enable userpass 
	vault write auth/userpass/users/stephen password=hashicorp token_policies="plugins/pwmgr-user-default,pwmgr/entity/stephen"

login:
	vault write -format json auth/userpass/login/stephen password=hashicorp

entity:
	vault write /identity/entity/id/$(shell vault write -format json auth/userpass/login/stephen password=hashicorp | jq -r '.auth.entity_id' | tr -d "\n") name=stephen
	vault policy write plugins/pwmgr-user-default policies/pwmgr_user_default.hcl

pwmgr:
	vault write pwmgr/role/stephen username=dontneed



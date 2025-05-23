# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)
PWD = $(shell pwd)

setup:
	export VAULT_ADDR=http://127.0.0.1:8200 && \
	export VAULT_TOKEN=root && \
	$(MAKE) dev

dev: build-plugin debug sleep unseal sleep register enable appRole config secretMount userpass entity cors

# Kill the vault server. Run before running setup again.
kill:
	$(shell kill $$(ps aux | grep 'vault server' | awk '{print $2}'))

sleep:
	sleep 5

build-plugin:
	go clean -testcache && \
	go clean && \
	go build -o vault/plugins/local/pwmanager cmd/vault-plugin-secrets-pwmanager/main.go

debug:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins/local -log-level=debug &


unseal:
	$(eval KEY := $(shell cat ${PWD}/vault/json | jq -r '.unseal_keys_b64 | first ' | tr -d '\n'))
	vault operator unseal $(KEY)

token:
	cat ${PWD}/vault/json | jq -r '.root_token' | tr -d '\n' > ${PWD}/vault/dev_root_token.txt
	export VAULT_TOKEN=$(shell cat vault/dev_root_token.txt)

register:
	vault plugin register -sha256=$(shell sha256sum ./vault/plugins/local/pwmanager | cut -d ' ' -f1 | tr -d '\n') -version=v0.1.0 secret pwmanager

unregister:
	vault plugin deregister -version=v0.1.0 secret pwmanager

enable:
	vault secrets enable -path=pwmanager pwmanager

appRole:
	vault policy write pwmanager/approle policies/pwmanager_approle.hcl && \
	vault auth enable approle && \
	vault write auth/approle/role/pwmanager \
    token_type=batch \
    secret_id_ttl=0 \
    token_ttl=60m \
    token_max_ttl=120m \
    secret_id_num_uses=0 \
	token_policies="pwmanager/approle"

readAppRole:
	vault read auth/approle/role/pwmanager

readConfig:
	vault read pwmanager/config

cors:
	vault write /sys/config/cors enabled=true allowed_origins='*' allowed_headers='Content-Type,X-Requested-With,X-Vault-AWS-IAM-Server-ID,X-Vault-No-Request-Forwarding,X-Vault-Token,Authorization,X-Vault-Wrap-Format,X-Vault-Wrap-TTL'

#.SILENT:
config:
	vault write pwmanager/config role_id=$(shell vault read -field=role_id  auth/approle/role/pwmanager/role-id ) \
	secret_id=$(shell vault write -f -field=secret_id auth/approle/role/pwmanager/secret-id) \
	url="localhost:8200"

secretMount:
	vault secrets enable -version=2 -path=bundles kv
	vault write bundles/config cas_required=true

userpass:
	vault auth enable userpass
	vault write auth/userpass/users/stephen password=hashicorp token_policies="plugins/pwmanager-user-default,pwmanager/entity/stephen" token_ttl="24h"
	vault write auth/userpass/users/bob password=hashicorp token_policies="plugins/pwmanager-user-default,pwmanager/entity/bob" token_ttl="24h"
	vault write auth/userpass/users/alice password=hashicorp token_policies="plugins/pwmanager-user-default,pwmanager/entity/alice" token_ttl="24h"

login-bob:
	vault write -format json auth/userpass/login/bob password=hashicorp

login-alice:
	vault write -format json auth/userpass/login/alice password=hashicorp

login:
	vault write -format json auth/userpass/login/stephen password=hashicorp

entity:
	vault write /identity/entity/id/$(shell vault write -format json auth/userpass/login/stephen password=hashicorp | jq -r '.auth.entity_id' | tr -d "\n") name=stephen
	vault write /identity/entity/id/$(shell vault write -format json auth/userpass/login/bob password=hashicorp | jq -r '.auth.entity_id' | tr -d "\n") name=bob
	vault write /identity/entity/id/$(shell vault write -format json auth/userpass/login/alice password=hashicorp | jq -r '.auth.entity_id' | tr -d "\n") name=alice
	vault policy write plugins/pwmanager-user-default policies/pwmanager_user_default.hcl

# set up vault with file storage for client development
# root token is stored in ./vault/dev_root_token.txt
# export VAULT_TOKEN=$(cat vault/dev_root_token.txt)
# TODO: need to fix plugin not found error when using hashicorp/vault image https://github.com/hashicorp/vault/issues/26550
client: docker-run unseal token

docker-run:
	mkdir -p vault/data && \
	chmod -R 777 vault && \
	docker run -d --cap-add=IPC_LOCK -e 'VAULT_LOCAL_CONFIG={"api_addr":"http://0.0.0.0:8200", "plugin_directory": "/plugins", "storage": {"file": {"path": "/data/file"}}, "listener": [{"tcp": { "address": "0.0.0.0:8200", "tls_disable": true}}], "default_lease_ttl": "168h", "max_lease_ttl": "720h", "ui": true}' -p 8200:8200 -v ${PWD}/vault/plugins/local:/plugins  -v ${PWD}/vault/data:/data --name pwmanager-vault hashicorp/vault  server && \
	sleep 5 && \
	vault operator init -key-shares=1 -key-threshold=1 -format=json > ${PWD}/vault/json

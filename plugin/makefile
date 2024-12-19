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

config:
	vault write pwmgr/config username=pwmgr password=secret url="http://localhost:8200"

readConfig:
	vault read pwmgr/config
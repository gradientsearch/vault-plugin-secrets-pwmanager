FROM alpine:latest

WORKDIR /app

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o vault/plugins/pwmanager cmd/vault-plugin-secrets-pwmanager/main.go


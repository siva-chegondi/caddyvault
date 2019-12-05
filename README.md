# CaddyVault

A TLS Clustering plugin of caddyserver to use [Vault](https://vaultproject.io) as TLS storage.

state: **ALPHA**

## Building with CaddyVault plugin
To extend caddy with CaddyVault plugin, you need to append following import statements
in github.com/caddyserver/caddy/caddy/caddymain/run.go file.
```
import (
   _ "github.com/caddyserver/caddy/caddyhttp"
   _ "github.com/siva-chegondi/caddyvault"
)
```

## Configuration
We are using Vault's official Go Client API.
Use following ENV variables to config plugin.

*  VAULT_ADDR
*  VAULT_TOKEN


## Docker file
Reference for building docker image
[CaddyVault Docker](https://github.com/siva-chegondi/caddyvault-docker)

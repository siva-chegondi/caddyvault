# CaddyVault

A TLS clustering plugin for caddyserver to use [Vault](https://vaultproject.io) as storage for storing TLS data like certificates, keys etc.,

state: **ALPHA**

## Prerequisite
This plugin expects the following environment. 
* You need a VAULT server running and accessible from the machine/s on which caddy is running.

## Extending caddy with CaddyVault plugin
To extend caddy with CaddyVault plugin, we need to include following `import statement`
in github.com/mholt/caddy/caddy/caddymain/run.go file.
```
import (
   _ "github.com/mholt/caddy/caddyhttp"
   _ "github.com/siva-chegondi/caddyvault"
)
```

## Docker file

Checkout following project for reference to build your own docker file.
[CaddyVault Docker](https://github.com/siva-chegondi/caddyvault-docker)

## Configuration

### Vault configuration
* We need to enable KV2 secrets engine on the path `certpaths`.

### Caddy configuration
* We can enable `CaddyVault` plugin by setting environment variable `CADDY_CLUSTERING` to `vault`.
* Now set the following environment variables.
   
    * CADDY_CLUSTERING_VAULT_ENDPOINT
    * CADDY_CLUSTERING_VAULT_KEY

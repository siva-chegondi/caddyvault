# Caddyvault

A TLS clustering plugin for caddyserver to use [Vault](https://vaultproject.io) as storage for storing TLS data like certificates, keys etc.,

state: **ALPHA**

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

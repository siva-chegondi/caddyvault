# Caddyvault

A TLS plugin of caddyserver to store certificates in hashicorp Vault in state of Clustering.

**State: [ InProgress ]**

## Extending caddy with this plugin
I configured the caddy source code to extend `vault` plugin, you just clone the following repo and install the project to have binary in your GO environment as follows.

[Caddy - Vault Plugin](https://github.com/siva-chegondi/caddy)

```
cd $GOPATH/src/github.com/siva-chegondi/caddy
go install
```

Or Simply you can do following

```
go get github.com/siva-chegondi/caddy
```

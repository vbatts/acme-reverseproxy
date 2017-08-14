# acme-reverseproxy

A multi-domain, TLS, reverse proxy that uses [Let's Encrypt](https://letsencrypt.org/) as the automatic CA.

## Install

```shell
go get github.com/vbatts/acme-reverseproxy
```

## Configure

To get started with a configuration file, do:

```shell
acme-reverseproxy gen config > config.toml
```

Then edit as needed.

## Usage

This uses the default listener on `:https`/`:443` so it will need privilege or to be inside a container for port mapping.
```shell
acme-reverseproxy srv --config ./config.toml
```

## Good to know

As the certificates for the domains are issued once validated, the domains configured ought to be public facing so that Let's Encrypt can attest it.

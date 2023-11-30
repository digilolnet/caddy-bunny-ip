# Bunny.net Edge IP source module for Caddy

> Retrieves Bunny.net Edge IPs for use in Caddy `trusted_proxies` directives.

[![Go report](https://goreportcard.com/badge/github.com/digilolnet/caddy-bunny-ip)](https://goreportcard.com/report/github.com/digilolnet/caddy-bunny-ip)
[![GoDoc](https://godoc.org/github.com/digilolnet/caddy-bunny-ip?status.svg)](https://godoc.org/github.com/digilolnet/caddy-bunny-ip)
[![License](https://img.shields.io/github/license/digilolnet/caddy-bunny-ip.svg)](https://github.com/digilolnet/caddy-bunny-ip/blob/master/LICENSE.txt)
[![Code with hearth by Stnby](https://img.shields.io/badge/%3C%2F%3E%20with%20%E2%99%A5%20by-Stnby-ff1414.svg)](https://github.com/stnby)

## Caddy module name
```
http.ip_sources.bunny
```

## Config example
Put following config in global options under corresponding server options:
```
trusted_proxies bunny {
    interval 12h
    timeout 15s
}
```

## License
This project is licensed under the Apache License, Version 2.0 - see the [LICENSE.txt](https://github.com/digilolnet/caddy-bunny-ip/blob/master/LICENSE.txt) file for details.

This project is based on [caddy-cloudflare-ip](https://github.com/WeidiDeng/caddy-cloudflare-ip) module. Thanks [WeidiDeng](https://github.com/WeidiDeng) for your hard work.

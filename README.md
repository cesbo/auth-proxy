# Auth Proxy

Auth Proxy for HTTP Authorization with multiply backends.

## Install

Save binary file from [Releases](https://github.com/cesbo/auth-proxy/releases) to `/usr/local/bin/auth-proxy` on your server.

Set permission to execute file: `chmod +x /usr/local/bin/auth-proxy`

## Config

Create configuration file `/etc/auth-proxy.conf` on your server:

```json
{
    "listen": ":1064",
    "backend": [
        "http://mw1.example.com",
        "http://mw2.example.com"
    ]
}
```

- `listen` - address and port for built-in HTTP server. For example: `127.0.0.1:1064` to listen port 1064 on localhost, `:1064` to listen port 1064 on any interface.
- `backend` - list of the backend URLs

## Astra Configuration

To configure HTTP Authorization in [Astra](https://www.cesbo.com) open Settings -> HTTP Auth.
Select `HTTP Request` in the `Backend Type` and specify URL to the Auth Proxy: `http://127.0.0.1:1064`.

## Backend URL

Backend URL depends on the middleware being used.

### Ministra / Stalker

Backend URL:

```
http://example.com/stalker_portal/server/api/chk_flussonic_tmp_link.php
```

In the Ministra / Stalker settings turn on option "Temporary URL - Flussonic support"

### IPTVPORTAL

Backend URL:

```
https://go.iptvportal.cloud/auth/arescrypt/
```

In the portal settings open "Keys" menu and create a new key:

- Name: Astra
- Algorithm: ARESSTREAM
- Mode: SM
- Key Length: 1472 bit
- Update Rate: 1:00:00

In channel settings:

- Auth: arescrypt
- Encoded: turn on
- Key: Astra

### Microimpulse Smarty

Backend URL:

```
http://example.com/tvmiddleware/api/streamservice/token/check/
```

## Auto Start

To start Auth Proxy with system save to the file `/etc/systemd/system/auth-proxy.service` next config:

```ini
[Unit]
Description=Auth Proxy service
After=network-online.target
StartLimitBurst=5
StartLimitIntervalSec=10

[Service]
Type=simple
TimeoutStartSec=10
LimitNOFILE=65536
ExecStart=/usr/local/bin/auth-proxy /etc/auth-proxy.conf
KillMode=process
Restart=on-failure
RestartSec=1

[Install]
WantedBy=multi-user.target
```

Command to manage service:

| Command                        | Description                      |
|--------------------------------|----------------------------------|
| `systemctl restart auth-proxy` | Restart service                  |
| `systemctl start auth-proxy`   | Start service                    |
| `systemctl stop auth-proxy`    | Stop service                     |
| `systemctl status auth-proxy`  | Service status                   |
| `systemctl enable auth-proxy`  | Launch service on system startup |
| `systemctl disable auth-proxy` | Disable autorun                  |
| `journalctl -fu auth-proxy`    | Service logs                     |

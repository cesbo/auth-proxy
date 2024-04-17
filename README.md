# Auth Proxy

Auth Proxy for HTTP authorization with multiply backends.

## Config

```json
{
    "listen": ":1064",
    "backend": [
        "http://mw1.example.com",
        "http://mw2.example.com"
    ]
}
```

- `listen` - address and port for built-in HTTP server
- `backend` - list of the backend URLs

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

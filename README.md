# mocksmtp

`mocksmtp` is a simple tool to develop and test mail sending features.

## Features

- SMTP-Server to send Mails to
- HTTP-Server to investigate, what was sent
- Cleanup mails after some time

## Quickstart

```bash
# Help
mocksmtp -help

# HTTP to the loopback-interface port 8080 and smtp to all interfaces port 8025
mocksmtp -http-bind=127.0.0.1:8080 -smtp-bind=:8025

# Set retention time to 15min
mocksmtp -rentention-time=15m
```


## (Possible) further features

- Simple "fake" authentication
- Save messages for integration into a CI
- API to integrate into CI
- Persist mails (?)

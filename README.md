# flared
Cloudflare DDNS utility

## Installation

We recommend using `docker` install. Consider the following `compose.yml` example as a starting point:

```yaml
version: '3'
services:
  flared:
    container_name: flared
    image: ghcr.io/boggydigital/flared:latest
    environment:
      # - FD_TOKEN=(CLOUDFLARE DNS API TOKEN)
      # - FD_SYNC_FILENAME=/var/lib/flared/domains.txt
    volumes:
      # state
      - /docker/flared:/var/lib/flared
      # logs
      - /docker/flared:/var/log/flared
      # sharing timezone from the host
      - /etc/localtime:/etc/localtime:ro
    ports:
      # https://en.wikipedia.org/wiki/Galileo_Galilei
      - "1564:1564"
    restart: always
```

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
      # - FD_SYNC_FILENAME=domains.txt
    volumes:
      # input
      - /docker/flared:/var/lib/flared/input
      # metadata
      - /docker/flared/metadata:/var/log/flared/metadata
      # sharing timezone from the host
      - /etc/localtime:/etc/localtime:ro
    ports:
      # https://en.wikipedia.org/wiki/Galileo_Galilei
      - "1564:1564"
    restart: always
```

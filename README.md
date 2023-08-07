# cf_ddns
Cloudflare DDNS utility

## Installation

We recommend using `docker` install. Consider the following `compose.yml` example as a starting point:

```yaml
version: '3'
services:
  cf_ddns:
    container_name: cf_ddns
    image: ghcr.io/boggydigital/cf_ddns:latest
    environment:
      # - CF_DDNS_TOKEN=(CLOUDFLARE DNS API TOKEN)
      # - CF_DDNS_SYNC_FILENAME=/var/lib/cf_ddns/domains.txt
    volumes:
      # state
      - /docker/cf_ddns:/var/lib/cf_ddns
      # logs
      - /docker/cf_ddns:/var/log/cf_ddns
      # sharing timezone from the host
      - /etc/localtime:/etc/localtime:ro
    ports:
      # https://en.wikipedia.org/wiki/Galileo_Galilei
      - "1564:1564"
    restart: always
```

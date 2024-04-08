# vidlp - Free Video Downloader

## Overview

**vidlp** is a robust video downloader, constructed on the [yt-dlp](https://github.com/yt-dlp/yt-dlp) library, It enables you to download videos from a variety of platforms.

[Preview](https://dl.sunls.de)

## Deploy

### Docker

```shell
docker run --name vidlp -d --restart unless-stopped -e 'HOST=0.0.0.0' -p 3003:3003 sunls24/vidlp
```

#### cookies.txt

**Some platforms require cookies to download clearer videos.** [Look here.](https://github.com/yt-dlp/yt-dlp/wiki/FAQ#how-do-i-pass-cookies-to-yt-dlp)

```shell
touch cookies.txt
docker run --name vidlp -d --restart unless-stopped -e 'HOST=0.0.0.0' -v $(pwd)/cookies.txt:/yt-dlp/cokoies.txt -p 3003:3003 sunls24/vidlp
```

### Docker Compose & Caddy (Recommend)

**The same requires cookies.txt**

```yaml
version: '3.0'

services:
  vidlp:
    container_name: vidlp
    image: sunls24/vidlp:latest
    network_mode: host
    restart: unless-stopped
    volumes:
    - ./cookies.txt:/yt-dlp/cookies.txt:rw
```

#### Caddyfile

```text
dl.example.com {
    reverse_proxy 127.0.0.1:3003
}
```

## Support

- [x] Douyin / TikTok
- [x] Youtube
- [x] BiliBili
- [ ] [More](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md) (Should be supported, no filter format and tested)


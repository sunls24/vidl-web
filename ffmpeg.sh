#!/usr/bin/env ash

set -e

ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    target=linux64
elif [ "$ARCH" = "aarch64" ]; then
    target=linuxarm64
else
    echo "unsupported ""$ARCH"
    exit 1
fi

url=https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-${target}-gpl.tar.xz
wget -qO- $url | tar -xJ --strip-components=1 -C .
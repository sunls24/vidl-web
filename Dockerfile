FROM node:20-alpine as builder

RUN wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.shrc" SHELL="$(which sh)" sh -

WORKDIR /web
COPY ./web .
RUN source /root/.shrc && pnpm install && pnpm build

FROM golang:1.22-alpine AS builder2

WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=builder /web/dist ./web/dist
RUN CGO_ENABLED=0 go build -o vidlp && ./ffmpeg.sh

FROM python:3.11-slim

WORKDIR /yt-dlp

ENV GIN_MODE=release
RUN python3 -m pip install --no-deps -U --pre yt-dlp

COPY --from=builder2 /build/vidlp .
COPY --from=builder2 /build/bin/ffmpeg ffmpeg
EXPOSE 3003
ENTRYPOINT ["/yt-dlp/vidlp"]
# zynu-server

High-performance Go backend for the ZynU platform.
Handles HLS video streaming, upload pre-signing, and webhook processing.

![Go](https://img.shields.io/badge/Go-1.22-00ADD8)
![license](https://img.shields.io/badge/license-MIT-blue)
![build](https://img.shields.io/badge/build-passing-brightgreen)

## Features

- **HLS Streaming** — Serve `.m3u8` playlists and `.ts` segments with immutable cache headers
- **Upload Pre-signing** — Generate short-lived pre-signed upload URLs
- **Webhook Processing** — Handle `video.processed` and `video.failed` events
- **Middleware** — Bearer auth, per-IP rate limiting, request logging, CORS

## Setup

```bash
cp .env.example .env   # fill in your secrets
go build -o zynu-server ./cmd/server
./zynu-server
```

## Docker

```bash
docker build -t zynu-server .
docker run -p 8080:8080 --env-file .env zynu-server
```

## Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/stream/{id}/{quality}/{segment}` | HLS segment delivery |
| POST | `/upload/presign?video_id=...` | Get pre-signed upload URL |
| POST | `/webhook/video` | Video processing webhook |

## License

MIT © ZynU

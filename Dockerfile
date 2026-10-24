# ── Stage 1: build ────────────────────────────────────────────
FROM golang:1.27-alpine AS builder

WORKDIR /app

# cache deps first
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY *.go ./

# static binary — no cgo, no external libs needed
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o checker .

# ── Stage 2: minimal runtime ───────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/checker .
COPY sites.txt .

EXPOSE 8081

CMD ["./checker"]

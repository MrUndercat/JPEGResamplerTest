FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o jpeg_resampler ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/jpeg_resampler .

EXPOSE 8085

CMD ["./jpeg_resampler", "-path-orig", "/tmp/img_orig", "-path-res", "/tmp/img_res"]

# syntax=docker/dockerfile:1.4

##
## 1) Build Flutter Web (amd64 seulement)
##
FROM --platform=linux/amd64 neosu/flutter-web:3.13.9 AS flutter-web-builder

WORKDIR /ui
COPY ui/ /ui/

RUN flutter pub get
RUN flutter build web \
    --dart-define=BUILD_TYPE=prod \
    --tree-shake-icons \
    --pwa-strategy none \
    --web-renderer html \
    --release

##
## 2) Build Golang (multi-arch)
##
FROM --platform=$BUILDPLATFORM golang:1.21.3 AS go-builder

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH

WORKDIR /app

# On copie d'abord les go.mod / go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copie du reste du code
COPY . .

# On copie les fichiers web compilés par Flutter
COPY --from=flutter-web-builder /ui/build/web /app/web

# Compilation Go
RUN go build -a -ldflags="-w -s" -o /app/cmd/kubebadges/main ./cmd/kubebadges

##
## 3) Image finale (multi-arch)
##
FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /app
COPY --from=go-builder /app/cmd/kubebadges/main /app/main
RUN chmod +x /app/main

EXPOSE 8080 8090

CMD ["/app/main"]

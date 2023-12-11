# Build stage
FROM golang:1.19.0 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .

RUN go mod download

# Build
COPY . .

RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 GOARCH=amd64 go build -o solax-cloud-prometheus-exporter -ldflags "-X main.GitCommit=${GIT_COMMIT}"

RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 GOARCH=arm64 go build -o solax-cloud-prometheus-exporter-arm64 -ldflags "-X main.GitCommit=${GIT_COMMIT}"

# Final stage
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy the x86 binary
COPY --from=builder /build/solax-cloud-prometheus-exporter /app

# Copy the arm64 binary
COPY --from=builder /build/solax-cloud-prometheus-exporter-arm64 /app

EXPOSE 8886

ENTRYPOINT ["/app/solax-cloud-prometheus-exporter"]

CMD ["/app/solax-cloud-prometheus-exporter", "-listen", "0.0.0.0:8886"]


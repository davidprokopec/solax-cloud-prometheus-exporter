ARG ARCH=
FROM golang:1.19.0 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .

RUN go mod download

# Build
COPY . .

RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 GOARCH=${ARCH} go build -o solax-cloud-prometheus-exporter -ldflags "-X main.GitCommit=${GIT_COMMIT}"


# Final stage
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /build/solax-cloud-prometheus-exporter /app

EXPOSE 8886

ENTRYPOINT ["/app/solax-cloud-prometheus-exporter"]

CMD ["/app/solax-cloud-prometheus-exporter", "-listen", "0.0.0.0:8886"]


# solax-cloud-prometheus-exporer

Prometheus exporter for inverter data readouts from public solax cloud api. (currently only exports data i need, will maybe add support for selecting wanted data later)

Using Go 1.19 or newer

```shell
go install github.com/davidprokopec/solax-cloud-prometheus-exporter@latest
```

## Usage

Metrics prefix is `solaxcloud_`

### Run from shell

If you downloaded the source code and built it, you can run it from shell:

```shell
./solax-cloud-prometheus-exporter -address "<your_api_address>" -listen 0.0.0.0:8886
```

### Run from docker

The docker image is built for amd64 and arm64 platforms, since I myself use arm64 on my server.

```shell
docker run ghcr.io/davidprokopec/solax-cloud-prometheus-exporter:latest -address "<your_api_address>"
```

### Deployment

Configuration for scraping from Prometheus

```yml
metrics:
  configs:
    - name: default
      scrape_configs:
        - job_name: 'solaxcloud_exporter'
          scrape_interval: 120s
          static_configs:
            - targets: ['solaxcloud-exporter:8886']
```

Example docker-compose configuration

```yml
version: '3.3'
services:
  solaxcloud-exporter:
    container_name: monitoring-solaxcloud_exporter
    image: ghcr.io/davidprokopec/solax-cloud-prometheus-exporter
    command:
      - '-address=https://www.solaxcloud.com/proxyApp/proxy/api/getRealtimeInfo.do?tokenId=<your_token>&sn=<your_sn>'
    restart: unless-stopped
    networks:
      - monitoring
```

## Acknowledgement

Original project taken from loafoe: https://github.com/loafoe/prometheus-solaxrt-exporter (Would not be able to make this in Go without it!)

## License

License is MIT

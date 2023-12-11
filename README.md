# solax-cloud-prometheus-exporer

Prometheus exporter for Inverter data readouts from public solax cloud api. (currently only exports data i need, will maybe add support for selecting wanted data later)

## Usage


## Install

Using Go 1.19 or newer

```shell
go install github.com/davidprokopec/solax-cloud-prometheus-exporter@latest
```

## Usage

### Run from shell

If you downloaded the source code and built it, you can run it from shell:

```shell
solax-cloud-prometheus-exporter -address "<your_api_address>" -listen 0.0.0.0:8886
```

### Run from docker

```shell
docker run ghcr.io/davidprokopec/solax-cloud-prometheus-exporter:latest -address "<your_api_address>"
```

### Deployment

```prometheus.yml
metrics:
  configs:
    - name: default
      scrape_configs:
        - job_name: 'solaxcloud_exporter'
          scrape_interval: 2s
          static_configs:
            - targets: ['solaxcloud-exporter:8886']
```


## Acknowledgement

Original project taken from loafoe: https://github.com/loafoe/prometheus-solaxrt-exporter (Would not be able to make this in Go without it!)

## License

License is MIT

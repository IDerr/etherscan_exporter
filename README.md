# etherscan_exporter
A etherscan exporter to monitor your blockchain nodes via [Prometheus](https://prometheus.io).
Metrics are retrieved using the EtherScan REST API.

To run it:

    go build app.go -o bin/etherscan_exporter
    bin/etherscan_exporter [flags]

## Exported Metrics
| Metric | Description |
| ------ | ------- |
| etherscan_up | Was the last etherscan query successful |
| etherscan_last_block | Last block of chain |

## Flags
    ./bin/etherscan_exporter --help

| Flag | Description | Default |
| ---- | ----------- | ------- |
| log.level | Logging level | `info` |
| web.listen-address | Address to listen on for telemetry | `:9142` |
| web.telemetry-path | Path under which to expose metrics | `/metrics` |

## Env Variables

Use a .env file in the local folder, or /etc/sysconfig/etherscan_exporter
```
CHAIN="main" # Could be main,ropsten,rinkby,kovan
APIKEY="eth"
```
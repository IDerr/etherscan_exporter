package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nanmu42/etherscan-api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apikey        = ""
	chain         = "main"
	ns            = "etherscan"
	es            *etherscan.Client
	listenAddress = flag.String("web.listen-address", ":9142",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(ns, "", "up"),
		"Was the last query successful.",
		nil, nil,
	)
	lastBlock = prometheus.NewDesc(
		prometheus.BuildFQName(ns, "", "last_block"),
		"Last block number",
		[]string{"chain"}, nil,
	)
)

type Exporter struct {
	chain string
}

func NewExporter(chain string) *Exporter {
	return &Exporter{
		chain: chain,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- lastBlock
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	blockNumber, err := es.BlockNumber(time.Now().Unix(), "before")
	if err != nil {
		fmt.Println(err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1,
		)
		ch <- prometheus.MustNewConstMetric(
			lastBlock, prometheus.GaugeValue, float64(blockNumber), chain,
		)
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assume env variables are set.")
	}

	if c := os.Getenv("APIKEY"); c != "" {
		apikey = c
	} else {
		panic("APIKEY must be set")
	}
	if c := os.Getenv("CHAIN"); c != "" {
		chain = c
	}

	if chain == "main" {
		es = etherscan.New(etherscan.Mainnet, apikey)
	}
	if chain == "ropsten" {
		es = etherscan.New(etherscan.Ropsten, apikey)
	}
	if chain == "rinkby" {
		es = etherscan.New(etherscan.Rinkby, apikey)
	}
	if chain == "kovan" {
		es = etherscan.New(etherscan.Kovan, apikey)
	}
	flag.Parse()

	exporter := NewExporter(chain)
	prometheus.MustRegister(exporter)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Etherscan Exporter</title></head>
            <body>
            <h1>Etherscan Exporter</h1>
            <p><a href='` + *metricsPath + `'>Metrics</a></p>
            </body>
            </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

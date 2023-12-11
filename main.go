package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/davidprokopec/solax-cloud-prometheus-exporter/solax"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var listenAddr string
var ethDevice string
var apiAddr string
var debug bool

var (
	metricNamePrefix = "solaxrt_"
	registry         = prometheus.NewRegistry()
)

var (
	acPowerMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "ac_power",
		Help: "Current power generation (Wh)",
	}, []string{
		"inverter_sn",
	})

	yieldTodayMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "yield_today",
		Help: "The yield for today (KWh)",
	}, []string{
		"inverter_sn",
	})

	yieldTotalMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "yield_total",
		Help: "The total yield of the system (KWh)",
	}, []string{
		"inverter_sn",
	})

	feedInPowerMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "feed_in_power",
		Help: "The feed in power (W)",
	}, []string{
		"inverter_sn",
	})

	feedInEnergyMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "feed_in_energy",
		Help: "The feed in energy (Wh)",
	}, []string{
		"inverter_sn",
	})

	consumeEnergyMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "consume_energy",
		Help: "The consume energy (Wh)",
	}, []string{
		"inverter_sn",
	})

	upMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "up",
		Help: "The inverter power on status",
	}, []string{
		"sn",
	})
)

func init() {
	registry.MustRegister(acPowerMetric)
	registry.MustRegister(yieldTotalMetrics)
	registry.MustRegister(yieldTodayMetric)
	registry.MustRegister(feedInPowerMetrics)
	registry.MustRegister(feedInEnergyMetrics)
	registry.MustRegister(consumeEnergyMetrics)
	registry.MustRegister(upMetric)
}

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8886", "Listen address for HTTP metrics")
	flag.StringVar(&apiAddr, "address", "http://solaxcloud.com/proxyApp/proxy/api/getRealTimeInfo.do?tokenId=?&sn=?", "The URL address of the Solax API with tokenId and SN")
	flag.Parse()

	go func() {
		sleep := false
		for {
			if sleep {
				time.Sleep(time.Second * 2)
			}
			sleep = true
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

      fmt.Printf("checking address %s...\n", apiAddr)

			_, err := solax.UrlValid(apiAddr)

			if err != nil {
				fmt.Printf("address %s is not locally reachableeeee, skipping refresh...\n", apiAddr)
        fmt.Printf("error: %v\n", err)
				continue
			}

			fmt.Printf("calling Realtime API at %s...\n", apiAddr)

			resp, err := solax.GetRealtimeInfo(ctx,
				solax.WithURL(apiAddr),
				solax.WithDebug(debug))

			cancel()

			data := resp.Result

			if err != nil {

				fmt.Printf("error: %v\n", err)
				upMetric.WithLabelValues("").Set(0)

				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Printf("not sleeping\n")
					sleep = false
				}

				continue
			}

			if data.InverterStatus == "100" {
				upMetric.WithLabelValues("").Set(1)
			} else {
				upMetric.WithLabelValues("").Set(0)
			}

			acPowerMetric.WithLabelValues(data.SN).Set(data.ACPower)
			yieldTodayMetric.WithLabelValues(data.SN).Set(data.YieldToday)
			yieldTotalMetrics.WithLabelValues(data.SN).Set(data.YieldTotal)
			feedInPowerMetrics.WithLabelValues(data.SN).Set(data.FeedInPower)
			feedInEnergyMetrics.WithLabelValues(data.SN).Set(data.FeedInEnergy)
			consumeEnergyMetrics.WithLabelValues(data.SN).Set(data.ConsumeEnergy)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	_ = http.ListenAndServe(listenAddr, nil)
}

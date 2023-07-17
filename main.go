package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var queues = []string{"queue_a", "queue_b", "queue_c", "queue_d", "queue_e", "queue_f"}

func generateData() map[string]float64 {
	data := make(map[string]float64)

	for _, queue := range queues {
		// Chance to generate
		if rand.Intn(2)%2 != 0 {
			data[queue] = rand.Float64()
		}
	}

	return data
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "testMetric",
			Help: "A not very useful metric",
		},
		[]string{"queue"},
	)

	reg.MustRegister(vec)

	go func() {
		for true {
			data := generateData()
			secondVec := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "testMetric",
					Help: "A not very useful metric",
				},
				[]string{"queue"},
			)

			for queue, val := range data {
				secondVec.WithLabelValues(queue).Set(val)
			}

			start := time.Now()

			reg.Unregister(vec)
			reg.MustRegister(secondVec)

			duration := time.Since(start).Microseconds()

			vec = secondVec

			fmt.Printf("Data regenerated, swap took %dms\n", duration)
			time.Sleep(10 * time.Second)
		}
	}()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

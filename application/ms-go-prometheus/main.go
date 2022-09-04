package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	productsInInvetory = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "productsInInventory",
		Help: "Quantity of products available in the inventory",
		ConstLabels: map[string]string{
			"inventory": "products-1232sj-2324ehm",
		},
	})

	totalProductsInsertions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "totalProductsInsertions",
		Help: "Quantity INSERT product operations",
		ConstLabels: map[string]string{
			"inventory": "products-1232sj-2324ehm",
		},
	})

	httpTotalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "httpTotalRequests",
		Help: "Total http requests",
	}, []string{})

	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "httpRequestDuration",
		Help: "duration of http request",
	}, []string{"handler"})
)

func main() {
	register := prometheus.NewRegistry()
	register.MustRegister(productsInInvetory, totalProductsInsertions, httpDuration, httpTotalRequests)

	go randomQuantity()

	productsHandler := http.HandlerFunc(products)

	d := promhttp.InstrumentHandlerDuration(httpDuration.MustCurryWith(prometheus.Labels{"handler": "products"}),
		promhttp.InstrumentHandlerCounter(httpTotalRequests, productsHandler))

	http.Handle("/", d)

	http.Handle("/metrics", promhttp.HandlerFor(register, promhttp.HandlerOpts{}))

	log.Fatal(http.ListenAndServe(":8881", nil))
}

func products(rw http.ResponseWriter, r *http.Request) {
	randomT := rand.Intn(1000)
	time.Sleep(time.Millisecond * time.Duration(randomT))
	if r.Method == http.MethodPost {
		log.Printf("Inserting new product in inventory: %d", randomT)
		totalProductsInsertions.Inc()
		rw.Write([]byte("POST"))
	}
}

func randomQuantity() {
	var randomQ int
	for {
		time.Sleep(time.Second * 5)
		randomQ = rand.Intn(1000)
		log.Printf("Setting new random quantity in inventory: %d", randomQ)
		productsInInvetory.Set(float64(randomQ))
	}
}

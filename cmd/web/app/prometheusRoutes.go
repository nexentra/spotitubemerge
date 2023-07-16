package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/alice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// devices for prometheus guage
type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

type Metrics struct {
	devices prometheus.Gauge
	info   *prometheus.GaugeVec
	upgrades *prometheus.CounterVec
	duration *prometheus.HistogramVec
	loginDuration prometheus.Summary
}

// will need to move it to a dependency struct later on to be able to use it in the routes
var reg = prometheus.NewRegistry()
var metrics *Metrics = NewMetrics(reg)
func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "spotitubemerge",
			Name:      "connected_devices",
			Help:      "Number of connected devices",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "spotitubemerge",
			Name:      "info",
			Help:      "Info about the app",
		}, []string{"version"}),
		upgrades: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "spotitubemerge",
			Name:      "upgraded_devices",
			Help:      "Number of upgraded devices",
		}, []string{"type"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "spotitubemerge",
			Name:      "request_duration_seconds",
			Help:      "Time taken to process request",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status","method"}),
		loginDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "spotitubemerge",
			Name:      "login_request_duration_seconds",
			Help:      "Time taken to process login request",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}),
	}

	reg.MustRegister(m.devices, m.info, m.upgrades, m.duration, m.loginDuration)
	return m
}

var dvs []Device
var version string

func init() {
	fmt.Println("init")
	dvs = []Device{
		{1, "5F-33-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
	version = "0.0.1"
}



// PrometheusRoutes returns a handler for prometheus metrics
func (app *Application) PrometheusRoutes(mux *http.ServeMux) http.Handler {
	router := echo.New()
	// reg := prometheus.NewRegistry()
	// metrics = NewMetrics(reg)
	m := metrics

	m.devices.Set(float64(len(dvs)))
	m.info.WithLabelValues(version).Set(1)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	router.Use(middleware.Recover())

	// prometheus middleware for echo
	// router.Use(echoprometheus.NewMiddleware("spotitubemerge"))
	router.GET("/metrics", echo.WrapHandler(promHandler))

	standard := alice.New(app.logRequest, secureHeaders)
	return standard.Then(router)
}

func getDevices(c echo.Context) error {
	now := time.Now()
	b, err := json.Marshal(dvs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	sleep(200)

	metrics.duration.WithLabelValues("200", "GET").Observe(time.Since(now).Seconds())

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(b)
	return nil
}


func createDevice(c echo.Context) error {
    var dv Device
	m := metrics

    err := c.Bind(&dv)
    if err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }

    dvs = append(dvs, dv)

	m.devices.Set(float64(len(dvs)))

    return c.String(http.StatusCreated, "Device created!")
}

func upgradeDevice(c echo.Context) error {
	path := c.Param("id")

	id, err := strconv.Atoi(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var dv Device
	err = c.Bind(&dv)
    if err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }

	for i := range dvs {
		if dvs[i].ID == id {
			dvs[i].Firmware = dv.Firmware
		}
	}

	sleep(1000)

	// metrics.upgrades.With(prometheus.Labels{"type": "router"}).Inc()
	metrics.upgrades.WithLabelValues("router").Inc()

	return c.String(http.StatusCreated, "Upgrading device!")
}

func sleep(ms int){
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	n := rand.Intn(ms + now.Second())
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func login(c echo.Context) error {
	sleep(200)
	return c.String(http.StatusOK, "login")
}

func loginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		now := time.Now()
		metrics.loginDuration.Observe(time.Since(now).Seconds())
		return next(c)
	}
}
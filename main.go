package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	AwairURL = os.Getenv("AWAIR_URL")

	scoreGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_score",
		},
		[]string{},
	)

	dewPointGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_dew_point",
		},
		[]string{},
	)

	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_temp",
		},
		[]string{},
	)

	humidGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_humid",
		},
		[]string{},
	)

	absHumidGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_abs_humid",
		},
		[]string{},
	)

	co2Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2",
		},
		[]string{},
	)

	co2EstGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2_est",
		},
		[]string{},
	)

	co2EstBaselineGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2_est_baseline",
		},
		[]string{},
	)

	vocGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc",
		},
		[]string{},
	)

	vocBaselineGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_baseline",
		},
		[]string{},
	)

	vocH2RawGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_h2_raw",
		},
		[]string{},
	)

	vocEthanolRawGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_ethanol_raw",
		},
		[]string{},
	)

	pm25Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_pm25",
		},
		[]string{},
	)

	pm10EstGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_pm10_est",
		},
		[]string{},
	)
)

type AirData struct {
	Timestamp      time.Time `json:"timestamp"`
	Score          float64   `json:"score"`
	DewPoint       float64   `json:"dew_point"`
	Temp           float64   `json:"temp"`
	Humid          float64   `json:"humid"`
	AbsHumid       float64   `json:"abs_humid"`
	Co2            float64   `json:"co2"`
	Co2Est         float64   `json:"co2_est"`
	Co2EstBaseline float64   `json:"co2_est_baseline"`
	VOC            float64   `json:"voc"`
	VOCBaseline    float64   `json:"voc_baseline"`
	VOCH2Raw       float64   `json:"voc_h2_raw"`
	VOCEthanolRaw  float64   `json:"voc_ethanol_raw"`
	PM25           float64   `json:"pm25"`
	PM10Est        float64   `json:"pm10_est"`
}

func init() {
	prometheus.MustRegister(scoreGauge)
	prometheus.MustRegister(dewPointGauge)
	prometheus.MustRegister(tempGauge)
	prometheus.MustRegister(humidGauge)
	prometheus.MustRegister(absHumidGauge)
	prometheus.MustRegister(co2Gauge)
	prometheus.MustRegister(co2EstGauge)
	prometheus.MustRegister(co2EstBaselineGauge)
	prometheus.MustRegister(vocGauge)
	prometheus.MustRegister(vocBaselineGauge)
	prometheus.MustRegister(vocH2RawGauge)
	prometheus.MustRegister(vocEthanolRawGauge)
	prometheus.MustRegister(pm25Gauge)
	prometheus.MustRegister(pm10EstGauge)
}

func fetchAirData() {
	resp, err := http.Get(AwairURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		airData := AirData{}

		if err := json.Unmarshal(body, &airData); err != nil {
			return
		}

		scoreGauge.With(prometheus.Labels{}).Set(airData.Score)
		dewPointGauge.With(prometheus.Labels{}).Set(airData.DewPoint)
		tempGauge.With(prometheus.Labels{}).Set(airData.Temp)
		humidGauge.With(prometheus.Labels{}).Set(airData.Humid)
		absHumidGauge.With(prometheus.Labels{}).Set(airData.AbsHumid)
		co2Gauge.With(prometheus.Labels{}).Set(airData.Co2)
		co2EstGauge.With(prometheus.Labels{}).Set(airData.Co2Est)
		co2EstBaselineGauge.With(prometheus.Labels{}).Set(airData.Co2EstBaseline)
		vocGauge.With(prometheus.Labels{}).Set(airData.VOC)
		vocBaselineGauge.With(prometheus.Labels{}).Set(airData.VOCBaseline)
		vocH2RawGauge.With(prometheus.Labels{}).Set(airData.VOCH2Raw)
		vocEthanolRawGauge.With(prometheus.Labels{}).Set(airData.VOCEthanolRaw)
		pm25Gauge.With(prometheus.Labels{}).Set(airData.PM25)
		pm10EstGauge.With(prometheus.Labels{}).Set(airData.PM10Est)
	}
}

func pollAirData(d time.Duration) {
	timer := time.NewTicker(d)
	for {
		<-timer.C
		fetchAirData()
	}
}

func main() {
	if AwairURL == "" {
		panic("You must supply an Awair URL")
	}

	duration := os.Getenv("POLL_DURATION")
	if duration == "" {
		duration = "30s"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	go pollAirData(d)

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8181", nil); err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	AwairURL = os.Getenv("AWAIR_URL")

	labelNames = []string{"device_uuid", "ip"}

	scoreGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_score",
		},
		labelNames,
	)

	dewPointGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_dew_point",
		},
		labelNames,
	)

	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_temp",
		},
		labelNames,
	)

	humidGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_humid",
		},
		labelNames,
	)

	absHumidGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_abs_humid",
		},
		labelNames,
	)

	co2Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2",
		},
		labelNames,
	)

	co2EstGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2_est",
		},
		labelNames,
	)

	co2EstBaselineGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2_est_baseline",
		},
		labelNames,
	)

	vocGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc",
		},
		labelNames,
	)

	vocBaselineGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_baseline",
		},
		labelNames,
	)

	vocH2RawGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_h2_raw",
		},
		labelNames,
	)

	vocEthanolRawGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_ethanol_raw",
		},
		labelNames,
	)

	pm25Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_pm25",
		},
		labelNames,
	)

	pm10EstGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_pm10_est",
		},
		labelNames,
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

type Config struct {
	DeviceUUID      string `json:"device_uuid"`
	WiFiMAC         string `json:"wifi_mac"`
	SSID            string `json:"ssid"`
	IP              string `json:"ip"`
	Netmask         string `json:"netmask"`
	Gateway         string `json:"gateway"`
	FirmwareVersion string `json:"fw_version"`
	Timezone        string `json:"timezone"`
	Display         string `json:"display"`
}

func (c Config) ToPrometheusLabels() prometheus.Labels {
	return prometheus.Labels{
		"device_uuid": c.DeviceUUID,
		"ip":          c.IP,
	}
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

func fetchAirData(labels prometheus.Labels) {
	latestDataURL := fmt.Sprintf("%s/air-data/latest", AwairURL)
	resp, err := http.Get(latestDataURL)
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

		scoreGauge.With(labels).Set(airData.Score)
		dewPointGauge.With(labels).Set(airData.DewPoint)
		tempGauge.With(labels).Set(airData.Temp)
		humidGauge.With(labels).Set(airData.Humid)
		absHumidGauge.With(labels).Set(airData.AbsHumid)
		co2Gauge.With(labels).Set(airData.Co2)
		co2EstGauge.With(labels).Set(airData.Co2Est)
		co2EstBaselineGauge.With(labels).Set(airData.Co2EstBaseline)
		vocGauge.With(labels).Set(airData.VOC)
		vocBaselineGauge.With(labels).Set(airData.VOCBaseline)
		vocH2RawGauge.With(labels).Set(airData.VOCH2Raw)
		vocEthanolRawGauge.With(labels).Set(airData.VOCEthanolRaw)
		pm25Gauge.With(labels).Set(airData.PM25)
		pm10EstGauge.With(labels).Set(airData.PM10Est)
	}
}

func fetchAwairConfig() (Config, error) {
	configURL := fmt.Sprintf("%s/settings/config/data", AwairURL)
	resp, err := http.Get(configURL)
	if err != nil {
		return Config{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Config{}, errors.New("invalid status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := json.Unmarshal(body, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func pollAirData(d time.Duration, l prometheus.Labels) {
	timer := time.NewTicker(d)
	for {
		<-timer.C
		fetchAirData(l)
	}
}

func main() {
	if AwairURL == "" {
		panic("You must supply an Awair URL")
	}

	config, err := fetchAwairConfig()
	if err != nil {
		panic(err)
	}

	duration := os.Getenv("POLL_DURATION")
	if duration == "" {
		duration = "30s"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	go pollAirData(d, config.ToPrometheusLabels())

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8181", nil); err != nil {
		panic(err)
	}
}

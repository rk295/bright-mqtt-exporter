package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/certifi/gocertifi"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	electricityTopic = "electricitymeter"
	gasTopic         = "gasmeter"

	electricityMetricName = "electricity"
	gasMetricName         = "gas"

	mqttHostEnv     = "MQTT_HOST"
	mqttUserEnv     = "MQTT_USER"
	mqttPassEnv     = "MQTT_PASS"
	mqttTopicEnv    = "MQTT_TOPIC"
	exporterPortEnv = "PORT"

	mqttDefaultHost = "192.168.0.50:1883"
	mqttDefaultUser = "admin"

	exporterDefaultPort = "9999"
)

type Meters map[string]float64

type Data struct {
	Usage          Meters
	UnitRate       Meters
	StandingCharge Meters
}

type config struct {
	mqttHost     string
	mqttUser     string
	mqttPass     string
	mqttTopic    string
	exporterPort string
}

var (
	currentValues Data

	electricityUsageDetails = prometheus.NewDesc(
		prometheus.BuildFQName("uk_riviera", "monitoring", "electricity"),
		"electricity power usage readings from the smart meter in kWh",
		[]string{}, nil,
	)

	gasUsageDetails = prometheus.NewDesc(
		prometheus.BuildFQName("uk_riviera", "monitoring", "gas"),
		"gas usage readings from the smart meter in kWh",
		[]string{}, nil,
	)

	rateDetails = prometheus.NewDesc(
		prometheus.BuildFQName("uk_riviera", "monitoring", "price_per_unit"),
		"price per power (kWh) unit",
		[]string{"source"}, nil,
	)

	standingChartDetails = prometheus.NewDesc(
		prometheus.BuildFQName("uk_riviera", "monitoring", "standing_charge"),
		"price per power (kWh) unit",
		[]string{"source"}, nil,
	)
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.Debug("starting...")
	currentValues.Usage = make(map[string]float64)
	currentValues.UnitRate = make(map[string]float64)
	currentValues.StandingCharge = make(map[string]float64)
}

func main() {

	config, err := newConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	var qos byte

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.mqttHost)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetPassword(config.mqttPass)
	opts.SetUsername(config.mqttUser)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Debugf("subscribing to topic %s", config.mqttTopic)
	_ = client.Subscribe(config.mqttTopic, qos, currentValues.newMessage)

	certPool, err := gocertifi.CACerts()
	if err != nil {
		log.Fatalln("failed to initialize root CA pool:", err)
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	prometheus.MustRegister(currentValues)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("starting metrics server on port %s", config.exporterPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.exporterPort), nil); err != nil {
		log.Fatal(err)
	}
}

func (d Data) newMessage(c mqtt.Client, m mqtt.Message) {

	switch {
	case strings.HasSuffix(m.Topic(), electricityTopic):
		t := &BrightElectricitysMsg{}
		if err := json.Unmarshal(m.Payload(), &t); err != nil {
			log.Error(err)
			return
		}

		err := d.updateCurrent(t.Electricitymeter, electricityMetricName)
		if err != nil {
			log.Error(err)
		}

	case strings.HasSuffix(m.Topic(), gasTopic):
		t := &BrightGasMsg{}
		if err := json.Unmarshal(m.Payload(), &t); err != nil {
			log.Error(err)
			return
		}

		err := d.updateCurrent(t.Gasmeter, gasMetricName)
		if err != nil {
			log.Error(err)
		}
	default:
		return
	}

}

func (d Data) updateCurrent(m Meter, kind string) error {

	switch kind {
	case electricityMetricName:
		log.Debugf("mqtt: updating %s with %v", electricityMetricName, m.Power)
		d.Usage[kind] = float64(m.Power)
	case gasMetricName:
		log.Debugf("mqtt: updating %s with %v", gasMetricName, m.Energy.Import.Cummulative)
		d.Usage[kind] = float64(m.Energy.Import.Cummulative)
	default:
		return fmt.Errorf("unknown meter kind %s", kind)
	}

	d.UnitRate[kind] = float64(m.Price.Unitrate)             // Unit rate update
	d.StandingCharge[kind] = float64(m.Price.Standingcharge) // Standing charge update

	return nil
}

func (d Data) Describe(ch chan<- *prometheus.Desc) {
	ch <- electricityUsageDetails
}

func (d Data) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		electricityUsageDetails,
		prometheus.GaugeValue,
		d.Usage[electricityMetricName],
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		gasUsageDetails,
		prometheus.CounterValue,
		d.Usage[gasMetricName],
		[]string{}...,
	)

	for r, d := range d.UnitRate {
		ch <- prometheus.MustNewConstMetric(
			rateDetails,
			prometheus.GaugeValue,
			d,
			[]string{r}...,
		)
	}

	for r, d := range d.StandingCharge {
		ch <- prometheus.MustNewConstMetric(
			standingChartDetails,
			prometheus.GaugeValue,
			d,
			[]string{r}...,
		)
	}

}

func newConfig() (*config, error) {
	c := &config{}

	mqttHost := os.Getenv(mqttHostEnv)
	if mqttHost == "" {
		log.Debugf("%s not set, using default host of %s", mqttHostEnv, mqttDefaultHost)
		mqttHost = mqttDefaultHost
	}
	c.mqttHost = mqttHost

	mqttUser := os.Getenv(mqttUserEnv)
	if mqttUser == "" {
		log.Debugf("%s not set, using default user of %s", mqttUserEnv, mqttDefaultUser)
		mqttUser = mqttDefaultUser
	}
	c.mqttUser = mqttUser

	mqttPass := os.Getenv(mqttPassEnv)
	if mqttPass == "" {
		return c, fmt.Errorf("the %s variable must be set to the connection password", mqttPassEnv)
	}
	c.mqttPass = mqttPass

	mqttTopic := os.Getenv(mqttTopicEnv)
	if mqttTopic == "" {
		return c, fmt.Errorf("the %s variable must be set to the topic", mqttTopicEnv)
	}
	c.mqttTopic = mqttTopic

	exporterPort := os.Getenv(exporterPortEnv)
	if exporterPort == "" {
		log.Debugf("%s not set, using default port of %s", exporterPortEnv, exporterDefaultPort)
		exporterPort = exporterDefaultPort
	}
	c.exporterPort = exporterPort

	log.Debugf("mqtt config: host=%s user=%s topic=%s exporter-port=%s", mqttHost, mqttUser, mqttTopic, exporterPort)

	return c, nil

}

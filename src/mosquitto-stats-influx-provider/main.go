package main

import (
	"flag"
	"fmt"
	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
	"time"
)

// TODO: Is it really neccessary to declare this as a global variable ?
var mutex sync.Mutex
var mqttStats mqttStatistics
var config configuration

func main() {
	var configFile = flag.String("config", "", "Path to configuration file")
	var help = flag.Bool("help", false, "Show help text")
	var version = flag.Bool("version", false, "Show version")
	var verbose = flag.Bool("verbose", false, "Verbose output")
	var _fmt = new(log.TextFormatter)
	var wg sync.WaitGroup
	var err error

	_fmt.FullTimestamp = true
	_fmt.TimestampFormat = time.RFC3339
	log.SetFormatter(_fmt)

	flag.Usage = showUsage
	flag.Parse()

	if *help {
		showUsage()
		os.Exit(0)
	}

	if *version {
		showVersion()
		os.Exit(0)
	}

	if *configFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Missing configuration file\n\n")
		showUsage()
		os.Exit(1)
	}

	if *verbose {
		log.WithFields(log.Fields{
			"config_file": *configFile,
		}).Info("Parsing configuration file")
	}

	config, err = parseConfigFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	config.verbose = *verbose

	if config.MQTT.Password != "" {
		config.MQTT.mqttPassword = "<redacted>"
	}

	if config.verbose {
		log.WithFields(log.Fields{
			"service.listen":              config.Service.Listen,
			"service.influx_endpoint":     config.Service.InfluxEndpoint,
			"service.prometheus_endpoint": config.Service.PrometheusEndpoint,
		}).Info("Validating service configuration")

		log.WithFields(log.Fields{
			"mqtt.server":        config.MQTT.Server,
			"mqtt.client_id":     config.MQTT.ClientID,
			"mqtt.tls_ca":        config.MQTT.TLSCA,
			"mqtt.tls_cert":      config.MQTT.TLSCert,
			"mqtt.tls_key":       config.MQTT.TLSKey,
			"mqtt.insecure_ssl":  config.MQTT.InsecureSSL,
			"mqtt.qos":           config.MQTT.QoS,
			"mqtt.username":      config.MQTT.Username,
			"mqtt.password":      config.MQTT.mqttPassword,
			"mqtt.password_file": config.MQTT.PasswordFile,
		}).Info("Validating MQTT configuration")
	}
	err = validateConfiguration(config)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc(config.Service.InfluxEndpoint, influxHandler)
	router.HandleFunc(config.Service.PrometheusEndpoint, prometheusHandler)

	wg.Add(2)

	if config.verbose {
		log.WithFields(log.Fields{
			"mqtt.server":        config.MQTT.Server,
			"mqtt.client_id":     config.MQTT.ClientID,
			"mqtt.tls_ca":        config.MQTT.TLSCA,
			"mqtt.tls_cert":      config.MQTT.TLSCert,
			"mqtt.tls_key":       config.MQTT.TLSKey,
			"mqtt.insecure_ssl":  config.MQTT.InsecureSSL,
			"mqtt.qos":           config.MQTT.QoS,
			"mqtt.username":      config.MQTT.Username,
			"mqtt.password":      config.MQTT.mqttPassword,
			"mqtt.password_file": config.MQTT.PasswordFile,
		}).Info("Starting MQTT client to connect to broker")
	}

	go mqttConnect(config, &wg)

	if config.verbose {
		log.WithFields(log.Fields{
			"service.listen":              config.Service.Listen,
			"service.influx_endpoint":     config.Service.InfluxEndpoint,
			"service.prometheus_endpoint": config.Service.PrometheusEndpoint,
		}).Info("Starting web server to provide data at the configured endpoints")
	}
	go func() {
		defer wg.Done()
		http.ListenAndServe(config.Service.Listen, router)
	}()
	wg.Wait()

}

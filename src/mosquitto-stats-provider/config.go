package main

import (
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
	ini "gopkg.in/ini.v1"
	"net/url"
)

func parseConfigFile(f string) (configuration, error) {
	var cfg configuration

	c, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, f)

	if err != nil {
		return cfg, err
	}

	svcSect, err := c.GetSection("service")
	if err != nil {
		return cfg, err
	}
	err = svcSect.MapTo(&cfg.Service)
	if err != nil {
		return cfg, err
	}
	mqttSect, err := c.GetSection("mosquitto")
	if err != nil {
		return cfg, err
	}
	err = mqttSect.MapTo(&cfg.MQTT)
	if err != nil {
		return cfg, err
	}
	if cfg.MQTT.Password != "" && cfg.MQTT.PasswordFile != "" {
		return cfg, fmt.Errorf("Either password or password_file should be set but not both")
	}
	if cfg.MQTT.PasswordFile != "" {
		cfg.MQTT.Password, err = readPasswordFile(cfg.MQTT.PasswordFile)
		if err != nil {
			return cfg, err
		}
	}

	if cfg.MQTT.Username != "" {
		cfg.MQTT.authMethod = AuthMethodUserPass
	}
	if cfg.MQTT.TLSCert != "" {
		cfg.MQTT.authMethod = AuthMethodSSLClientCertificate
	}

	setDefaultConfigurationValues(&cfg)

	return cfg, nil
}

func setDefaultConfigurationValues(c *configuration) {
	if c.Service.Listen == "" {
		c.Service.Listen = "localhost:8383"
	}
	if c.Service.InfluxEndpoint == "" {
		c.Service.InfluxEndpoint = influxDefaultEndpoint
	}
	if c.Service.PrometheusEndpoint == "" {
		c.Service.PrometheusEndpoint = prometheusDefaultEndpoint
	}
	if c.Service.KeyValueEndpoint == "" {
		c.Service.KeyValueEndpoint = kvDefaultEndpoint
	}

	if c.MQTT.ClientID == "" {
		_uuid, err := uuid.NewV4()
		if err != nil {
			log.Panic(fmt.Sprintf("Unable to generate UUID for go routine: %s", err.Error()))
		}
		c.MQTT.ClientID = _uuid.String()
	}
}

func validateConfiguration(c configuration) error {
	if c.MQTT.Server == "" {
		return fmt.Errorf("No MQTT server configured")
	}

	u, err := url.Parse(c.MQTT.Server)
	if err != nil {
		return err
	}
	if u.Scheme != "tcp" && u.Scheme != "ssl" && u.Scheme != "ws" {
		return fmt.Errorf("Can't parse MQTT server URL")
	}

	if (c.MQTT.Username == "" && c.MQTT.Password != "") || (c.MQTT.Username != "" && c.MQTT.Password == "") {
		return fmt.Errorf("For user/password authentication both username and password must be set")
	}
	if (c.MQTT.TLSCert != "" && c.MQTT.TLSKey == "") || (c.MQTT.TLSCert == "" && c.MQTT.TLSKey != "") {
		return fmt.Errorf("For authentication using SSL client certificates both public and private key must be set")
	}
	if c.MQTT.Username == "" && c.MQTT.Password == "" && c.MQTT.TLSCert == "" && c.MQTT.TLSKey == "" {
		return fmt.Errorf("No authentication methode configured")
	}
	if c.MQTT.Username != "" && c.MQTT.Password != "" && c.MQTT.TLSCert != "" && c.MQTT.TLSKey != "" {
		return fmt.Errorf("Multiple authentication methods configured. Please configure either user/password or SSL client certificate authentication, but not both")
	}
	if c.MQTT.QoS < 0 || c.MQTT.QoS > 2 {
		return fmt.Errorf("Invalid QoS value. QoS must be 0, 1 or 2")
	}

	return nil
}

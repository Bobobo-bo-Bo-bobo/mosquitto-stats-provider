package main

type mqttConfiguration struct {
	Server       string `ini:"server"`
	ClientID     string `ini:"client_id"`
	TLSCA        string `ini:"tls_ca"`
	TLSCert      string `ini:"tls_cert"`
	TLSKey       string `ini:"tls_key"`
	InsecureSSL  bool   `ini:"insecure_ssl"`
	QoS          uint   `ini:"qos"`
	Username     string `ini:"username"`
	Password     string `ini:"password"`
	PasswordFile string `ini:"password_file"`
	authMethod   uint
	mqttPassword string
}

type serviceConfiguration struct {
	Listen             string `ini:"listen"`
	InfluxEndpoint     string `ini:"influx_endpoint"`
	PrometheusEndpoint string `ini:"prometheus_endpoint"`
	KeyValueEndpoint   string `ini:"key_value_endpoint"`
}

type configuration struct {
	Service serviceConfiguration
	MQTT    mqttConfiguration
	verbose bool
}

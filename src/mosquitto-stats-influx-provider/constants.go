package main

const version = "1.0.0-20191030"

const (
	// QoSAtMostOnce - QoS 0
	QoSAtMostOnce uint = iota
	// QoSAtLeastOnce - QoS 1
	QoSAtLeastOnce
	// QoSExactlyOnce - QoS 2
	QoSExactlyOnce
)

const (
	// AuthMethodUserPass - authenticate using username/password
	AuthMethodUserPass uint = iota
	// AuthMethodSSLClientCertificate - authenticate using SSL client certificates
	AuthMethodSSLClientCertificate
)

const (
	_ uint = iota
	_
	_
	// MQTTv3_1 - use MQTT 3.1 protocol
	MQTTv3_1
	// MQTTv3_1_1 - use MQTT 3.1.1 protocol
	MQTTv3_1_1
)

const mosquittoStatisticsTopic = "$SYS/#"

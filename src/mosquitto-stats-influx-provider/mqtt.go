package main

import (
	"crypto/tls"
	"crypto/x509"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"sync"
)

func mqttConnect(cfg configuration, wg *sync.WaitGroup) {
	var tlsCfg = new(tls.Config)

	defer wg.Done()

	mqttOptions := mqtt.NewClientOptions()
	mqttOptions.AddBroker(cfg.MQTT.Server)

	if cfg.MQTT.authMethod == AuthMethodUserPass {
		mqttOptions.SetUsername(cfg.MQTT.Username)
		mqttOptions.SetPassword(cfg.MQTT.Password)
	}

	if cfg.MQTT.authMethod == AuthMethodSSLClientCertificate {
		cert, err := tls.LoadX509KeyPair(cfg.MQTT.TLSCert, cfg.MQTT.TLSKey)
		if err != nil {
			log.Fatal(err)
		}

		tlsCfg.Certificates = make([]tls.Certificate, 1)
		tlsCfg.Certificates[0] = cert

		mqttOptions.SetTLSConfig(tlsCfg)
	}

	if cfg.MQTT.TLSCA != "" {
		tlsCfg.RootCAs = x509.NewCertPool()
		cacert, err := ioutil.ReadFile(cfg.MQTT.TLSCA)
		if err != nil {
			log.Fatal(err)
		}

		tlsok := tlsCfg.RootCAs.AppendCertsFromPEM(cacert)
		if !tlsok {
			log.Fatal("Can't add CA certificate to x509.CertPool")
		}
	}

	mqttOptions.SetClientID(cfg.MQTT.ClientID)

	// XXX: this could be read from the configuration file
	mqttOptions.SetAutoReconnect(true)
	mqttOptions.SetConnectRetry(true)
	mqttOptions.SetProtocolVersion(MQTTv3_1_1)

	mqttClient := mqtt.NewClient(mqttOptions)

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
		}).Info("Connecting to MQTT broker")
	}

	mqttToken := mqttClient.Connect()
	mqttToken.Wait()
	if mqttToken.Error() != nil {
		log.Fatal(mqttToken.Error())
	}

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
			"mqtt.topic":         "$SYS/#",
		}).Info("Subscribing to statistics topic")
	}
	mqttClient.Subscribe(mosquittoStatisticsTopic, byte(cfg.MQTT.QoS), mqttMessageHandler)
}

func mqttMessageHandler(c mqtt.Client, msg mqtt.Message) {
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
			"message.topic":      msg.Topic(),
			"message.payload":    string(msg.Payload()),
		}).Info("Received MQTT message on subscribed topic")
	}
	mutex.Lock()
	{
		processMessages(msg)
	}
	mutex.Unlock()
	return
}

func processMessages(msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	if topic == "$SYS/broker/clients/total" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.totalClients = u
		}
	} else if topic == "$SYS/broker/clients/maximum" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.maximumClients = u
		}
	} else if topic == "$SYS/broker/clients/inactive" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.incativeClients = u
		}
	} else if topic == "$SYS/broker/clients/disconnected" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.disconnectedClients = u
		}
	} else if topic == "$SYS/broker/clients/active" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.activeClients = u
		}
	} else if topic == "$SYS/broker/clients/connected" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.connectedClients = u
		}
	} else if topic == "$SYS/broker/clients/expired" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.expiredClients = u
		}
	} else if topic == "$SYS/broker/load/messages/received/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesReceived1min = f
		}
	} else if topic == "$SYS/broker/load/messages/received/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesReceived5min = f
		}
	} else if topic == "$SYS/broker/load/messages/received/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesReceived15min = f
		}
	} else if topic == "$SYS/broker/load/messages/sent/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesSent1min = f
		}
	} else if topic == "$SYS/broker/load/messages/sent/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesSent5min = f
		}
	} else if topic == "$SYS/broker/load/messages/sent/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesSent15min = f
		}
	} else if topic == "$SYS/broker/load/publish/dropped/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishDropped1min = f
		}
	} else if topic == "$SYS/broker/load/publish/dropped/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishDropped5min = f
		}
	} else if topic == "$SYS/broker/load/publish/dropped/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishDropped15min = f
		}
	} else if topic == "$SYS/broker/load/publish/received/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishReceived1min = f
		}
	} else if topic == "$SYS/broker/load/publish/received/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishReceived5min = f
		}
	} else if topic == "$SYS/broker/load/publish/received/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishReceived15min = f
		}
	} else if topic == "$SYS/broker/load/publish/sent/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishSent1min = f
		}
	} else if topic == "$SYS/broker/load/publish/sent/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishSent5min = f
		}
	} else if topic == "$SYS/broker/load/publish/sent/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishSent15min = f
		}
	} else if topic == "$SYS/broker/load/bytes/received/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesReceived1min = f
		}
	} else if topic == "$SYS/broker/load/bytes/received/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesReceived5min = f
		}
	} else if topic == "$SYS/broker/load/bytes/received/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesReceived15min = f
		}
	} else if topic == "$SYS/broker/load/bytes/sent/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesSent1min = f
		}
	} else if topic == "$SYS/broker/load/bytes/sent/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesSent5min = f
		}
	} else if topic == "$SYS/broker/load/bytes/sent/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesSent15min = f
		}
	} else if topic == "$SYS/broker/load/sockets/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.sockets1min = f
		}
	} else if topic == "$SYS/broker/load/sockets/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.sockets5min = f
		}
	} else if topic == "$SYS/broker/load/sockets/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.sockets15min = f
		}
	} else if topic == "$SYS/broker/load/connections/1min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.connections1min = f
		}
	} else if topic == "$SYS/broker/load/connections/5min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.connections5min = f
		}
	} else if topic == "$SYS/broker/load/connections/15min" {
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.connections15min = f
		}
	} else if topic == "$SYS/broker/messages/stored" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesStored = u
		}
	} else if topic == "$SYS/broker/messages/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesReceived += u
		}
	} else if topic == "$SYS/broker/messages/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesSent += u
		}
	} else if topic == "$SYS/broker/store/messages/count" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesStoredCount += u
		}
	} else if topic == "$SYS/broker/store/messages/bytes" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.messagesStoredBytes += u
		}
	} else if topic == "$SYS/broker/subscriptions/count" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.subscriptionsCount += u
		}
	} else if topic == "$SYS/broker/retained messages/count" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.retainedMessagesCount += u
		}
	} else if topic == "$SYS/broker/heap/current" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.heapCurrent = u
		}
	} else if topic == "$SYS/broker/heap/maximum" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.heapMaximum = u
		}
	} else if topic == "$SYS/broker/publish/messages/dropped" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishMessagesDropped += u
		}
	} else if topic == "$SYS/broker/publish/messages/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishMessagesReceived += u
		}
	} else if topic == "$SYS/broker/publish/messages/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishMessagesSent += u
		}
	} else if topic == "$SYS/broker/publish/bytes/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishBytesReceived += u
		}
	} else if topic == "$SYS/broker/publish/bytes/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.publishBytesSent += u
		}
	} else if topic == "$SYS/broker/bytes/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesReceived += u
		}
	} else if topic == "$SYS/broker/bytes/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			logStringConversionError(topic, payload, err)
		} else {
			mqttStats.bytesSent += u
		}
	} else if topic == "$$SYS/broker/version" {
		// skip
	} else if topic == "$SYS/broker/uptime" {
		// skip
	} else {
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
			"message.topic":      topic,
			"message.payload":    payload,
		}).Warn("Unhandled MQTT message on statistics topic")
		return
	}
}

func logStringConversionError(topic string, payload string, err error) {
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
		"message.topic":      topic,
		"message.payload":    payload,
		"error":              err,
	}).Error("Can't convert payload to a number")
}

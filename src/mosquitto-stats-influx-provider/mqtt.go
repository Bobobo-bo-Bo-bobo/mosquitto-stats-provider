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

	mqttToken := mqttClient.Connect()
	mqttToken.Wait()
	if mqttToken.Error() != nil {
		log.Fatal(mqttToken.Error())
	}

	mqttClient.Subscribe(mosquittoStatisticsTopic, byte(cfg.MQTT.QoS), mqttMessageHandler)
}

func mqttMessageHandler(c mqtt.Client, msg mqtt.Message) {
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
			return
		} else {
			mqttStats.totalClients = u
		}
		return
	} else if topic == "$SYS/broker/clients/maximum" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.maximumClients = u
		}

		return
	} else if topic == "$SYS/broker/clients/inactive" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.incativeClients = u
		}

		return
	} else if topic == "$SYS/broker/clients/disconnected" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.disconnectedClients = u
		}

		return
	} else if topic == "$SYS/broker/clients/active" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.activeClients = u
		}

		return
	} else if topic == "$SYS/broker/clients/connected" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.connectedClients = u
		}

		return
	} else if topic == "$SYS/broker/clients/expired" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.expiredClients = u
		}

		return
	} else if topic == "$SYS/broker/load/messages/received/1min" {
		return
	} else if topic == "$SYS/broker/load/messages/received/5min" {
		return
	} else if topic == "$SYS/broker/load/messages/received/15min" {
		return
	} else if topic == "$SYS/broker/load/messages/sent/1min" {
		return
	} else if topic == "$SYS/broker/load/messages/sent/5min" {
		return
	} else if topic == "$SYS/broker/load/messages/sent/15min" {
		return
	} else if topic == "$SYS/broker/load/publish/dropped/1min" {
		return
	} else if topic == "$SYS/broker/load/publish/dropped/5min" {
		return
	} else if topic == "$SYS/broker/load/publish/dropped/15min" {
		return
	} else if topic == "$SYS/broker/load/publish/received/1min" {
		return
	} else if topic == "$SYS/broker/load/publish/received/15min" {
		return
	} else if topic == "$SYS/broker/load/publish/sent/1min" {
		return
	} else if topic == "$SYS/broker/load/publish/sent/5min" {
		return
	} else if topic == "$SYS/broker/load/publish/sent/15min" {
		return
	} else if topic == "$SYS/broker/load/bytes/received/1min" {
		return
	} else if topic == "$SYS/broker/load/bytes/received/5min" {
		return
	} else if topic == "$SYS/broker/load/bytes/received/15min" {
		return
	} else if topic == "$SYS/broker/load/bytes/sent/1min" {
		return
	} else if topic == "$SYS/broker/load/bytes/sent/5min" {
		return
	} else if topic == "$SYS/broker/load/bytes/sent/15min" {
		return
	} else if topic == "$SYS/broker/load/sockets/1min" {
		return
	} else if topic == "$SYS/broker/load/sockets/5min" {
		return
	} else if topic == "$SYS/broker/load/sockets/15min" {
		return
	} else if topic == "$SYS/broker/load/connections/1min" {
		return
	} else if topic == "$SYS/broker/load/connections/5min" {
		return
	} else if topic == "$SYS/broker/load/connections/15min" {
		return
	} else if topic == "$SYS/broker/messages/stored" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.messagesStored = u
		}

		return
	} else if topic == "$SYS/broker/messages/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.messagesReceived += u
		}

		return
	} else if topic == "$SYS/broker/messages/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.messagesSent += u
		}

		return
	} else if topic == "$SYS/broker/store/messages/count" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.messagesStoredCount += u
		}

		return
	} else if topic == "$SYS/broker/store/messages/bytes" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.messagesStoredBytes += u
		}

		return
	} else if topic == "$SYS/broker/subscriptions/count" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.subscriptionsCount += u
		}

		return
	} else if topic == "$SYS/broker/retained messages/count" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.retainedMessagesCount += u
		}

		return
	} else if topic == "$SYS/broker/heap/current" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.heapCurrent = u
		}

		return
	} else if topic == "$SYS/broker/heap/maximum" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.heapMaximum = u
		}

		return
	} else if topic == "$SYS/broker/publish/messages/dropped" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.publishMessagesDropped += u
		}

		return
	} else if topic == "$SYS/broker/publish/messages/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.publishMessagesReceived += u
		}

		return
	} else if topic == "$SYS/broker/publish/messages/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.publishMessagesSent += u
		}

		return
	} else if topic == "$SYS/broker/publish/bytes/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.publishBytesReceived += u
		}

		return
	} else if topic == "$SYS/broker/publish/bytes/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.publishBytesSent += u
		}

		return
	} else if topic == "$SYS/broker/bytes/received" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.bytesReceived += u
		}
		return
	} else if topic == "$SYS/broker/bytes/sent" {
		u, err := strconv.ParseUint(payload, 10, 64)
		if err != nil {
			return
		} else {
			mqttStats.bytesSent += u
		}
		return
	} else {
	}
	return
}

package main

import (
	"fmt"
	"net/http"
	"time"
)

func influxHandler(writer http.ResponseWriter, request *http.Request) {
	var str []byte

	now := time.Now().Unix()
	mutex.Lock()
	{
		str = append(str, fmt.Sprintf("mosquitto_statistics,broker=%s total_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s maximum_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s incative_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s disconnected_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s active_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s connected_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s expired_clients=%d %d\n"+
			"mosquitto_statistics,broker=%s messages_received_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s messages_received_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s messages_received_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s messages_sent_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s messages_sent_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s messages_sent_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_dropped_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_dropped_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_dropped_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_received_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_received_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_received_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_sent_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_sent_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s publish_sent_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s bytes_received_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s bytes_received_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s bytes_received_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s bytes_sent_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s bytes_sent_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s bytes_sent_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s sockets_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s sockets_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s sockets_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s connections_1min=%f %d\n"+
			"mosquitto_statistics,broker=%s connections_5min=%f %d\n"+
			"mosquitto_statistics,broker=%s connections_15min=%f %d\n"+
			"mosquitto_statistics,broker=%s messages_stored=%d %d\n"+
			"mosquitto_statistics,broker=%s messages_received=%d %d\n"+
			"mosquitto_statistics,broker=%s messages_sent=%d %d\n"+
			"mosquitto_statistics,broker=%s messages_stored_count=%d %d\n"+
			"mosquitto_statistics,broker=%s messages_stored_bytes=%d %d\n"+
			"mosquitto_statistics,broker=%s subscriptions_count=%d %d\n"+
			"mosquitto_statistics,broker=%s retained_messages_count=%d %d\n"+
			"mosquitto_statistics,broker=%s heap_current=%d %d\n"+
			"mosquitto_statistics,broker=%s heap_maximum=%d %d\n"+
			"mosquitto_statistics,broker=%s publish_messages_dropped=%d %d\n"+
			"mosquitto_statistics,broker=%s publish_messages_received=%d %d\n"+
			"mosquitto_statistics,broker=%s publish_messages_sent=%d %d\n"+
			"mosquitto_statistics,broker=%s publish_bytes_received=%d %d\n"+
			"mosquitto_statistics,broker=%s publish_bytes_sent=%d %d\n"+
			"mosquitto_statistics,broker=%s bytes_received=%d %d\n"+
			"mosquitto_statistics,broker=%s bytes_sent=%d %d\n", config.MQTT.Server, mqttStats.totalClients, now,
			config.MQTT.Server, mqttStats.maximumClients, now,
			config.MQTT.Server, mqttStats.incativeClients, now,
			config.MQTT.Server, mqttStats.disconnectedClients, now,
			config.MQTT.Server, mqttStats.activeClients, now,
			config.MQTT.Server, mqttStats.connectedClients, now,
			config.MQTT.Server, mqttStats.expiredClients, now,
			config.MQTT.Server, mqttStats.messagesReceived1min, now,
			config.MQTT.Server, mqttStats.messagesReceived5min, now,
			config.MQTT.Server, mqttStats.messagesReceived15min, now,
			config.MQTT.Server, mqttStats.messagesSent1min, now,
			config.MQTT.Server, mqttStats.messagesSent5min, now,
			config.MQTT.Server, mqttStats.messagesSent15min, now,
			config.MQTT.Server, mqttStats.publishDropped1min, now,
			config.MQTT.Server, mqttStats.publishDropped5min, now,
			config.MQTT.Server, mqttStats.publishDropped15min, now,
			config.MQTT.Server, mqttStats.publishReceived1min, now,
			config.MQTT.Server, mqttStats.publishReceived5min, now,
			config.MQTT.Server, mqttStats.publishReceived15min, now,
			config.MQTT.Server, mqttStats.publishSent1min, now,
			config.MQTT.Server, mqttStats.publishSent5min, now,
			config.MQTT.Server, mqttStats.publishSent15min, now,
			config.MQTT.Server, mqttStats.bytesReceived1min, now,
			config.MQTT.Server, mqttStats.bytesReceived5min, now,
			config.MQTT.Server, mqttStats.bytesReceived15min, now,
			config.MQTT.Server, mqttStats.bytesSent1min, now,
			config.MQTT.Server, mqttStats.bytesSent5min, now,
			config.MQTT.Server, mqttStats.bytesSent15min, now,
			config.MQTT.Server, mqttStats.sockets1min, now,
			config.MQTT.Server, mqttStats.sockets5min, now,
			config.MQTT.Server, mqttStats.sockets15min, now,
			config.MQTT.Server, mqttStats.connections1min, now,
			config.MQTT.Server, mqttStats.connections5min, now,
			config.MQTT.Server, mqttStats.connections15min, now,
			config.MQTT.Server, mqttStats.messagesStored, now,
			config.MQTT.Server, mqttStats.messagesReceived, now,
			config.MQTT.Server, mqttStats.messagesSent, now,
			config.MQTT.Server, mqttStats.messagesStoredCount, now,
			config.MQTT.Server, mqttStats.messagesStoredBytes, now,
			config.MQTT.Server, mqttStats.subscriptionsCount, now,
			config.MQTT.Server, mqttStats.retainedMessagesCount, now,
			config.MQTT.Server, mqttStats.heapCurrent, now,
			config.MQTT.Server, mqttStats.heapMaximum, now,
			config.MQTT.Server, mqttStats.publishMessagesDropped, now,
			config.MQTT.Server, mqttStats.publishMessagesReceived, now,
			config.MQTT.Server, mqttStats.publishMessagesSent, now,
			config.MQTT.Server, mqttStats.publishBytesReceived, now,
			config.MQTT.Server, mqttStats.publishBytesSent, now,
			config.MQTT.Server, mqttStats.bytesReceived, now,
			config.MQTT.Server, mqttStats.bytesSent, now)...)
	}
	mutex.Unlock()
	writer.Write(str)
}

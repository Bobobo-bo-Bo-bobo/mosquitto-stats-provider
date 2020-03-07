package main

import (
	"fmt"
	"net/http"
)

func keyValueHandler(writer http.ResponseWriter, request *http.Request) {
	var str []byte

	mutex.Lock()
	{
		str = append(str, fmt.Sprintf("total_clients=%d\n"+
			"maximum_clients=%d\n"+
			"incative_clients=%d\n"+
			"disconnected_clients=%d\n"+
			"active_clients=%d\n"+
			"connected_clients=%d\n"+
			"expired_clients=%d\n"+
			"messages_received_1min=%f\n"+
			"messages_received_5min=%f\n"+
			"messages_received_15min=%f\n"+
			"messages_sent_1min=%f\n"+
			"messages_sent_5min=%f\n"+
			"messages_sent_15min=%f\n"+
			"publish_dropped_1min=%f\n"+
			"publish_dropped_5min=%f\n"+
			"publish_dropped_15min=%f\n"+
			"publish_received_1min=%f\n"+
			"publish_received_5min=%f\n"+
			"publish_received_15min=%f\n"+
			"publish_sent_1min=%f\n"+
			"publish_sent_5min=%f\n"+
			"publish_sent_15min=%f\n"+
			"bytes_received_1min=%f\n"+
			"bytes_received_5min=%f\n"+
			"bytes_received_15min=%f\n"+
			"bytes_sent_1min=%f\n"+
			"bytes_sent_5min=%f\n"+
			"bytes_sent_15min=%f\n"+
			"sockets_1min=%f\n"+
			"sockets_5min=%f\n"+
			"sockets_15min=%f\n"+
			"connections_1min=%f\n"+
			"connections_5min=%f\n"+
			"connections_15min=%f\n"+
			"messages_stored=%d\n"+
			"messages_received=%d\n"+
			"messages_sent=%d\n"+
			"messages_stored_count=%d\n"+
			"messages_stored_bytes=%d\n"+
			"subscriptions_count=%d\n"+
			"retained_messages_count=%d\n"+
			"heap_current=%d\n"+
			"heap_maximum=%d\n"+
			"publish_messages_dropped=%d\n"+
			"publish_messages_received=%d\n"+
			"publish_messages_sent=%d\n"+
			"publish_bytes_received=%d\n"+
			"publish_bytes_sent=%d\n"+
			"bytes_received=%d\n"+
			"bytes_sent=%d\n", mqttStats.totalClients,
			mqttStats.maximumClients,
			mqttStats.incativeClients,
			mqttStats.disconnectedClients,
			mqttStats.activeClients,
			mqttStats.connectedClients,
			mqttStats.expiredClients,
			mqttStats.messagesReceived1min,
			mqttStats.messagesReceived5min,
			mqttStats.messagesReceived15min,
			mqttStats.messagesSent1min,
			mqttStats.messagesSent5min,
			mqttStats.messagesSent15min,
			mqttStats.publishDropped1min,
			mqttStats.publishDropped5min,
			mqttStats.publishDropped15min,
			mqttStats.publishReceived1min,
			mqttStats.publishReceived5min,
			mqttStats.publishReceived15min,
			mqttStats.publishSent1min,
			mqttStats.publishSent5min,
			mqttStats.publishSent15min,
			mqttStats.bytesReceived1min,
			mqttStats.bytesReceived5min,
			mqttStats.bytesReceived15min,
			mqttStats.bytesSent1min,
			mqttStats.bytesSent5min,
			mqttStats.bytesSent15min,
			mqttStats.sockets1min,
			mqttStats.sockets5min,
			mqttStats.sockets15min,
			mqttStats.connections1min,
			mqttStats.connections5min,
			mqttStats.connections15min,
			mqttStats.messagesStored,
			mqttStats.messagesReceived,
			mqttStats.messagesSent,
			mqttStats.messagesStoredCount,
			mqttStats.messagesStoredBytes,
			mqttStats.subscriptionsCount,
			mqttStats.retainedMessagesCount,
			mqttStats.heapCurrent,
			mqttStats.heapMaximum,
			mqttStats.publishMessagesDropped,
			mqttStats.publishMessagesReceived,
			mqttStats.publishMessagesSent,
			mqttStats.publishBytesReceived,
			mqttStats.publishBytesSent,
			mqttStats.bytesReceived,
			mqttStats.bytesSent)...)
	}
	mutex.Unlock()
	writer.Write(str)
	str = nil
}

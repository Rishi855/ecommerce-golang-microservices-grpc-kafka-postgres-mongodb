package events

const(
	NOTIFICATION_TOPIC = "notification-events"
	ORDER_LOGS_TOPIC = "logs.order-service"
)
// docker exec -it kafka bash
// kafka-topics --create --topic notification-events --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
// kafka-topics --create --topic logs.order-service --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
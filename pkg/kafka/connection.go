package kafka

import (
	"context"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func InitKafka() (*kafka.Conn, error) {
	topic := os.Getenv("KAFKA_ROOT_TOPIC")
	conn, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_URL"), topic, 0)

	if err != nil {
		return nil, err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn, nil
}

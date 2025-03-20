package kafka

import (
	"auth-service/pkg/logger"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type LogProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewLogProducer(brokers []string, topic string) *LogProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		RequiredAcks: kafka.RequireOne,
	}

	return &LogProducer{
		writer: writer,
		topic:  topic,
	}
}

func (p *LogProducer) Close() error {
	return p.writer.Close()
}

func (p *LogProducer) SendLog(l logger.Log) error {
	logJSON, err := json.Marshal(l)
	if err != nil {
		log.Printf("Failed to convert log to JSON: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte(l.Category),
		Value: logJSON,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{
				Key:   "level",
				Value: []byte(l.Level),
			},
			{
				Key:   "category",
				Value: []byte(l.Category),
			},
		},
	}

	err = p.writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}

func (p *LogProducer) SendLogToTopic(l logger.Log, topic string) error {
	tempWriter := &kafka.Writer{
		Addr:         p.writer.Addr,
		Topic:        topic,
		Balancer:     p.writer.Balancer,
		BatchTimeout: p.writer.BatchTimeout,
	}
	defer tempWriter.Close()

	logJSON, err := json.Marshal(l)
	if err != nil {
		log.Printf("Failed to convert log to JSON: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte(l.Category),
		Value: logJSON,
		Time:  time.Now(),
	}

	err = tempWriter.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send message to topic %s: %v", topic, err)
		return err
	}

	return nil
}

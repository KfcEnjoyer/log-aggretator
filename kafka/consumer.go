package kafka

import (
	"auth-service/internal/database"
	"auth-service/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type LogConsumer struct {
	reader      *kafka.Reader
	logRepo     *database.LogRepoService
	sampleRates map[string]float64
}

func NewLogConsumer(brokers []string, topic, groupID string, logRepo *database.LogRepoService) *LogConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset,
		MaxWait:     time.Second,
	})

	sampleRates := map[string]float64{
		logger.ERROR: 1.0,
		logger.WARN:  1.0,
		logger.INFO:  0.5,
		logger.DEBUG: 0.1,
	}

	return &LogConsumer{
		reader:      reader,
		logRepo:     logRepo,
		sampleRates: sampleRates,
	}
}

func (c *LogConsumer) Start(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- c.consume(ctx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		log.Println("Shutting down log consumer...")
		if err := c.reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
		wg.Wait()
		return ctx.Err()
	}
}

func (c *LogConsumer) consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var logEntry logger.Log
		if err := json.Unmarshal(m.Value, &logEntry); err != nil {
			log.Printf("Error unmarshaling log entry: %v", err)
			continue
		}

		if c.shouldStoreLog(logEntry) {
			if err := c.logRepo.CreateLog(logEntry); err != nil {
				log.Printf("Error storing log in database: %v", err)
			}
		}

		c.processLogForAlerts(logEntry)
	}
}

func (c *LogConsumer) shouldStoreLog(l logger.Log) bool {
	if l.Level == logger.ERROR || l.Level == logger.WARN {
		return true
	}

	if l.Category == logger.SECURITY {
		return true
	}

	rate, exists := c.sampleRates[l.Level]
	if !exists {
		return false
	}

	return rand.Float64() < rate
}

func (c *LogConsumer) processLogForAlerts(l logger.Log) {
	if l.Level == logger.ERROR || (l.Level == logger.WARN && l.Category == logger.SECURITY) {
		fmt.Printf("ALERT: %s - %s - User: %s\n", l.Level, l.Message, l.Username)
	}
}

func (c *LogConsumer) Close() error {
	return c.reader.Close()
}

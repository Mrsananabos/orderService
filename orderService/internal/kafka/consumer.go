package eventHandler

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"orderService/configs"
	"orderService/internal/models"
	"orderService/internal/service"
	"time"
)

type Consumer struct {
	consumer     *kafka.Consumer
	orderService service.IOrderService
	retry        int
	backoff      time.Duration
}

const (
	ORDER_TOPIC = "Orders"
)

func CreateConsumer(cnf configs.Kafka, service service.IOrderService) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s%s", cnf.Host, cnf.Port),
		"group.id":           "1",
		"enable.auto.commit": false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	if err = consumer.Subscribe(ORDER_TOPIC, nil); err != nil {
		return nil, err
	}
	return &Consumer{consumer, service, cnf.Retry, time.Duration(cnf.Backoff) * time.Millisecond}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer start")
	for {
		select {
		case <-ctx.Done():
			if err := c.Stop(); err != nil {
				log.Printf("%v", err.Error())
			}
			log.Println("Kafka consumer Stopped")
			return
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Received message in %s topic: %s\n", *msg.TopicPartition.Topic, string(msg.Value))
				err = c.handleMessage(msg.Value)
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				log.Printf("Consumer error: %v\n", err)
			}
		}
	}
}

func (c *Consumer) Stop() error {
	if err := c.commitLastOffset(); err != nil {
		log.Printf("Kafka consumer Stop failed: %s\n", err.Error())
	}
	return c.consumer.Close()
}

func (c *Consumer) commitLastOffset() error {
	var err error
	backoff := c.backoff
	for attempt := 1; attempt <= c.retry; attempt++ {
		_, err = c.consumer.Commit()
		if err == nil {
			return nil
		}
		log.Printf("Error while commit kafka offset: %v. %dth attempt out of %d, wait %v\n", err, attempt, c.retry, backoff)
		<-time.After(backoff)
		backoff *= 2
	}

	return err
}

func (c *Consumer) handleMessage(msg []byte) error {
	var order models.Order
	if err := order.UnmarshalJSON(msg); err != nil {
		return err
	}

	return c.orderService.Create(order)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

var brokers = []string{"localhost:9092"}

func main() {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true

    topic := "test-topic"
    group := "test-group"

    admin, err := sarama.NewClusterAdmin(brokers, config)
    if err != nil {
        log.Fatalf("error while creating cluster admin: %v", err)
    }
    defer admin.Close()

    err = admin.CreateTopic(topic, &sarama.TopicDetail{
        NumPartitions:     1,
        ReplicationFactor: 1,
    }, false)
    if err != nil {
        if errors.Is(err, sarama.ErrTopicAlreadyExists) {
            log.Println("topic already exists")
        } else {
            log.Fatalf("error while creating topic: %v", err)
        }
    } else {
        log.Println("INFO: topic created")
    }

    go produceMessagesAsync(topic)

    consumeMessagesWithGroup(topic, group)
}

func NewAsyncProducer() (sarama.AsyncProducer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true
    producer, err := sarama.NewAsyncProducer(brokers, config)
    return producer, err
}

func produceMessagesAsync(topic string) {
    producer, err := NewAsyncProducer()
    if err != nil {
        log.Fatalf("error while creating producer: %v", err)
    }
    defer producer.AsyncClose()

    go func() {
        for success := range producer.Successes() {
            log.Printf("sent message to partition %d with offset %d", success.Partition, success.Offset)
        }
    }()

    go func() {
        for err := range producer.Errors() {
            log.Printf("error while sending message: %v", err)
        }
    }()

    for i := 0; i < 10; i++ {
        message := fmt.Sprintf("Message %d", i)
        producer.Input() <- &sarama.ProducerMessage{
            Topic: topic,
            Value: sarama.StringEncoder(message),
        }

        log.Printf("message sent: %s", message)
    }
}

type ConsumerGroupHandler struct {
    ready chan bool
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error  {
    close(h.ready)
    return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h * ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession,  claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {
        log.Printf("got message: topic=%s partition=%d offset=%d: %s", message.Topic, message.Partition, message.Offset, string(message.Value))
        session.MarkMessage(message, "")
    }
    return nil
}

func consumeMessagesWithGroup(topic string, group string) {
    consumer := ConsumerGroupHandler {
        ready: make(chan bool),
    }

    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true
    config.Consumer.Offsets.Initial = sarama.OffsetOldest

    client, err := sarama.NewConsumerGroup(brokers, group, config)
    if err != nil {
        log.Fatalf("error while creating consumer group: %v", err)
    }
    defer client.Close()

    go func() {
        for err := range client.Errors() {
            log.Printf("error in consumer group: %v", err)
        }
    }()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        for {
            err := client.Consume(ctx, []string{topic}, &consumer)
            if err != nil {
                log.Printf("error while consuming: %v", err)
            }

            if ctx.Err() != nil {
                return
            }
            consumer.ready = make(chan bool)
        }
    }()

    <- consumer.ready
    log.Println("consumer ready and joined the group")

    sigterm := make(chan os.Signal, 1)
    signal.Notify(sigterm, os.Interrupt)
    <-sigterm
    log.Println("got os interrupt. exiting")
    log.Println("done")
}

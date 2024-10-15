package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

var brokers = []string{"localhost:9092"}

func main() {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true

    topic := "test-topic"

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

    consumeMessages(topic)
}

func NewAsyncProducer() (sarama.AsyncProducer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true
    producer, err := sarama.NewAsyncProducer(brokers, config)
    return producer, err
}

func NewConsumer() (sarama.Consumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true
    consumer, err := sarama.NewConsumer(brokers, config)
    return consumer, err
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
        time.Sleep(500 * time.Millisecond)
    }
}

func consumeMessages(topic string) {
    consumer, err := NewConsumer()
    if err != nil {
        log.Fatalf("error while creating consumer: %v", err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
    if err != nil {
        log.Fatalf("error while creating partition consumer: %v", err)
    }
    defer partitionConsumer.Close()

    doneCh := make(chan struct{})
    go func() {
        for msg := range partitionConsumer.Messages() {
            log.Printf("got message: %s", string(msg.Value))
        }
        doneCh <- struct{}{}
    }()

    sigterm := make(chan os.Signal, 1)
    signal.Notify(sigterm, os.Interrupt)
    select {
    case <-sigterm:
        log.Println("got exit signal")
    case <-doneCh:
    }

    log.Println("done")
}

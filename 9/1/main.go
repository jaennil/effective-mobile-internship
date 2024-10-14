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

func main() {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Consumer.Return.Errors = true

    brokers := []string{"localhost:9092"}
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

    go produceMessages(brokers, topic, config)

    consumeMessages(brokers, topic, config)
}

func produceMessages(brokers []string, topic string, config *sarama.Config) {
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        log.Fatalf("error while creating producer: %v", err)
    }
    defer producer.Close()

    for i := 0; i < 10; i++ {
        msg := &sarama.ProducerMessage{
            Topic: topic,
            Value: sarama.StringEncoder(fmt.Sprintf("Message %d", i)),
        }
        partition, offset, err := producer.SendMessage(msg)
        if err != nil {
            log.Printf("error while sending message: %v", err)
        } else {
            log.Printf("message sent to partition %d with offset %d", partition, offset)
        }

        time.Sleep(500 * time.Millisecond)
    }
}

func consumeMessages(brokers []string, topic string, config *sarama.Config) {
    consumer, err := sarama.NewConsumer(brokers, config)
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

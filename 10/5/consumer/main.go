package main

import (
	"5/service"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

func main() {
    brokers := []string{"localhost:9092"}
    kafkaConsumer, err := service.InitKafkaConsumer(brokers)
    if err != nil {
        log.Printf("failed to init kafka consumer: %v", err)
    }
    defer kafkaConsumer.Close()

    partitionConsumer, err := kafkaConsumer.ConsumePartition("user-registration", 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatalf("failed to start consumer for partition: %v", err)
    }
    defer partitionConsumer.Close()

    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)

    log.Println("notification service started")

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            sendNotification(msg)
        case <-signals:
            log.Println("shutting down notification service")
            return
        }
    }
}

func sendNotification(message *sarama.ConsumerMessage) {
    fmt.Println("new notification!!!")
    fmt.Println(string(message.Value))
}

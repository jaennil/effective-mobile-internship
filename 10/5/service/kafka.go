package service

import "github.com/IBM/sarama"

func InitKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    return sarama.NewSyncProducer(brokers, config)
}

func InitKafkaConsumer(brokers []string) (sarama.Consumer, error) {
    return sarama.NewConsumer(brokers, nil)
}

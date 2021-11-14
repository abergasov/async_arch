package broker

import (
	"async_arch/internal/config"

	"github.com/segmentio/kafka-go"
)

func InitKafkaProducer(conf *config.BrokerConf, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(conf.Address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func InitKafkaConsumer(conf *config.BrokerConf, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{conf.Address},
		Topic:    topic,
		GroupID:  "tasker",
		MinBytes: 5,
		MaxBytes: 10e6,
		//Partition: 0,
	})
}

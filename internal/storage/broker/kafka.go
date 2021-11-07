package broker

import (
	"async_arch/internal/config"

	"github.com/segmentio/kafka-go"
)

func InitKafkaProducer(conf *config.BrokerConf, topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(conf.Address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return w
}

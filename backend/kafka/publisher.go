package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/imc-trading/ifwatch/backend"

	"github.com/Shopify/sarama"
)

type publisher struct {
	topic string
	sarama.SyncProducer
}

func NewPublisher(brokers []string, topic string) (*publisher, error) {
	//	cfg := sarama.NewConfig()
	//	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	p, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("new kafka producer: %v", err)
	}

	return &publisher{
		topic:        topic,
		SyncProducer: p,
	}, nil
}

func (p *publisher) Backend() backend.Backend {
	return backend.Kafka
}

func (p *publisher) Send(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal [%v]: %s", v, err)
	}

	_, _, err = p.SendMessage(&sarama.ProducerMessage{Topic: p.topic, Key: sarama.ByteEncoder([]byte(k)), Value: sarama.ByteEncoder(b)})
	if err != nil {
		return fmt.Errorf("kafka produce message [%s]: %v", k, err)
	}
	return nil
}

func (p *publisher) Close() error {
	return p.Close()
}

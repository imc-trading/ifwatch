package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/imc-trading/ifwatch/backend"

	"github.com/Shopify/sarama"
	"github.com/mickep76/compress"
	_ "github.com/mickep76/compress/gzip"
	_ "github.com/mickep76/compress/lzw"
	_ "github.com/mickep76/compress/snappy"
	_ "github.com/mickep76/compress/xz"
	_ "github.com/mickep76/compress/zlib"
)

type publisher struct {
	topic string
	algo  compress.Algorithm
	sarama.SyncProducer
}

type PublisherOption publisher

func NewPublisher(brokers []string, topic string, algo string) (*publisher, error) {
	p := &publisher{
		topic: topic,
	}

	var err error
	if algo != "none" {
		if p.algo, err = compress.NewAlgorithm(algo); err != nil {
			return nil, err
		}
	}

	if p.SyncProducer, err = sarama.NewSyncProducer(brokers, nil); err != nil {
		return nil, fmt.Errorf("new kafka producer: %v", err)
	}

	return p, nil
}

func (p *publisher) Backend() backend.Backend {
	return backend.Kafka
}

func (p *publisher) Send(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal [%v]: %s", v, err)
	}

	if p.algo != nil {
		enc, err := p.algo.Encode(b)
		if err != nil {
			return fmt.Errorf("compress [%v]: %s", v, err)
		}

		_, _, err = p.SendMessage(&sarama.ProducerMessage{Topic: p.topic, Key: sarama.ByteEncoder([]byte(k)), Value: sarama.ByteEncoder(enc)})
		if err != nil {
			return fmt.Errorf("kafka produce message [%s]: %v", k, err)
		}
	} else {
		_, _, err = p.SendMessage(&sarama.ProducerMessage{Topic: p.topic, Key: sarama.ByteEncoder([]byte(k)), Value: sarama.ByteEncoder(b)})
		if err != nil {
			return fmt.Errorf("kafka produce message [%s]: %v", k, err)
		}
	}

	return nil
}

func (p *publisher) Close() error {
	return p.Close()
}

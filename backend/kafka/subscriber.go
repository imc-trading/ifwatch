package kafka

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/imc-trading/ifwatch/backend"

	"github.com/Shopify/sarama"
)

type subscriber struct {
	topic    string
	handlers []backend.MessageHandler
	wait     sync.WaitGroup
	done     chan interface{}
	client   sarama.Client
	sarama.Consumer
}

func NewSubscriber(brokers []string, topic string) (*subscriber, error) {
	cons, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	cli, err := sarama.NewClient(brokers, sarama.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("new kafka client: %v", err)
	}

	return &subscriber{
		topic:    topic,
		wait:     sync.WaitGroup{},
		done:     make(chan interface{}),
		client:   cli,
		Consumer: cons,
	}, nil
}

func (s *subscriber) Backend() backend.Backend {
	return backend.Kafka
}

func (s *subscriber) AddHandler(handler backend.MessageHandler) {
	s.handlers = append(s.handlers, handler)
}

func (s *subscriber) getPartition(k string) (int32, error) {
	parts, err := s.Partitions(s.topic)
	if err != nil {
		return 0, fmt.Errorf("get partitions for topic [%s]: %v", s.topic, err)
	}
	numParts := int32(len(parts))

	p := sarama.NewHashPartitioner(s.topic)
	msg := sarama.ProducerMessage{Topic: s.topic, Key: sarama.ByteEncoder([]byte(k)), Value: sarama.ByteEncoder([]byte{})}
	return p.Partition(&msg, numParts)
}

func (s *subscriber) Versions(k string, handler backend.MessageHandler) error {
	part, err := s.getPartition(k)
	if err != nil {
		return err
	}

	// Exit if it's an empty partition
	off, _ := s.client.GetOffset(s.topic, part, sarama.OffsetNewest)
	if off == 0 {
		return nil
	}

	pc, err := s.ConsumePartition(s.topic, part, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("start subscriber for topic [%s]: %v", s.topic, err)
	}

	defer pc.Close()

	for m := range pc.Messages() {
		if string(m.Key) == k {
			handler(string(m.Key), m.Value)
		}
		if m.Offset+1 == pc.HighWaterMarkOffset() {
			return nil
		}
	}

	return nil
}

func (s *subscriber) Keys() ([]string, error) {
	return nil, errors.New("not supported")
}

func (s *subscriber) Start() error {
	parts, err := s.Partitions(s.topic)
	if err != nil {
		return fmt.Errorf("get partitions for topic [%s]: %v", s.topic, err)
	}

	for _, part := range parts {
		pc, err := s.ConsumePartition(s.topic, part, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("start subscriber for topic [%s]: %v", s.topic, err)
		}

		s.wait.Add(1)
		go func() {
			defer s.wait.Done()

		Loop:
			for {
				select {
				case m := <-pc.Messages():
					for _, handler := range s.handlers {
						go handler(string(m.Key), m.Value)
					}
				case err := <-pc.Errors():
					log.Printf("consume message: %v", err)
				case <-s.done:
					pc.Close()
					break Loop
				}
			}
		}()
	}

	s.wait.Wait()

	return nil
}

func (s *subscriber) Stop() error {
	close(s.done)
	s.wait.Wait()

	return s.Close()
}

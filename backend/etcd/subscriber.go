package etcd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imc-trading/ifwatch/backend"

	"github.com/coreos/etcd/clientv3"
)

type subscriber struct {
	prefix   string
	handlers []backend.MessageHandler
	*clientv3.Client
}

func NewSubscriber(endpoints []string, tmout int, prefix string) (*subscriber, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(tmout) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd conn: %v", err)
	}

	return &subscriber{
		prefix: prefix,
		Client: cli,
	}, nil
}

func (s *subscriber) Backend() backend.Backend {
	return backend.Etcd
}

func (s *subscriber) AddHandler(handler backend.MessageHandler) {
	s.handlers = append(s.handlers, handler)
}

func (s *subscriber) Versions(k string, handler backend.MessageHandler) error {
	resp, err := s.Client.Get(context.TODO(), k, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("get etcd key [%s]: %v", k, err)
	}

	if len(resp.Kvs) == 0 {
		return nil
	}

	firstRev := resp.Kvs[0].Version
	lastRev := resp.Kvs[0].Version
	lastVers := resp.Kvs[0].Version
	prev := map[string]int64{}
	for _, kv := range resp.Kvs {
		if firstRev > kv.CreateRevision {
			firstRev = kv.CreateRevision
		}
		if lastRev < kv.ModRevision {
			lastRev = kv.ModRevision
		}
		if lastVers < kv.Version {
			lastVers = kv.Version
		}
		prev[string(kv.Key)] = 0
	}

	for rev := firstRev; rev <= lastRev; rev++ {
		resp, err := s.Client.Get(context.TODO(), k, clientv3.WithPrefix(), clientv3.WithRev(rev))
		if err != nil {
			return fmt.Errorf("get etcd key [%s] rev [%d]: %v", k, rev, err)
		}

		for _, kv := range resp.Kvs {
			if prev[string(kv.Key)] < kv.Version {
				handler(string(kv.Key), kv.Value)
				prev[string(kv.Key)] = kv.Version
			}
		}
	}

	return nil
}

func (s *subscriber) Keys() ([]string, error) {
	resp, err := s.Client.Get(context.TODO(), s.prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, fmt.Errorf("get etcd keys: %v", err)
	}

	keys := []string{}
	for _, kv := range resp.Kvs {
		keys = append(keys, strings.TrimPrefix(string(kv.Key), s.prefix+"/"))
	}

	return keys, nil
}

func (s *subscriber) Start() error {
	rch := s.Client.Watch(context.Background(), s.prefix, clientv3.WithPrefix())
	for resp := range rch {
		for _, ev := range resp.Events {
			for _, handler := range s.handlers {
				handler(string(ev.Kv.Key), ev.Kv.Value)
			}
		}
	}

	return nil
}

func (s *subscriber) Stop() error {
	return s.Close()
}

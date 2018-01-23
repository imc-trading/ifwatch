package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imc-trading/ifwatch/backend"

	"github.com/coreos/etcd/clientv3"
)

type publisher struct {
	prefix string
	*clientv3.Client
}

func NewPublisher(endpoints []string, tmout int, prefix string) (*publisher, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(tmout) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("new etcd client: %v", err)
	}

	return &publisher{
		prefix: prefix,
		Client: cli,
	}, nil
}

func (p *publisher) Backend() backend.Backend {
	return backend.Etcd
}

func (p *publisher) Send(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal [%v]: %s", v, err)
	}

	kvc := clientv3.NewKV(p.Client)

	if _, err := kvc.Put(context.TODO(), k, string(b)); err != nil {
		return fmt.Errorf("etcd put key [%s]: %v", k, err)
	}
	return nil
}

func (p *publisher) Close() error {
	return p.Close()
}

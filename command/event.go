package command

import (
	"time"

	"github.com/imc-trading/ifwatch/network"
)

type EventType string

const (
	ActionAdd    = EventType("add")
	ActionModify = EventType("modify")
	ActionDelete = EventType("delete")
)

type Event struct {
	Created time.Time `json:"created"`
	Action  EventType `json:"action"`
	Host    string    `json:"host"`
	*network.Interface
}

package event

import (
	"time"

	"github.com/imc-trading/ifwatch/network"
)

type EventType string

const (
	ActionAdd         = EventType("add")
	ActionModify      = EventType("modify")
	ActionRefresh     = EventType("refresh")
	ActionDelete      = EventType("delete")
	ActionLimitStop   = EventType("limit-stop")
	ActionLimitResume = EventType("limit-resume")
)

type Event struct {
	Created time.Time `json:"created"`
	Action  EventType `json:"action"`
	Host    string    `json:"host"`
	*network.Interface
}

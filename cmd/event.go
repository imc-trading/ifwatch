package cmd

import (
	"fmt"
	"time"

	"github.com/imc-trading/ifwatch/netx"
)

type EventType string

const (
	ActionAdd     = EventType("add")
	ActionModify  = EventType("modify")
	ActionDelete  = EventType("delete")
	ActionRefresh = EventType("refresh")
)

type Event struct {
	Created time.Time `json:"created"`
	Action  EventType `json:"action"`
	Host    string    `json:"host"`
	*netx.Interface
}

type EventList []*Event
type EventMap map[string]*Event

func (e *Event) String() string {
	return fmt.Sprintf("%s/%s:\n", e.Host, e.Name) +
		fmt.Sprintf("\t%-58s : %s\n", "Created", e.Created.Format("2006-01-02 15:04:05")) +
		fmt.Sprintf("\t%-58s : %v", "Action", e.Action) + e.Interface.String()
}

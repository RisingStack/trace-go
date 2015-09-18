package trace

import (
	"time"

	"github.com/twinj/uuid"
)

const (
	RequestReceived = iota
)

type Event struct {
	RequestId string
	CreatedAt time.Time
	Type      int
}

func NewEvent(t int) Event {
	return Event{
		RequestId: uuid.NewV4().String(),
		Type:      t,
		CreatedAt: time.Now(),
	}
}

package trace

import (
	"log"
	"sync"
)

const NumberOfEventsBeforeFlush = 1000

type Collector interface {
	Record(e Event) error
}

type MemoryCollector struct {
	list []Event
	in   chan Event
	lock sync.Mutex
}

func NewMemoryCollector() *MemoryCollector {
	c := MemoryCollector{}
	c.list = make([]Event, 0, NumberOfEventsBeforeFlush)
	c.in = make(chan Event)
	c.lock = sync.Mutex{}
	go func() {
		for {
			e := <-c.in
			c.lock.Lock()
			c.list = append(c.list, e)
			c.lock.Unlock()
		}
	}()
	return &c
}

func (c *MemoryCollector) Record(event Event) error {
	c.in <- event
	return nil
}

func (c *MemoryCollector) GetEvents() []Event {
	c.lock.Lock()
	defer c.lock.Unlock()
	listLength := len(c.list)
	newList := make([]Event, listLength)
	copied := copy(newList, c.list)
	if copied != listLength {
		log.Panicln("Failed to copy Events from list")
	}
	return newList
}

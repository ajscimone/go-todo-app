package main

import (
	"sync"
)

type Repository interface {
	Get(id int) (Event, bool)
	GetAll() []Event
	Set(e Event) error
	Delete(id int) error
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{events: make([]Event, 0)}
}

type InMemRepo struct {
	sync.Mutex
	events []Event
}

func (imr *InMemRepo) Get(id int) (Event, bool) {
	for _, e := range imr.events {
		if e.ID == id {
			return e, true
		}
	}
	return Event{}, false
}

func (imr *InMemRepo) GetAll() []Event {
	return imr.events
}

func (imr *InMemRepo) Set(e Event) error {
	imr.Lock()
	defer imr.Unlock()

	existing, ok := imr.Get(e.ID)
	if ok {
		err := imr.Delete(existing.ID)
		if err != nil {
			return err
		}
	}

	imr.events = append(imr.events, e)
	return nil
}

func (imr *InMemRepo) Delete(id int) error {
	imr.Lock()
	defer imr.Unlock()

	for idx, e := range imr.events {
		if e.ID == id {
			imr.events = append(imr.events[:idx], imr.events[idx:]...)
		}
	}
	return nil
}

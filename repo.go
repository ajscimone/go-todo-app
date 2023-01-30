package main

import (
	"sync"
)

type repository interface {
	Get(id int) (event, bool)
	GetAll() []event
	Set(e event) error
	Delete(id int) error
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{events: make([]event, 0)}
}

type InMemRepo struct {
	sync.Mutex
	events []event
}

func (imr *InMemRepo) Get(id int) (event, bool) {
	for _, e := range imr.events {
		if e.ID == id {
			return e, true
		}
	}
	return event{}, false
}

func (imr *InMemRepo) GetAll() []event {
	return imr.events
}

func (imr *InMemRepo) Set(e event) error {
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

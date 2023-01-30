package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"strconv"

	"github.com/gorilla/mux"
)

type autoInc struct {
    sync.Mutex
    id int
}

func (a *autoInc) ID() (id int) {
    a.Lock()
    defer a.Unlock()

    id = a.id
    a.id++
    return
}

var unique autoInc

type event struct {
	ID          int		`json:"ID"`
	Title       string	`json:"Title"`
	Description string	`json:"Description"`
}

var events = []event{
	{
		ID:          unique.ID(),
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	newEvent.ID = unique.ID()

	reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventIDParam := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(eventIDParam)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventIDParam := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(eventIDParam)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
			return
		}
	}
	
	http.Error(w, err.Error(), http.StatusNotFound)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventIDParam := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(eventIDParam)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
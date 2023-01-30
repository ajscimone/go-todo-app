package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type autoInc struct {
	sync.Mutex
	id int
}

func (a *autoInc) ID() int {
	a.Lock()
	defer a.Unlock()

	a.id++
	return a.id
}

type event struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type Handler struct {
	counter autoInc
	repo    repository
}

func NewHandler(repo repository) *Handler {
	return &Handler{counter: autoInc{}, repo: repo}
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	newEvent.ID = h.counter.ID()

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.Unmarshal(reqBody, &newEvent)
	err = h.repo.Set(newEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func (h *Handler) getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventIDParam, ok := mux.Vars(r)["id"]
	if !ok {

	}
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e, ok := h.repo.Get(eventID)
	if !ok {
		http.Error(w, "event not found", http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(e)
}

func (h *Handler) getAllEvents(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.repo.GetAll())
}

func (h *Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	eventIDParam := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := h.repo.Get(eventID)
	if !ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedEvent event
	err = json.Unmarshal(reqBody, &updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.repo.Set(updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedEvent)
	return
}

func (h *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventIDParam := mux.Vars(r)["id"]
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	rep := NewInMemRepo()
	h := NewHandler(rep)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", h.createEvent).Methods("POST")
	router.HandleFunc("/events", h.getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", h.getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", h.updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", h.deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

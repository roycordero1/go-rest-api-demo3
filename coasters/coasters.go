package coasters

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Coaster struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID           string `json:"id"`
	InPark       string `json:"in_park"`
	Height       int    `json:"height"`
}

type coastersHandler struct {
	sync.Mutex
	store map[string]Coaster
}

func NewCoastersHandler() *coastersHandler {
	return &coastersHandler{
		store: map[string]Coaster{},
	}
}

func (h *coastersHandler) ListCoasters(w http.ResponseWriter, r *http.Request) {

	// List the coasters
	coasters := make([]Coaster, len(h.store))

	h.Lock()
	i := 0
	for _, coaster := range h.store {
		coasters[i] = coaster
		i++
	}
	h.Unlock()

	// Write body as json to return response
	jsonBytes, err := json.Marshal(coasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *coastersHandler) GetCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Search for existing coaster
	h.Lock()
	coaster, ok := h.store[id]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write body as json to return response
	jsonBytes, err := json.Marshal(coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *coastersHandler) CreateCoaster(w http.ResponseWriter, r *http.Request) {

	// Read body to create the coaster
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json' but got '%s'", ct)))
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var coaster Coaster
	err = json.Unmarshal(bodyBytes, &coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Save the coaster
	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[coaster.ID] = coaster
	defer h.Unlock()
}

func (h *coastersHandler) GetRandomCoaster(w http.ResponseWriter, r *http.Request) {

	ids := make([]string, len(h.store))

	h.Lock()
	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header().Add("location", fmt.Sprintf("/api/v1/coasters/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *coastersHandler) UpdateCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Search for existing coaster
	h.Lock()
	_, ok := h.store[id]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Read body to update the coaster
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json' but got '%s'", ct)))
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var newCoaster Coaster
	err = json.Unmarshal(bodyBytes, &newCoaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Save the coaster
	newCoaster.ID = id
	h.Lock()
	h.store[id] = newCoaster
	h.Unlock()
}

func (h *coastersHandler) DeleteCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Search for existing coaster
	h.Lock()
	_, ok := h.store[id]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the coaster
	h.Lock()
	delete(h.store, id)
	h.Unlock()
}

package coasters

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"time"
)

type Coaster struct {
	ID           string `json:"id" db:", primarykey"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	InPark       string `json:"in_park" db:"in_park"`
	Height       int    `json:"height"`
}

type coastersHandler struct {
	DB *gorp.DbMap
}

const (
	SqlListCoasters = "SELECT * FROM coasters"
	SqlGetCoasterById = "SELECT * FROM coasters WHERE id = ?"
)

func NewCoastersHandler(db *gorp.DbMap) *coastersHandler {
	return &coastersHandler{
		DB: db,
	}
}

func (h *coastersHandler) ListCoasters(w http.ResponseWriter, r *http.Request) {

	// List the coasters
	var coasters []Coaster
	_, err := h.DB.Select(&coasters, SqlListCoasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write body as json to return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coasters)
}

func (h *coastersHandler) GetCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Search for existing coaster
	var coaster Coaster
	if err := h.DB.SelectOne(&coaster, SqlGetCoasterById, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Coaster Not Found!"))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// Write body as json to return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coaster)
}

func (h *coastersHandler) GetRandomCoaster(w http.ResponseWriter, r *http.Request) {

	// List the coasters
	// List the coasters
	var coasters []Coaster
	_, err := h.DB.Select(&coasters, SqlListCoasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var target string
	if len(coasters) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Coasters Found!"))
		return
	} else if len(coasters) == 1 {
		target = coasters[0].ID
	} else {
		rand.Seed(time.Now().UnixNano())
		target = coasters[rand.Intn(len(coasters))].ID
	}

	w.Header().Add("location", fmt.Sprintf("/api/v1/coasters/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *coastersHandler) CreateCoaster(w http.ResponseWriter, r *http.Request) {

	// Read body to create the coaster
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json' but got '%s'", ct)))
		return
	}

	var coaster Coaster
	err := json.NewDecoder(r.Body).Decode(&coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Save the coaster
	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	err = h.DB.Insert(&coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write body as json to return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(coaster)
}

func (h *coastersHandler) UpdateCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Read body to update the coaster
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json' but got '%s'", ct)))
		return
	}

	var coaster Coaster
	err := json.NewDecoder(r.Body).Decode(&coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Save the coaster
	coaster.ID = id
	count, err := h.DB.Update(&coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Coaster Not Found or Data Not Changing!"))
		return
	}

	// Write body as json to return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coaster)
}

func (h *coastersHandler) DeleteCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Search for existing coaster
	var coaster Coaster
	if err := h.DB.SelectOne(&coaster, SqlGetCoasterById, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Coaster Not Found!"))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// Delete the coaster
	_, err := h.DB.Delete(&coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coaster Deleted!"))
}

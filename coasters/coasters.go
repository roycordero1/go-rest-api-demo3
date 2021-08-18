package coasters

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"time"
)

type Coaster struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	InPark       string `json:"in_park"`
	Height       int    `json:"height"`
}

type coastersHandler struct {
	DB *sql.DB
}

func NewCoastersHandler(db *sql.DB) *coastersHandler {
	return &coastersHandler{
		DB: db,
	}
}

func (h *coastersHandler) ListCoasters(w http.ResponseWriter, r *http.Request) {

	// List the coasters
	rows, err := h.DB.Query("SELECT id, name, manufacturer, in_park, height FROM coasters")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	coasters := []Coaster{}
	for rows.Next() {
		var coaster Coaster
		if err = rows.Scan(&coaster.ID, &coaster.Name, &coaster.Manufacturer, &coaster.InPark, &coaster.Height); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		coasters = append(coasters, coaster)
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
	row := h.DB.QueryRow("SELECT id, name, manufacturer, in_park, height FROM coasters WHERE id = ?", id)

	var coaster Coaster
	if err := row.Scan(&coaster.ID, &coaster.Name, &coaster.Manufacturer, &coaster.InPark, &coaster.Height); err != nil {
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
	rows, err := h.DB.Query("SELECT id, name, manufacturer, in_park, height FROM coasters")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var coasters []Coaster
	for rows.Next() {
		var coaster Coaster
		if err = rows.Scan(&coaster.ID, &coaster.Name, &coaster.Manufacturer, &coaster.InPark, &coaster.Height); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		coasters = append(coasters, coaster)
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
	statement, err := h.DB.Prepare("INSERT INTO coasters(id, name, manufacturer, in_park, height) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = statement.Exec(coaster.ID, coaster.Name, coaster.Manufacturer, coaster.InPark, coaster.Height)
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
	statement, err := h.DB.Prepare("UPDATE coasters SET name = ?, manufacturer = ?, in_park = ?, height = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := statement.Exec(coaster.Name, coaster.Manufacturer, coaster.InPark, coaster.Height, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if qtyRowsAffected, _ := result.RowsAffected(); qtyRowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Coaster Not Found or Data Not Changing!"))
		return
	}

	// Write body as json to return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	coaster.ID = id
	json.NewEncoder(w).Encode(coaster)
}

func (h *coastersHandler) DeleteCoaster(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	// Delete the coaster
	statement, err := h.DB.Prepare("DELETE FROM coasters WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := statement.Exec(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if qtyRowsAffected, _ := result.RowsAffected(); qtyRowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Coaster Not Found!"))
		return
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coaster Deleted!"))
}

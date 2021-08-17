package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"rest-api-tutorial2/admin"
	"rest-api-tutorial2/coasters"
)

func main() {

	admin := admin.NewAdminHandler()
	coastersHandler := coasters.NewCoastersHandler()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/admin", admin.Handler)
	router.HandleFunc("/api/v1/coasters", coastersHandler.ListCoasters).Methods("GET")
	router.HandleFunc("/api/v1/coasters/random", coastersHandler.GetRandomCoaster).Methods("GET")
	router.HandleFunc("/api/v1/coasters/{id}", coastersHandler.GetCoaster).Methods("GET")
	router.HandleFunc("/api/v1/coasters", coastersHandler.CreateCoaster).Methods("POST")
	router.HandleFunc("/api/v1/coasters/{id}", coastersHandler.UpdateCoaster).Methods("PUT")
	router.HandleFunc("/api/v1/coasters/{id}", coastersHandler.DeleteCoaster).Methods("DELETE")

	http.ListenAndServe(":8081", router)
}

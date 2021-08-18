package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"rest-api-tutorial2/admin"
	"rest-api-tutorial2/coasters"
	"time"
)

var _ = godotenv.Load(".env")
var (
	connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("user"),
		os.Getenv("pass"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("db_name"))
)

type Env struct {
	DB *sql.DB
}

func (env *Env) initDatabase() error {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("Error with database. The error is: " + err.Error())
		return err
	} else {
		err = db.Ping()
		if err != nil {
			log.Printf("Error communicating to database, please check credentials. The error is: " + err.Error())
			return err
		}
	}

	log.Printf("Database connected!")

	env.DB = db
	return nil
}

func main() {

	env := Env{}
	err := env.initDatabase()
	if err != nil {
		log.Fatal("Error initializing the database. The error is: " + err.Error())
		return
	}

	admin := admin.NewAdminHandler()
	coastersHandler := coasters.NewCoastersHandler(env.DB)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/admin", admin.Handler)
	router.HandleFunc("/api/v1/coasters", coastersHandler.ListCoasters).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/coasters/random", coastersHandler.GetRandomCoaster).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/coasters/{id}", coastersHandler.GetCoaster).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/coasters", coastersHandler.CreateCoaster).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/coasters/{id}", coastersHandler.UpdateCoaster).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/coasters/{id}", coastersHandler.DeleteCoaster).Methods(http.MethodDelete)

	port := ":8081"
	server := &http.Server{
		Handler: router,
		Addr: port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Printf("Server started at %s", port)
	log.Fatal(server.ListenAndServe())
}

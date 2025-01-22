package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:root@(127.0.0.1:3333)/testdb")

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	initDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)

	http.ListenAndServe(":3000", router)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var version string

	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		log.Fatal(err)
	}

	w.Write([]byte("db version: " + version))
}

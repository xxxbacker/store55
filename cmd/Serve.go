package main

import (
	"Store55/internal"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

func loadConfig() *internal.Config {
	config := internal.NewConfig()
	yamlFile, err := os.ReadFile("./Config/local.yaml")
	if err != nil {
		log.Fatal("щшибка чтения конфиг файла", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal("ошибка анмаршелизации")
	}
	return config
}

func main() {

	config := loadConfig()

	db, err := sql.Open("postgres", config.Db_string)
	if err != nil {
		log.Fatal(err)
	}
	internal.Database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", internal.IndexHandler)
	router.HandleFunc("/create", internal.CreateHandler)
	router.HandleFunc("/edit/{id:[0-9]+}", internal.EditPage).Methods("GET")
	router.HandleFunc("/edit/{id:[0-9]+}", internal.EditHandler).Methods("POST")
	router.HandleFunc("/delete/{id:[0-9]+}", internal.DeleteHandler)

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(config.Address, nil)
}

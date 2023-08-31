package internal

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Product struct {
	Id      int
	Model   string
	Company string
	Price   int
}

var Database *sql.DB

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := Database.Exec("delete from mytable where id = $1", id)
	if err != nil {
		log.Println("26 строка", err)
	}

	http.Redirect(w, r, "/", 301)
}

func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := Database.QueryRow("select * from mytable where id = $1", id)
	prod := Product{}
	err := row.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
	if err != nil {
		log.Println("40 строка", err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit.html")
		tmpl.Execute(w, prod)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	model := r.FormValue("model")
	company := r.FormValue("company")
	price := r.FormValue("price")

	_, err = Database.Exec("update mytable set model=$1, company=$2, price = $3 where id = $4", model, company, price, id)

	if err != nil {
		log.Println("60 строка", err)
	}
	http.Redirect(w, r, "/", 301)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		model := r.FormValue("model")
		company := r.FormValue("company")
		price := r.FormValue("price")

		_, err = Database.Exec("INSERT INTO mytable (model, company, price) VALUES ($1, $2, $3)", model, company, price)

		if err != nil {
			log.Println("80 строка", err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := Database.Query("select * from mytable")
	if err != nil {
		log.Println("94 строка", err)
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println("103 строка", err)
			continue
		}
		products = append(products, p)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, products)
}

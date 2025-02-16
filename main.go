package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("template/*.html"))
}

func initDB() {
	var err error

	db, err = sql.Open("mysql", "root:root@(127.0.0.1:3333)/testdb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	gRouter := mux.NewRouter()

	gRouter.HandleFunc("/", HomeHandler)

	http.ListenAndServe(":3000", gRouter)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(),
			http.StatusInternalServerError)
	}
}

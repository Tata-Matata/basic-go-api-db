package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func (app *App) Initialize() error {
	connStr := fmt.Sprintf("%v:%v@tcp(localhost:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", connStr)
	checkErr(err)

	app.Router = mux.NewRouter().StrictSlash(true)
	app.HandleRoutes()
	return nil

}

func sendResponse(respWriter http.ResponseWriter, status int, payload interface{}) {

	content, _ := json.Marshal(payload)

	respWriter.Header().Set("Content-type", "application/json")
	respWriter.WriteHeader(status)
	respWriter.Write(content)
}

func sendError(respWriter http.ResponseWriter, status int, err string) {
	message := map[string]string{"error": err}

	sendResponse(respWriter, status, message)

}

func (app *App) getProducts(respWriter http.ResponseWriter, req *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(respWriter, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(respWriter, http.StatusOK, products)
}

func (app *App) getProduct(respWriter http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(respWriter, http.StatusBadRequest, "invalid product ID")
		return
	}
	product, err := getProduct(app.DB, id)
	if err != nil || product == nil {
		sendError(respWriter, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(respWriter, http.StatusOK, product)
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func (app *App) HandleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
}

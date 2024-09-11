package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/cryptotechgeorgia/mocker/payload"
	payloadstorer "github.com/cryptotechgeorgia/mocker/payload/store"
	"github.com/cryptotechgeorgia/mocker/project"
	projstorer "github.com/cryptotechgeorgia/mocker/project/store"
	"github.com/cryptotechgeorgia/mocker/request"
	reqstorer "github.com/cryptotechgeorgia/mocker/request/store"
	"github.com/cryptotechgeorgia/mocker/response"
	respstorer "github.com/cryptotechgeorgia/mocker/response/store"
	"github.com/cryptotechgeorgia/mocker/router"

	"github.com/cryptotechgeorgia/mocker/web/controllers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates/*.html
var templateFS embed.FS

func main() {
	// Initialize DB
	db, err := sqlx.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS project (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		base_addr TEXT
	);
	CREATE TABLE IF NOT EXISTS request (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		path TEXT,
		method TEXT
	);

	CREATE TABLE IF NOT EXISTS request_payload (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		payload TEXT,
		request_id INTEGER,
		content_type TEXT
	);
	CREATE TABLE IF NOT EXISTS response (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		payload TEXT,
		request_payload_id INTEGER,
		content_type TEXT
	);`

	db.MustExec(schema)

	applyRouterChan := make(chan struct{})
	projStorer := projstorer.NewRepo(db)
	projBus := project.NewBusiness(projStorer)

	reqStorer := reqstorer.NewRepo(db)
	reqBus := request.NewBusiness(reqStorer)

	payloadStorer := payloadstorer.NewRepo(db)
	payloadBus := payload.NewBusiness(payloadStorer)

	projectHandler := controllers.NewProjectHandler(projBus, reqBus, applyRouterChan, templateFS)
	respStorer := respstorer.NewRepo(db)

	respBus := response.NewBusiness(respStorer)
	requestHandler := controllers.NewRuesthandler(reqBus, respBus, projBus, payloadBus, templateFS)

	webMux := mux.NewRouter()
	// project
	webMux.HandleFunc("/", projectHandler.ListProjects).Methods("GET")
	webMux.HandleFunc("/projects", projectHandler.AddProject).Methods("POST")
	webMux.HandleFunc("/projects/{id}", projectHandler.ViewProject).Methods("GET")

	// request
	webMux.HandleFunc("/projects/{id}/requests/add", requestHandler.AddRequest).Methods("POST")
	webMux.HandleFunc("/projects/{id}/requests/{reqId}", requestHandler.ViewRequest).Methods("GET")
	webMux.HandleFunc("/projects/{id}/requests/{reqId}/pair/add", requestHandler.AddPair).Methods("POST")
	webMux.HandleFunc("/projects/{id}/requests/{reqId}/pair/{pairId}/delete", requestHandler.RemovePair).Methods("POST")
	webMux.HandleFunc("/projects/apply", projectHandler.ApplyProjects).Methods("POST")

	// router.router.HandleFunc("/projects/{id}/delete")

	go func() {
		log.Fatal(http.ListenAndServe(":8080", webMux))
	}()

	// router building

	mainRouter := router.NewMainRouter(
		applyRouterChan,
		projBus,
		reqBus,
		respBus,
		payloadBus,
		webMux,
	)
	ch := make(chan error)

	// apply projects on start

	go mainRouter.Listen(context.Background(), ch)

	applyRouterChan <- struct{}{}

	for e := range ch {
		fmt.Printf("errors %+v", e)

	}

}

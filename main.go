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
	"github.com/gorilla/mux"

	"github.com/cryptotechgeorgia/mocker/web/controllers"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates/*.html
var templateFS embed.FS

func main() {
	// Initialize DB
	db, err := sqlx.Open("sqlite3", "./data/mocker.db")
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

	projStorer := projstorer.NewRepo(db)
	projBus := project.NewBusiness(projStorer)

	reqStorer := reqstorer.NewRepo(db)
	reqBus := request.NewBusiness(reqStorer)

	payloadStorer := payloadstorer.NewRepo(db)
	payloadBus := payload.NewBusiness(payloadStorer)

	applyRouterChan := make(chan struct{}, 1)
	errChan := make(chan error)

	projectHandler := controllers.NewProjectHandler(projBus, reqBus, applyRouterChan, templateFS)
	respStorer := respstorer.NewRepo(db)
	respBus := response.NewBusiness(respStorer)
	requestHandler := controllers.NewRuesthandler(reqBus, respBus, projBus, payloadBus, templateFS)

	webMux := mux.NewRouter()

	webMux.HandleFunc("/", projectHandler.ListProjects).Methods("GET")
	webMux.HandleFunc("/projects", projectHandler.AddProject).Methods("POST")
	webMux.HandleFunc("/projects/{id}", projectHandler.ViewProject).Methods("GET")

	// request
	webMux.HandleFunc("/projects/{id}/requests/add", requestHandler.AddRequest).Methods("POST")
	webMux.HandleFunc("/projects/{id}/requests/{reqId}", requestHandler.ViewRequest).Methods("GET")
	webMux.HandleFunc("/projects/{id}/requests/{reqId}/pair/add", requestHandler.AddPair).Methods("POST")
	webMux.HandleFunc("/projects/{id}/requests/{reqId}/pair/{pairId}/delete", requestHandler.RemovePair).Methods("POST")
	webMux.HandleFunc("/projects/apply", projectHandler.ApplyChanges).Methods("POST")

	populator := router.NewPopulator(projBus, payloadBus, respBus, reqBus)
	mockerRouter := router.NewMockerRouter(
		applyRouterChan,
		populator,
		webMux,
	)

	go func() {
		fmt.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", mockerRouter); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// apply projects on start
	go mockerRouter.Listen(context.Background(), errChan)

	applyRouterChan <- struct{}{}

	for e := range errChan {
		fmt.Printf("errors %+v", e)

	}
}

package controllers

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	cnv "github.com/cryptotechgeorgia/mocker/foundation/convert"
	"github.com/cryptotechgeorgia/mocker/project"
	"github.com/cryptotechgeorgia/mocker/request"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	bus           *project.Bussiness
	reqBus        *request.Bussiness
	applyChan     chan struct{}
	doneApplyChan chan struct{}
	tmpl          embed.FS
}

func NewProjectHandler(bus *project.Bussiness, req *request.Bussiness, applyChan chan struct{}, doneApplyChan chan struct{}, tmpl embed.FS) ProjectHandler {
	return ProjectHandler{
		bus:           bus,
		reqBus:        req,
		applyChan:     applyChan,
		doneApplyChan: doneApplyChan,
		tmpl:          tmpl,
	}
}

func (p *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projs, err := p.bus.All(r.Context())
	if err != nil {
		http.Error(w, "Error fetching projects", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(p.tmpl, "templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}

	data := struct {
		Projects []project.Project
	}{
		Projects: projs,
	}

	tmpl.Execute(w, data)
}
func (p *ProjectHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	name := r.FormValue("name")
	baseAddr := r.FormValue("base_addr")

	// Create a new project object
	newProject := project.Project{
		Name:     name,
		BaseAddr: baseAddr,
	}

	if err := p.bus.Add(r.Context(), newProject); err != nil {
		http.Error(w, "Failed to add project", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (p *ProjectHandler) ViewProject(w http.ResponseWriter, r *http.Request) {
	// Get the project ID from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// get project
	proj, err := p.bus.Get(r.Context(), id)

	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	requests, err := p.reqBus.Filter(r.Context(), request.FilterBy{
		ProjectId: cnv.ToIntPtr(id),
	})

	if err != nil {
		http.Error(w, "Error fetching requests", http.StatusInternalServerError)
		return
	}

	fmt.Println("")

	// Load the template
	tmpl, err := template.ParseFS(p.tmpl, "templates/view_project.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render the template with the project and requests data
	data := struct {
		Project  project.Project
		Requests []request.Request
	}{
		Project:  proj,
		Requests: requests,
	}

	tmpl.Execute(w, data)
}
func (p *ProjectHandler) ApplyChanges(w http.ResponseWriter, r *http.Request) {
	// signaling
	p.applyChan <- struct{}{}

	w.Write([]byte("<h1>Applied </h1>"))

}

// func NewProjectHandler() handlers.Hand

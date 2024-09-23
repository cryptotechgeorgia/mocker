package controllers

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	cnv "github.com/cryptotechgeorgia/mocker/foundation/convert"
	"github.com/cryptotechgeorgia/mocker/project"
	"github.com/cryptotechgeorgia/mocker/request"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	bus       *project.Bussiness
	reqBus    *request.Bussiness
	applyChan chan struct{}
	tmpl      embed.FS
}

func NewProjectHandler(bus *project.Bussiness, req *request.Bussiness, applyChan chan struct{}, tmpl embed.FS) ProjectHandler {
	return ProjectHandler{
		bus:       bus,
		reqBus:    req,
		applyChan: applyChan,
		tmpl:      tmpl,
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

func (p *ProjectHandler) RemoveProject(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retreving project id  %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := p.bus.Delete(req.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting project %s", err.Error()), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	select {
	case p.applyChan <- struct{}{}:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "<h1>Applied Successfully</h1>")
	case <-ctx.Done():
		http.Error(w, "Timeout: Failed to apply changes", http.StatusRequestTimeout)
	}
}

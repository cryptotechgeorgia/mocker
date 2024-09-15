package controllers

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/cryptotechgeorgia/mocker/foundation/convert"
	"github.com/cryptotechgeorgia/mocker/payload"
	"github.com/cryptotechgeorgia/mocker/project"
	"github.com/cryptotechgeorgia/mocker/request"
	"github.com/cryptotechgeorgia/mocker/response"
	"github.com/gorilla/mux"
)

var (
	AddRequestErr = errors.New("Failed to add request")
	AddPairErr    = errors.New("Failed to add pair")
	RemovePairErr = errors.New("Failed to remove pari")
)

type RequestHandler struct {
	bus        *request.Bussiness
	respBus    *response.Bussiness
	projBus    *project.Bussiness
	payloadBus *payload.Bussiness
	tmpl       embed.FS
}

func NewRuesthandler(
	bus *request.Bussiness,
	resp *response.Bussiness,
	proj *project.Bussiness,
	payload *payload.Bussiness,
	tmpl embed.FS,
) *RequestHandler {
	return &RequestHandler{
		bus:        bus,
		respBus:    resp,
		projBus:    proj,
		payloadBus: payload,
		tmpl:       tmpl,
	}
}

func (r *RequestHandler) AddRequest(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	path := req.FormValue("path")
	method := req.FormValue("method")

	projectId, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, AddRequestErr.Error(), http.StatusInternalServerError)
		return
	}
	projId, err := strconv.Atoi(projectId)
	if err != nil {
		http.Error(w, AddRequestErr.Error(), http.StatusInternalServerError)
		return
	}

	id, err := r.bus.Add(req.Context(), request.Request{
		ProjectID: projId,
		Path:      path,
		Method:    method,
	})
	if err != nil {
		http.Error(w, AddRequestErr.Error(), http.StatusInternalServerError)
		return
	}

	idstr := strconv.Itoa(id)

	http.Redirect(w, req, fmt.Sprintf("/projects/%d/requests/%s", projId, idstr), http.StatusSeeOther)
}

func (r *RequestHandler) ViewRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	projId, _ := vars["id"]
	reqId, _ := vars["reqId"]

	projIdInt, _ := strconv.Atoi(projId)
	reqIdInt, _ := strconv.Atoi(reqId)

	projInfo, err := r.projBus.Get(req.Context(), projIdInt)
	if err != nil {
		fmt.Printf("errror:  %s\n", err.Error())
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	reqInfo, err := r.bus.Get(req.Context(), reqIdInt)
	if err != nil {
		fmt.Printf("errror:  %s\n", err.Error())
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	reqPayloads, err := r.payloadBus.Filter(req.Context(), payload.FilterBy{
		RequestId: convert.ToIntPtr(reqInfo.ID),
	})

	if err != nil {
		fmt.Printf("errror:  %s\n", err.Error())
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	var pairs []PayloadPair
	for _, p := range reqPayloads {

		resp, err := r.respBus.Filter(req.Context(), response.FilterBy{
			RequestPayloadId: convert.ToIntPtr(p.ID),
		})
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		pair := PayloadPair{
			RequestPath:         reqInfo.Path,
			RequestMethod:       reqInfo.Method,
			RequestPayload:      p.Payload,
			RequestPayloadId:    p.ID,
			ResponsePayload:     resp[0].Payload,
			RequestContentType:  p.ContentType,
			ResponseContentType: resp[0].ContentType,
		}

		pairs = append(pairs, pair)
	}

	tmpl, err := template.ParseFS(r.tmpl, "templates/view_request.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Project project.Project
		Request request.Request
		Pairs   []PayloadPair
	}{
		Project: projInfo,
		Request: reqInfo,
		Pairs:   pairs,
	}

	tmpl.Execute(w, data)
}

type PayloadPair struct {
	RequestPath         string
	RequestMethod       string
	RequestContentType  string
	ResponseContentType string
	RequestPayload      string
	ResponsePayload     string
	RequestPayloadId    int
}

func (r *RequestHandler) AddPair(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	reqId, _ := strconv.Atoi(vars["reqId"])

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	reqContentType := req.FormValue("content_type")
	reqPayload := req.FormValue("payload")
	respContentType := req.FormValue("resp_content_type")
	respPayload := req.FormValue("resp_payload")

	if reqContentType == "application/json" {
		var temp map[string]interface{}
		if err := json.Unmarshal([]byte(reqPayload), &temp); err != nil {
			http.Error(w, "Invalid request JSON", http.StatusBadRequest)
			return
		}
	}

	if respContentType == "application/json" {
		var temp map[string]interface{}
		if err := json.Unmarshal([]byte(respPayload), &temp); err != nil {
			http.Error(w, "Invalid  response JSON", http.StatusBadRequest)
			return
		}
	}

	payloadId, err := r.payloadBus.Add(req.Context(), payload.Payload{
		Payload:     reqPayload,
		RequestId:   reqId,
		ContentType: reqContentType,
	})
	if err != nil {
		http.Error(w, AddPairErr.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.respBus.Add(req.Context(), response.Response{
		Payload:          respPayload,
		ContentType:      respContentType,
		RequestPayloadId: payloadId,
	}); err != nil {
		http.Error(w, AddPairErr.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/projects/%s/requests/%s", vars["id"], vars["reqId"]), http.StatusSeeOther)
}

func (r *RequestHandler) RemoveRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reqId, _ := strconv.Atoi(vars["reqId"])
	if err := r.bus.Delete(req.Context(), reqId); err != nil {
		http.Error(w, RemovePairErr.Error(), http.StatusInternalServerError)
		return

	}

	http.Redirect(w, req, fmt.Sprintf("/projects/%s", vars["id"]), http.StatusSeeOther)
}

func (r *RequestHandler) RemovePair(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reqId, _ := strconv.Atoi(vars["reqId"])
	reqPayloadId, _ := strconv.Atoi(vars["pairId"])

	resp, err := r.respBus.Filter(req.Context(), response.FilterBy{
		RequestPayloadId: convert.ToIntPtr(reqPayloadId),
	})

	if err != nil {
		http.Error(w, RemovePairErr.Error(), http.StatusInternalServerError)
	}
	respInfo := resp[0]

	// remove request_paylod

	if err := r.payloadBus.Delete(req.Context(), reqPayloadId); err != nil {
		http.Error(w, RemovePairErr.Error(), http.StatusInternalServerError)
		return
	}

	// remove response
	if err := r.respBus.Delete(req.Context(), respInfo.ID); err != nil {
		http.Error(w, RemovePairErr.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(reqId)
	http.Redirect(w, req, fmt.Sprintf("/projects/%s/requests/%s", vars["id"], vars["reqId"]), http.StatusSeeOther)

}

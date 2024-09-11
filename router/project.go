package router

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	DefaultResponse ResponseData = ResponseData{
		ContentType: "application/json",
		Payload:     `{"error":"no match"}`,
	}

	DefaultNoRouteResponse ResponseData = ResponseData{
		ContentType: "application/json",
		Payload:     `{"error":"Oops! The route you're looking for does not exist."}`,
	}
)

type RequestData struct {
	ContentType string
	Payload     string
}

type ResponseData struct {
	ContentType string
	Payload     string
}

// gives ability to costumize
// request response pair
type PayloadPair struct {
	req  RequestData
	resp ResponseData
}

func NewRequestData(contentType, payload string) RequestData {
	return RequestData{
		ContentType: contentType,
		Payload:     payload,
	}
}

func NewResponseData(contentType, payload string) ResponseData {
	return ResponseData{
		ContentType: contentType,
		Payload:     payload,
	}
}

func NewPayloadPair(req RequestData, resp ResponseData) PayloadPair {
	return PayloadPair{
		req:  req,
		resp: resp,
	}
}

// func (p PayloadPair) HandlerFunc()

type Request struct {
	method       string
	path         string
	defaultResp  *ResponseData
	payloadPairs []PayloadPair
}

func (r *Request) SetDefaultResponse(data ResponseData) {
	r.defaultResp = &data
}

func NewRequest(method, path string) *Request {
	return &Request{
		method:       method,
		path:         path,
		defaultResp:  &DefaultResponse,
		payloadPairs: []PayloadPair{},
	}
}

func (r *Request) ChangeDefaultResponse(resp ResponseData) {
	r.defaultResp = &resp
}

func (r *Request) AddPayloadPair(pair PayloadPair) {
	r.payloadPairs = append(r.payloadPairs, pair)
}

type Project struct {
	ID       uuid.UUID
	Name     string
	BaseAddr string
	Requests []*Request
}

func (p *Project) AddRequest(req *Request) {
	p.Requests = append(p.Requests, req)
}

func NewProject(name string, addr string) *Project {
	return &Project{
		ID:       uuid.New(),
		Name:     name,
		BaseAddr: addr,
		Requests: []*Request{},
	}
}

func (p *Project) GetRoutes() map[string]*Request {
	router := map[string]*Request{}
	// building  project router
	// route protocol is "method:addr/path"
	for _, r := range p.Requests {
		route := fmt.Sprintf("%s/%s", p.BaseAddr, r.path)
		router[route] = r
	}

	return router
}

package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

// apply  chan activates project
// `*` for  activating all

type MockerRouter struct {
	applyProjects chan struct{}
	populator     *Populator
	handlers      map[string]map[string]http.HandlerFunc
	Mux           *mux.Router
	mu            sync.Mutex
}

func NewMockerRouter(
	applyChan chan struct{},
	populator *Populator,
	mux *mux.Router,

) *MockerRouter {
	return &MockerRouter{
		applyProjects: applyChan,
		populator:     populator,
		handlers:      make(map[string]map[string]http.HandlerFunc),
		Mux:           mux,
	}
}

func (m *MockerRouter) Handle(path string, method string, handler http.HandlerFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.handlers[path] == nil {
		m.handlers[path] = make(map[string]http.HandlerFunc)
	}
	m.handlers[path][method] = handler
}

func (m *MockerRouter) DispatchHandler(path string, method string) http.HandlerFunc {
	m.mu.Lock()
	defer m.mu.Unlock()
	if methods, ok := m.handlers[path]; ok {
		if handler, ok := methods[method]; ok {
			return handler
		}
	}
	return nil
}

func (m *MockerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// try to serve  dynamic route
	handler := m.DispatchHandler(r.URL.Path, r.Method)
	if handler != nil {
		fmt.Println("shemovida")
		handler(w, r)
		return
	}

	// back to static routes
	m.Mux.ServeHTTP(w, r)
}

func (m *MockerRouter) fixPath(path string) string {
	if path[:1] != "/" {
		return fmt.Sprintf("/%s", path)
	}
	return path
}

func (m *MockerRouter) ValidateJsonPair(headers http.Header, body []byte, pair PayloadPair) bool {
	reqJson, err := normalizeJson(body)
	if err != nil {
		return false
	}

	schema := gojsonschema.NewStringLoader(pair.req.Payload)
	reqBody := gojsonschema.NewStringLoader(reqJson)

	_, err = gojsonschema.Validate(schema, reqBody)
	if err != nil {
		return false
	}

	return true
}

func (m *MockerRouter) BuildHandler(request Request) func(resp http.ResponseWriter, req *http.Request) {
	handler := func(resp http.ResponseWriter, req *http.Request) {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprint(resp, err)
			return
		}

		for _, pair := range request.payloadPairs {
			fmt.Printf("validating pair  %+v\n ", pair)
			if req.Header.Get("Content-Type") != "application/json" && string(reqBody) == pair.req.Payload {

				resp.Header().Set("Content-Type", pair.resp.ContentType)
				fmt.Fprint(resp, pair.resp.Payload)
				return
			}

			valid := m.ValidateJsonPair(req.Header, reqBody, pair)
			if valid {
				fmt.Printf("pair %+v is  validdd\n", pair)
				resp.Header().Set("Content-Type", pair.resp.ContentType)
				fmt.Fprint(resp, pair.resp.Payload)
				return
			}

		}

		// nothing  matched
		resp.Header().Set("Content-Type", DefaultResponse.ContentType)
		fmt.Fprint(resp, DefaultResponse.Payload)
	}

	return handler
}

func (m *MockerRouter) applyPrefix(prefix string, path string) string {
	return fmt.Sprintf("/%s%s", prefix, path)
}

func (m *MockerRouter) Listen(ctx context.Context, errChan chan error) {
	for {
		select {
		case <-m.applyProjects:
			fmt.Println("shemovida ")

			projects, err := m.populator.Populate(ctx)
			if err != nil {
				errChan <- err
				return
			}

			for _, project := range projects {
				// sets  project name for mux router name
				for _, request := range project.Requests {
					// build one handlerFunc  per  request
					path := m.applyPrefix(project.Name, m.fixPath(request.path))
					handler := m.BuildHandler(request)

					m.Handle(path, request.method, handler)
					// m.Mux.HandleFunc(path, handler).Methods(request.method)
				}
			}

		case <-ctx.Done():
			fmt.Println("shutting down mocker router")
			return
		}
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", DefaultNoRouteResponse.ContentType)
	fmt.Fprintln(w, DefaultNoRouteResponse.Payload)
}

func normalizeJson(inp []byte) (string, error) {
	var marshalledReq interface{}
	if err := json.Unmarshal(inp, &marshalledReq); err != nil {
		return "", err
	}
	normalizedJson, err := json.Marshal(marshalledReq)
	if err != nil {
		return "", err
	}
	return string(normalizedJson), nil
}

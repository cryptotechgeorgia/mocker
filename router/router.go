package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cryptotechgeorgia/mocker/foundation/convert"
	"github.com/cryptotechgeorgia/mocker/payload"
	"github.com/cryptotechgeorgia/mocker/project"
	"github.com/cryptotechgeorgia/mocker/request"
	"github.com/cryptotechgeorgia/mocker/response"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

// apply  chan activates project
// `*` for  activating all

type MockerRouter struct {
	projBus       *project.Bussiness
	payloadBus    *payload.Bussiness
	respBus       *response.Bussiness
	reqBus        *request.Bussiness
	projects      []Project
	applyProjects chan struct{}
	router        *mux.Router
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

func NewMockerRouter(
	applyChan chan struct{},
	projbus *project.Bussiness,
	reqbus *request.Bussiness,
	respbus *response.Bussiness,
	paybus *payload.Bussiness,
	router *mux.Router,
) MockerRouter {
	return MockerRouter{
		applyProjects: applyChan,
		projBus:       projbus,
		reqBus:        reqbus,
		respBus:       respbus,
		payloadBus:    paybus,
		router:        router,
	}
}

func (m *MockerRouter) Listen(ctx context.Context, errChan chan error, doneApply chan struct{}) {
	for {
		<-m.applyProjects

		projects, err := m.projBus.All(ctx)
		if err != nil {
			errChan <- err
			return
		}

		// initialize projects
		activateProjects := []*Project{}
		for _, project := range projects {

			proj := NewProject(project.Name, project.BaseAddr)
			reqs, err := m.reqBus.Filter(ctx, request.FilterBy{
				ProjectId: convert.ToIntPtr(project.ID),
			})
			if err != nil {
				errChan <- err
				return
			}

			// adding requests to project
			for _, r := range reqs {
				req := NewRequest(r.Method, r.Path)

				// getting payloads
				payloads, err := m.payloadBus.Filter(ctx, payload.FilterBy{
					RequestId: convert.ToIntPtr(r.ID),
				})
				if err != nil {
					errChan <- err
					continue
				}

				// getting corresponting responses for this payloads
				for _, reqPayload := range payloads {
					respData, err := m.respBus.Filter(ctx, response.FilterBy{
						RequestPayloadId: convert.ToIntPtr(reqPayload.ID),
					})
					if err != nil {
						errChan <- err
						continue
					}
					req.AddPayloadPair(
						PayloadPair{
							req: RequestData{
								ContentType: reqPayload.ContentType,
								Payload:     reqPayload.Payload,
							},
							resp: ResponseData{
								ContentType: respData[0].ContentType,
								Payload:     respData[0].Payload,
							},
						})
				}
				proj.AddRequest(req)
			}
			activateProjects = append(activateProjects, proj)

			fmt.Printf("project %s : len(requests)=%d\n", proj.Name, len(proj.Requests))
		}

		for _, project := range activateProjects {
			projectRouter := project.Requests
			fmt.Println("project requests are ", projectRouter)

			// sets  project name for mux router name
			for _, projectReq := range projectRouter {

				// build one handlerFunc  per  request
				path := projectReq.path
				if projectReq.path[:1] != "/" {
					path = fmt.Sprintf("/%s", projectReq.path)
				}

				// project name as prefix
				path = fmt.Sprintf("/%s%s", project.Name, path)

				// building handler which  trys to match incoming request body
				// to  one of the request of the project, if  match does not exists
				// return  default error
				// if payload exists but the content-type  did not match
				// return wrapped default error

				handler := func(resp http.ResponseWriter, req *http.Request) {

					reqBody, err := io.ReadAll(req.Body)
					if err != nil {
						fmt.Fprint(resp, err)
						return
					}

					reqJson, err := normalizeJson(reqBody)
					if err != nil {
						fmt.Fprint(resp, err)
						return
					}

					for _, pair := range projectReq.payloadPairs {

						schemaLoader := gojsonschema.NewStringLoader(pair.req.Payload)
						requestLoader := gojsonschema.NewStringLoader(reqJson)

						res, err := gojsonschema.Validate(schemaLoader, requestLoader)
						if err != nil {
							log.Println("error while validating ", err.Error())
							continue
						}

						if err == nil {
							if res.Valid() {
								fmt.Println("valid ")
								// setting response content-type
								resp.Header().Set("Content-Type", pair.resp.ContentType)

								// check incoming request content-type  on  project inside  one
								if req.Header.Get("Content-Type") != pair.req.ContentType {
									// setting request specific default  contentType
									resp.Header().Set("Content-Type", projectReq.defaultResp.ContentType)
									fmt.Println("content type does not matched")
									fmt.Fprint(resp, projectReq.defaultResp.Payload)
									return
								}

								fmt.Fprint(resp, pair.resp.Payload)
								return
							}

						}

					}

					// nothing  matched
					resp.Header().Set("Content-Type", projectReq.defaultResp.ContentType)
					fmt.Fprint(resp, projectReq.defaultResp.Payload)
					return
				}

				m.router.Name(project.Name).Path(path).HandlerFunc(handler).Methods(projectReq.method)
			}
		}
	}

}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", DefaultNoRouteResponse.ContentType)
	fmt.Fprintln(w, DefaultNoRouteResponse.Payload)
}

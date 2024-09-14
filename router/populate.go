package router

import (
	"context"

	"github.com/cryptotechgeorgia/mocker/payload"
	"github.com/cryptotechgeorgia/mocker/project"
	"github.com/cryptotechgeorgia/mocker/request"
	"github.com/cryptotechgeorgia/mocker/response"
)

type Populator struct {
	projBus    *project.Bussiness
	payloadBus *payload.Bussiness
	respBus    *response.Bussiness
	reqBus     *request.Bussiness
}

func NewPopulator(
	projBus *project.Bussiness,
	payloadBus *payload.Bussiness,
	respBus *response.Bussiness,
	reqBus *request.Bussiness) *Populator {

	return &Populator{
		projBus:    projBus,
		payloadBus: payloadBus,
		respBus:    respBus,
		reqBus:     reqBus,
	}
}

func (p *Populator) payloadPairs(ctx context.Context, reqId int) ([]PayloadPair, error) {
	payloads, err := p.payloadBus.Filter(ctx, payload.FilterBy{RequestId: &reqId})
	if err != nil {
		return []PayloadPair{}, err
	}

	var pairs []PayloadPair
	for _, pl := range payloads {
		resps, err := p.respBus.Filter(ctx, response.FilterBy{RequestPayloadId: &pl.ID})
		if err != nil {
			return []PayloadPair{}, err
		}

		pair := NewPayloadPair(
			NewRequestData(pl.ContentType, pl.Payload),
			NewResponseData(resps[0].ContentType, resps[0].Payload),
		)

		pairs = append(pairs, pair)
	}

	return pairs, nil

}

func (p *Populator) requests(ctx context.Context, projId int) ([]Request, error) {
	reqs, err := p.reqBus.Filter(ctx, request.FilterBy{ProjectId: &projId})
	if err != nil {
		return []Request{}, err
	}

	var requests []Request

	for _, req := range reqs {
		pairs, err := p.payloadPairs(ctx, req.ID)
		if err != nil {
			return []Request{}, err
		}

		projReq := Request{
			method:       req.Method,
			path:         req.Path,
			defaultResp:  &DefaultResponse,
			payloadPairs: pairs,
		}

		requests = append(requests, projReq)
	}
	return requests, nil
}

func (p *Populator) Populate(ctx context.Context) ([]Project, error) {
	projs, err := p.projBus.All(ctx)
	if err != nil {
		return []Project{}, err
	}

	var projects []Project

	for _, proj := range projs {

		projReqs, err := p.requests(ctx, proj.ID)
		if err != nil {
			return []Project{}, err
		}
		project := Project{
			Name:     proj.Name,
			BaseAddr: proj.BaseAddr,
			Requests: projReqs,
		}

		projects = append(projects, project)
	}

	return projects, err
}

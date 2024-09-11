package store

import (
	"github.com/cryptotechgeorgia/mocker/request"
)

type dbRequest struct {
	ID        int    `db:"id"`
	ProjectID int    `db:"project_id"`
	Path      string `db:"path"`
	Method    string `db:"method"`
}

func fromRequest(a request.Request) dbRequest {
	return dbRequest{
		ID:        a.ID,
		ProjectID: a.ProjectID,
		Path:      a.Path,
		Method:    a.Method,
	}
}

func toRequest(a dbRequest) request.Request {
	return request.Request{
		ID:        a.ID,
		ProjectID: a.ProjectID,
		Path:      a.Path,
		Method:    a.Method,
	}
}

func toRequestSlice(dbReqs []dbRequest) []request.Request {
	var reqs []request.Request

	for _, p := range dbReqs {
		reqs = append(reqs, toRequest(p))
	}

	return reqs
}

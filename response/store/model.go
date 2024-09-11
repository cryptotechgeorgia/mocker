package store

import (
	"github.com/cryptotechgeorgia/mocker/response"
)

// id INTEGER PRIMARY KEY AUTOINCREMENT,
// payload TEXT,
// request_payload_id INTEGER,
// content_type TEXT

type dbResponse struct {
	ID               int    `db:"id"`
	Payload          string `db:"payload"`
	RequestPayloadId int    `db:"request_payload_id"`
	ContentType      string `db:"content_type"`
}

func fromResponse(a response.Response) dbResponse {
	return dbResponse{
		ID:               a.ID,
		Payload:          a.Payload,
		RequestPayloadId: a.RequestPayloadId,
		ContentType:      a.ContentType,
	}
}

func toResponse(a dbResponse) response.Response {
	return response.Response{
		ID:               a.ID,
		Payload:          a.Payload,
		RequestPayloadId: a.RequestPayloadId,
		ContentType:      a.ContentType,
	}
}

func toResponseSlice(dbResps []dbResponse) []response.Response {
	var resps []response.Response

	for _, resp := range dbResps {
		resps = append(resps, toResponse(resp))
	}

	return resps
}

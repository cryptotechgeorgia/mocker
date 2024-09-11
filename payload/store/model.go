package store

import (
	"github.com/cryptotechgeorgia/mocker/payload"
)

// CREATE TABLE IF NOT EXISTS request_payload (
//
//	id INTEGER PRIMARY KEY AUTOINCREMENT,
//	payload TEXT,
//	request_id INTEGER,
//	content_type TEXT
//
// );
type dbPayload struct {
	ID          int    `db:"id"`
	Payload     string `db:"payload"`
	RequestId   int    `db:"request_id"`
	ContentType string `db:"content_type"`
}

func fromPayload(a payload.Payload) dbPayload {
	return dbPayload{
		ID:          a.ID,
		Payload:     a.Payload,
		RequestId:   a.RequestId,
		ContentType: a.ContentType,
	}
}

func toPayload(a dbPayload) payload.Payload {
	return payload.Payload{
		ID:          a.ID,
		Payload:     a.Payload,
		RequestId:   a.RequestId,
		ContentType: a.ContentType,
	}
}

func toPayloadSlice(dbPs []dbPayload) []payload.Payload {
	var payloads []payload.Payload
	for _, p := range dbPs {
		payloads = append(payloads, toPayload(p))
	}

	return payloads
}

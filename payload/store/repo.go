package store

import (
	"bytes"
	"context"
	"log"

	"github.com/cryptotechgeorgia/mocker/payload"
	"github.com/jmoiron/sqlx"
)

// CREATE TABLE IF NOT EXISTS request_payload (
// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
// 	payload TEXT,
// 	request_id INTEGER,
// 	content_type TEXT
// );

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Get(ctx context.Context, id int) (payload.Payload, error) {
	var dbPayload dbPayload

	err := r.db.GetContext(ctx, &dbPayload, `SELECT * FROM request_payload WHERE id=$1`, id)
	if err != nil {
		log.Println("Error fetching projects:", err)
		return payload.Payload{}, err
	}
	return toPayload(dbPayload), nil
}

func (r *Repo) Add(ctx context.Context, p payload.Payload) (int, error) {
	var id int
	dbPayload := fromPayload(p)

	err := r.db.GetContext(ctx, &id, `INSERT INTO request_payload (payload, request_id, content_type) VALUES ($1, $2, $3) returning id`,
		dbPayload.Payload,
		dbPayload.RequestId,
		dbPayload.ContentType,
	)
	if err != nil {
		log.Println("Error adding project:", err)
		return 0, err
	}

	return id, nil
}

func (r *Repo) All(ctx context.Context) ([]payload.Payload, error) {
	var dbPayloads []dbPayload
	err := r.db.SelectContext(ctx, &dbPayloads, `SELECT * FROM request_payload`)
	if err != nil {
		log.Println("Error fetching projects:", err)
		return nil, err
	}

	return toPayloadSlice(dbPayloads), nil
}

func (r *Repo) Filter(ctx context.Context, filter payload.FilterBy) ([]payload.Payload, error) {
	var dbPayloads []dbPayload
	q := `SELECT * FROM request_payload `
	buf := bytes.NewBufferString(q)
	applyFilter(filter, buf)
	if err := r.db.SelectContext(ctx, &dbPayloads, buf.String()); err != nil {
		return []payload.Payload{}, err
	}

	return toPayloadSlice(dbPayloads), nil
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM request_payload WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting project:", err)
		return err
	}

	return nil
}

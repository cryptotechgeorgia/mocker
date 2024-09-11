package store

import (
	"bytes"
	"context"
	"log"

	"github.com/cryptotechgeorgia/mocker/response"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return Repo{
		db: db,
	}
}

func (r Repo) Add(ctx context.Context, p response.Response) error {
	dbResponse := fromResponse(p)
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO response (payload, request_payload_id,  content_type) VALUES (:payload, :request_payload_id, :content_type)`, dbResponse)
	if err != nil {
		log.Println("Error adding response:", err)
		return err
	}
	return nil
}

func (r Repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM response WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting project:", err)
		return err
	}
	return nil
}

func (r Repo) Filter(ctx context.Context, filter response.FilterBy) ([]response.Response, error) {
	var dbResps []dbResponse
	q := `SELECT * FROM response `
	buf := bytes.NewBufferString(q)

	applyFilter(filter, buf)

	if err := r.db.SelectContext(ctx, &dbResps, buf.String()); err != nil {
		return []response.Response{}, nil
	}

	return toResponseSlice(dbResps), nil

}

func (r Repo) Get(ctx context.Context, id int) (response.Response, error) {
	var dbResponse dbResponse
	if err := r.db.GetContext(ctx, &dbResponse, "SELECT * FROM response WHERE id = $1", id); err != nil {
		return response.Response{}, err
	}
	return toResponse(dbResponse), nil
}

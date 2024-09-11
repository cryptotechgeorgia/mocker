package store

import (
	"bytes"
	"context"
	"log"

	"github.com/cryptotechgeorgia/mocker/request"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Add(ctx context.Context, req request.Request) (int, error) {
	var id int
	dbReq := fromRequest(req)

	if err := r.db.GetContext(ctx, &id, `INSERT INTO request (project_id, path, method) 
VALUES ($1, $2, $3) returning id`, dbReq.ProjectID, dbReq.Path, dbReq.Method); err != nil {
		log.Println("Error adding request:", err)
		return 0, err
	}
	return id, nil
}

func (r *Repo) All(ctx context.Context) ([]request.Request, error) {
	var dbReqs []dbRequest

	err := r.db.SelectContext(ctx, &dbReqs, `SELECT * FROM request`)
	if err != nil {
		log.Println("Error fetching projects:", err)
		return nil, err
	}
	return toRequestSlice(dbReqs), nil
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM request WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting project:", err)
		return err
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, id int) (request.Request, error) {
	var dbReq dbRequest
	if err := r.db.GetContext(ctx, &dbReq, `SELECT * FROM request WHERE id = $1`, id); err != nil {
		return request.Request{}, err
	}

	return toRequest(dbReq), nil
}

func (r *Repo) Filter(ctx context.Context, filter request.FilterBy) ([]request.Request, error) {
	var dbReqs []dbRequest
	q := `SELECT * FROM request `
	buf := bytes.NewBufferString(q)

	applyFilter(filter, buf)

	if err := r.db.SelectContext(ctx, &dbReqs, buf.String()); err != nil {
		return []request.Request{}, nil
	}

	return toRequestSlice(dbReqs), nil
}

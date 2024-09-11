package store

import (
	"context"
	"log"

	"github.com/cryptotechgeorgia/mocker/project"
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

func (r *Repo) Get(ctx context.Context, id int) (project.Project, error) {
	var dbProject dbProject

	err := r.db.GetContext(ctx, &dbProject, `SELECT * FROM project WHERE id = $1`, id)
	if err != nil {
		log.Println("Error fetching projects:", err)
		return project.Project{}, err
	}
	return toProject(dbProject), nil
}

func (r *Repo) Add(ctx context.Context, p project.Project) error {
	dbProject := fromProject(p)

	_, err := r.db.NamedExecContext(ctx, `INSERT INTO project (name, base_addr) VALUES (:name, :base_addr)`, dbProject)
	if err != nil {
		log.Println("Error adding project:", err)
		return err
	}

	return nil
}

func (r *Repo) All(ctx context.Context) ([]project.Project, error) {
	var dbProjects []dbProject
	err := r.db.SelectContext(ctx, &dbProjects, `SELECT * FROM project`)
	if err != nil {
		log.Println("Error fetching projects:", err)
		return nil, err
	}

	return toProjectSlice(dbProjects), nil
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM project WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting project:", err)
		return err
	}

	return nil
}

package store

import "github.com/cryptotechgeorgia/mocker/project"

type dbProject struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	BaseAddr string `db:"base_addr"`
}

func fromProject(a project.Project) dbProject {
	return dbProject{
		ID:       a.ID,
		Name:     a.Name,
		BaseAddr: a.BaseAddr,
	}
}

func toProject(a dbProject) project.Project {
	return project.Project{
		ID:       a.ID,
		Name:     a.Name,
		BaseAddr: a.BaseAddr,
	}
}

func toProjectSlice(dbPs []dbProject) []project.Project {
	var projs []project.Project

	for _, p := range dbPs {
		projs = append(projs, toProject(p))
	}

	return projs
}

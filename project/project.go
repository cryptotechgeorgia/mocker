package project

import "context"

type Project struct {
	ID       int
	Name     string
	BaseAddr string
}

func (p Project) GetName() string {
	return p.Name
}

type Storer interface {
	Add(context.Context, Project) error
	Delete(context.Context, int) error
	All(context.Context) ([]Project, error)
	Get(context.Context, int) (Project, error)
}

type Bussiness struct {
	store Storer
}

func NewBusiness(repo Storer) *Bussiness {
	return &Bussiness{
		store: repo,
	}
}

func (b *Bussiness) Get(ctx context.Context, id int) (Project, error) {
	return b.store.Get(ctx, id)
}

func (b *Bussiness) Add(ctx context.Context, prj Project) error {
	return b.store.Add(ctx, prj)
}

func (b *Bussiness) Delete(ctx context.Context, id int) error {
	return b.store.Delete(ctx, id)
}

func (b *Bussiness) All(ctx context.Context) ([]Project, error) {
	return b.store.All(ctx)
}

// func (b *Bussiness) Get(ctx context.Context, id int) (Project, error) {
// 	resp, err := b.store.Get(ctx, id)
// 	if err != nil {
// 		return Request{}, err
// 	}

// 	return resp, nil
// }

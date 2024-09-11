package request

import (
	"context"
)

type Request struct {
	ID        int
	ProjectID int
	Path      string
	Method    string
}

type Storer interface {
	Add(context.Context, Request) (int, error)
	Delete(context.Context, int) error
	Get(context.Context, int) (Request, error)
	All(context.Context) ([]Request, error)
	Filter(context.Context, FilterBy) ([]Request, error)
}

type Bussiness struct {
	store Storer
}

func NewBusiness(repo Storer) *Bussiness {
	return &Bussiness{
		store: repo,
	}
}

func (b *Bussiness) Add(ctx context.Context, req Request) (int, error) {
	return b.store.Add(ctx, req)
}

func (b *Bussiness) Get(ctx context.Context, id int) (Request, error) {
	resp, err := b.store.Get(ctx, id)
	if err != nil {
		return Request{}, err
	}

	return resp, nil
}

func (b *Bussiness) Delete(ctx context.Context, id int) error {
	return b.store.Delete(ctx, id)
}

func (b *Bussiness) All(ctx context.Context) ([]Request, error) {

	return b.store.All(ctx)
}

func (b *Bussiness) Filter(ctx context.Context, fBy FilterBy) ([]Request, error) {
	return b.store.Filter(ctx, fBy)
}

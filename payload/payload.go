package payload

import "context"

type Payload struct {
	ID          int
	Payload     string
	RequestId   int
	ContentType string
}

type Storer interface {
	Add(context.Context, Payload) (int, error)
	Delete(context.Context, int) error
	Get(context.Context, int) (Payload, error)
	Filter(context.Context, FilterBy) ([]Payload, error)
	// All(context.Context)
}

type Bussiness struct {
	store Storer
}

func NewBusiness(repo Storer) *Bussiness {
	return &Bussiness{
		store: repo,
	}
}

func (b *Bussiness) Get(ctx context.Context, id int) (Payload, error) {
	return b.store.Get(ctx, id)
}

func (b *Bussiness) Add(ctx context.Context, p Payload) (int, error) {
	return b.store.Add(ctx, p)
}

func (b *Bussiness) Delete(ctx context.Context, id int) error {
	return b.store.Delete(ctx, id)
}

func (b *Bussiness) Filter(ctx context.Context, filter FilterBy) ([]Payload, error) {
	return b.store.Filter(ctx, filter)
}

// func (b *Bussiness) All(ctx context.Context) ([]Payload, error) {
// 	return b.store.All(ctx)
// }

// func (b *Bussiness) Get(ctx context.Context, id int) (Project, error) {
// 	resp, err := b.store.Get(ctx, id)
// 	if err != nil {
// 		return Request{}, err
// 	}

// 	return resp, nil
// }

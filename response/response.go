package response

import (
	"context"
)

// id INTEGER PRIMARY KEY AUTOINCREMENT,
// payload TEXT,
// request_payload_id INTEGER,
// content_type TEXT

type Response struct {
	ID               int
	Payload          string
	ContentType      string
	RequestPayloadId int
}

type Storer interface {
	Add(context.Context, Response) error
	Delete(context.Context, int) error
	Get(context.Context, int) (Response, error)
	Filter(context.Context, FilterBy) ([]Response, error)
}

type Bussiness struct {
	store Storer
}

func NewBusiness(repo Storer) *Bussiness {
	return &Bussiness{
		store: repo,
	}
}

func (b *Bussiness) Get(ctx context.Context, id int) (Response, error) {

	resp, err := b.store.Get(ctx, id)
	if err != nil {
		return Response{}, err
	}
	return resp, nil
}

func (b *Bussiness) Delete(ctx context.Context, id int) error {
	return b.store.Delete(ctx, id)
}

func (b *Bussiness) Add(ctx context.Context, resp Response) error {

	return b.store.Add(ctx, resp)

}
func (b *Bussiness) Filter(ctx context.Context, filter FilterBy) ([]Response, error) {
	return b.store.Filter(ctx, filter)
}

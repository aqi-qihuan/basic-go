package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type AnyESADO struct {
	client *elastic.Client
}

func NewAnyESDAO(client *elastic.Client) AnyDAO {
	return &AnyESADO{client: client}
}

func (a *AnyESADO) Input(ctx context.Context, index, docId, data string) error {
	_, err := a.client.Index().
		Index(index).Id(docId).BodyString(data).Do(ctx)
	return err
}

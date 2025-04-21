package redispresence

import (
	"GameApp/pkg/richerror"
	"context"
	"time"
)

func (d DB) Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error {
	const OP = "redispresences.Upsert"

	_, err := d.adaptor.Client().Set(ctx, key, timestamp, expTime).Result()
	if err != nil {
		return richerror.New(OP).WithKind(richerror.KindUnexpected).WithWrappedError(err)
	}
	return nil
}

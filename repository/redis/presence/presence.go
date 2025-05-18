package redispresence

import (
	"GameApp/param"
	"GameApp/pkg/richerror"
	"context"
	"fmt"
	"strconv"
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

func (d DB) GetPresence(ctx context.Context, userIDs []uint) ([]param.GetPresenceItem, error) {
	const OP = "redispresences.Upsert"
	keys := make([]string, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = fmt.Sprintf("presence:%d", userID) // Assuming presence data is stored with this key format
	}

	// Fetch all presence data in one request
	results, err := d.adaptor.Client().MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get presence data: %w", err)
	}

	var presenceItems []param.GetPresenceItem
	for i, res := range results {
		if res == nil {
			continue // Skip missing entries
		}

		timestamp, err := strconv.ParseInt(res.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp for user %d: %w", userIDs[i], err)
		}

		presenceItems = append(presenceItems, param.GetPresenceItem{
			UserID:    userIDs[i],
			Timestamp: timestamp,
		})
	}

	return presenceItems, nil
}

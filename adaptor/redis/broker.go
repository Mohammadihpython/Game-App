package redis

import (
	"GameApp/entity"
	"context"
	"github.com/labstack/gommon/log"
	"time"
)

func (a Adaptor) Publish(event entity.Event, payload string) {
	// TODO - add 1 to config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	if err := a.client.Publish(ctx, string(event), payload).Err(); err != nil {
		log.Errorf("publish err: %v\n", err)
		// TODO - log
		// TODO - update metrics
	}

	// TODO - update metrics
}

func (a Adaptor) Consumer(event entity.Event) string {
	// TODO - add 1 to config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	pubSub := a.client.Subscribe(ctx, string(event))
	defer pubSub.Close()
	ch := pubSub.Channel()
	for msg := range ch {
		return msg.Payload
		// Process event: Initialize game session, notify players, etc.

	}
	return ""

	// TODO - update metrics
}

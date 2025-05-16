package broker

import "GameApp/entity"

type Publisher interface {
	Publish(event entity.Event, payload string)
}
type Consumer interface {
	Consumer(event entity.Event) string
}

package broker

import "GameApp/entity"

type Publisher interface {
	Publish(event entity.Event, payload string)
}

package param

import (
	"GameApp/entity"
	"time"
)

type AddToWaitingListRequest struct {
	UserID   uint            `json:"user_id"`
	Category entity.Category `json:"category"`
}
type AddToWaitingResponse struct {
	Timeout time.Duration `json:"timeout"`
}

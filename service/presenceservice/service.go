package presenceservice

import (
	"GameApp/param"
	"GameApp/pkg/richerror"
	"context"
	"fmt"
	"time"
)

type Config struct {
	PresenceExpireTime time.Duration `koanf:"presence_expire_time"`
	PresencePrefix     string        `koanf:"presence_prefix"`
}

type Repo interface {
	Upsert(ctx context.Context, key string, timestamp int64, PresenceExpireTime time.Duration) error
}

type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}

func (s Service) Upsert(ctx context.Context, req param.UpsertPresenceRequest) (param.UpsertPresenceResponse, error) {
	const Op = "presenceservice.Upsert"
	err := s.repo.Upsert(ctx, fmt.Sprintf("%s:%d", s.config.PresencePrefix, req.UserID), req.Timestamp, s.config.PresenceExpireTime)
	if err != nil {
		return param.UpsertPresenceResponse{}, richerror.New(Op).WithWrappedError(err)
	}
	return param.UpsertPresenceResponse{}, nil

}

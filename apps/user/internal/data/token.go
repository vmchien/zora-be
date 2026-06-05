package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"vn.vato.zora.be.api/apps/user/internal/biz"
)

type TokenRepo struct {
	data *Data
	log  *log.Helper
}

func (t TokenRepo) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error {
	// TODO implement me
	return nil
}

func (t TokenRepo) GetRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	// TODO implement me
	panic("implement me")
}

func (t TokenRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	// TODO implement me
	panic("implement me")
}

func (t TokenRepo) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	// TODO implement me
	panic("implement me")
}

func (t TokenRepo) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func NewTokenRepo(data *Data, logger log.Logger) biz.TokenRepo {
	return &TokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

package account

import (
	"context"
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

const (
	accountKey = "account"
)

type Service struct {
	cache cache.Cache[account.Account]
	log   logger.Logger
}

func New(cache cache.Cache[account.Account], l logger.Logger) *Service {
	return &Service{
		cache: cache,
		log:   l,
	}
}

func (s *Service) buildLogger(ctx context.Context) logger.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)

	return s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)
}

func (s *Service) Get(ctx context.Context) (account.Account, error) {
	log := s.buildLogger(ctx)

	log.Debug("retrieving account from cache")
	acct, err := s.cache.Get(ctx, accountKey)
	if err != nil {
		if !errors.Is(err, bolt.ErrKeyNotFound) {
			log.Warn("failed to retrieve cached account",
				logger.Error(err),
			)
			return acct, err
		}
		log.Info("no account cached")
	} else {
		log.Info("account retrieved",
			logger.String("username", acct.Username),
		)
	}

	return acct, nil
}

func (s *Service) Put(ctx context.Context, acct account.Account) error {
	log := s.buildLogger(ctx)

	log.Debug("caching account",
		logger.String("username", acct.Username),
	)
	if err := s.cache.Set(ctx, accountKey, acct); err != nil {
		log.Warn("failed to cache account",
			logger.String("username", acct.Username),
			logger.Error(err),
		)
		return err
	}
	log.Info("account cached",
		logger.String("username", acct.Username),
	)

	return nil
}

func (s *Service) Delete(ctx context.Context) error {
	log := s.buildLogger(ctx)

	log.Debug("deleting cached account")
	if err := s.cache.Delete(ctx, accountKey); err != nil {
		log.Warn("failed to delete cached account",
			logger.Error(err),
		)
		return err
	}
	log.Info("deleted cached account")

	return nil
}

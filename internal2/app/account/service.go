package account

import (
	"context"
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

const (
	accountKey = "account"
)

type Service struct {
	cache  cache.Cache[account.Account]
	logger logging.Logger
}

func New(cache cache.Cache[account.Account], logger logging.Logger) *Service {
	return &Service{
		cache:  cache,
		logger: logger,
	}
}

func (s *Service) buildLogger(ctx context.Context) logging.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)

	return s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)
}

func (s *Service) Get(ctx context.Context) (account.Account, error) {
	logger := s.buildLogger(ctx)

	logger.Debug("retrieving account from cache")
	acct, err := s.cache.Get(ctx, accountKey)
	if err != nil {
		if !errors.Is(err, core.ErrKeyNotFound) {
			logger.Warn("failed to retrieve cached account",
				logging.Error(err),
			)
			return acct, err
		}
		logger.Info("no account cached")
	} else {
		logger.Info("account retrieved",
			logging.String("username", acct.Username),
		)
	}

	return acct, nil
}

func (s *Service) Put(ctx context.Context, acct account.Account) error {
	logger := s.buildLogger(ctx)

	logger.Debug("caching account",
		logging.String("username", acct.Username),
	)
	if err := s.cache.Set(ctx, accountKey, acct); err != nil {
		logger.Warn("failed to cache account",
			logging.String("username", acct.Username),
			logging.Error(err),
		)
		return err
	}
	logger.Info("account cached",
		logging.String("username", acct.Username),
	)

	return nil
}

func (s *Service) Delete(ctx context.Context) error {
	logger := s.buildLogger(ctx)

	logger.Debug("deleting cached account")
	if err := s.cache.Delete(ctx, accountKey); err != nil {
		logger.Warn("failed to delete cached account",
			logging.Error(err),
		)
		return err
	}
	logger.Info("deleted cached account")

	return nil
}

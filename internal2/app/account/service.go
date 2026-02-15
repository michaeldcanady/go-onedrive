package account

import (
	"context"
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

const (
	accountKey = "account"
)

var accountKeySerializer = func() ([]byte, error) { return json.Marshal(accountKey) }

type Service struct {
	cache  *abstractions.Cache2
	logger logging.Logger
}

func New(cache *abstractions.Cache2, logger logging.Logger) *Service {
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
	var acct account.Account

	logger := s.buildLogger(ctx)

	logger.Debug("retrieving account from cache")
	if err := s.cache.Get(ctx, accountKeySerializer, func(b []byte) error { return json.Unmarshal(b, &acct) }); err != nil {
		logger.Warn("failed to retrieve cached account",
			logging.Error(err),
		)
		return acct, err
	}
	logger.Info("account retrieved",
		logging.String("username", acct.Username),
	)

	return acct, nil
}

func (s *Service) Put(ctx context.Context, acct account.Account) error {
	logger := s.buildLogger(ctx)

	logger.Debug("caching account",
		logging.String("username", acct.Username),
	)
	if err := s.cache.Set(ctx, accountKeySerializer, func() ([]byte, error) { return json.Marshal(acct) }); err != nil {
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
	if err := s.cache.Delete(ctx, accountKeySerializer); err != nil {
		logger.Warn("failed to delete cached account",
			logging.Error(err),
		)
		return err
	}
	logger.Info("deleted cached account")

	return nil
}

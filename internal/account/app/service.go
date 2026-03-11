package app

import (
	"context"
	"errors"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

const (
	accountKey = "account"
)

type Service struct {
	cache pkgcache.Cache[domainaccount.Account]
	log   logger.Logger
}

func New(cache pkgcache.Cache[domainaccount.Account], l logger.Logger) *Service {
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

func (s *Service) Get(ctx context.Context) (domainaccount.Account, error) {
	log := s.buildLogger(ctx)

	log.Debug("retrieving account from cache")
	acct, err := s.cache.Get(ctx, accountKey)
	if err != nil {
		if !errors.Is(err, cache.ErrKeyNotFound) {
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

func (s *Service) Put(ctx context.Context, acct domainaccount.Account) error {
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

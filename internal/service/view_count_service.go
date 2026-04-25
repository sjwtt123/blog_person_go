package service

import (
	"blog/internal/repository"
	"blog/pkg/logger"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type viewCountService struct {
	articleRepo   repository.ArticleRepository
	viewCountRepo repository.ViewCountRepository
}

func NewViewCountService(articleRepo repository.ArticleRepository, viewCountRepo repository.ViewCountRepository) ViewCountService {
	return &viewCountService{
		articleRepo:   articleRepo,
		viewCountRepo: viewCountRepo,
	}
}

func (s *viewCountService) IncrementViewCount(articleID uint) error {
	return s.viewCountRepo.Increment(articleID)
}

func (s *viewCountService) SyncViewCounts() error {
	all, err := s.viewCountRepo.GetAll()
	if err != nil {
		logger.Error("Redis HGetAll failed", zap.Error(err))
		return err
	}

	if len(all) == 0 {
		return nil
	}

	for field, valStr := range all {
		articleID, err := strconv.ParseUint(field, 10, 64)
		if err != nil || articleID == 0 {
			continue
		}

		redisCount, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil || redisCount <= 0 {
			continue
		}

		article, err := s.articleRepo.FindByID(uint(articleID))
		if err != nil || article == nil {
			logger.Error("Find article failed during sync",
				zap.Uint("articleID", uint(articleID)),
				zap.Error(err),
			)
			continue
		}

		newCount := uint(article.ViewCount) + uint(redisCount)
		if err := s.articleRepo.UpdateViewCount(uint(articleID), newCount); err != nil {
			logger.Error("Update view count failed",
				zap.Uint("articleID", uint(articleID)),
				zap.Error(err),
			)
			continue
		}

		if err := s.viewCountRepo.Delete(uint(articleID)); err != nil {
			logger.Error("Delete Redis view count failed",
				zap.Uint("articleID", uint(articleID)),
				zap.Error(err),
			)
		}

		logger.Info("Synced view count",
			zap.Uint("articleID", uint(articleID)),
			zap.Int64("redisIncrement", redisCount),
			zap.Uint("newTotal", newCount),
		)
	}

	return nil
}

func (s *viewCountService) GetTotalViewCount(articleID uint) uint {
	article, err := s.articleRepo.FindByID(articleID)
	if err != nil || article == nil {
		return 0
	}

	mysqlCount := uint(article.ViewCount)

	redisCount, err := s.viewCountRepo.Get(articleID)
	if err != nil || redisCount <= 0 {
		return mysqlCount
	}

	newTotal := mysqlCount + uint(redisCount)
	if err := s.articleRepo.UpdateViewCount(articleID, newTotal); err == nil {
		if delErr := s.viewCountRepo.Delete(articleID); delErr != nil {
			logger.Error("Delete Redis view count failed after sync",
				zap.Uint("articleID", articleID),
				zap.Error(delErr),
			)
		}
		logger.Info("Synced view count on read",
			zap.Uint("articleID", articleID),
			zap.Int64("redisCount", redisCount),
			zap.Uint("newTotal", newTotal),
		)
		return newTotal
	}

	return mysqlCount
}

func (s *viewCountService) StartScheduledSync(stopChan <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	logger.Info("View count scheduled sync started")

	for {
		select {
		case <-stopChan:
			logger.Info("View count scheduled sync stopped")
			if err := s.SyncViewCounts(); err != nil {
				logger.Error("Final sync failed on stop", zap.Error(err))
			}
			return
		case <-ticker.C:
			if err := s.SyncViewCounts(); err != nil {
				logger.Error("Scheduled sync failed", zap.Error(err))
			}
		}
	}
}

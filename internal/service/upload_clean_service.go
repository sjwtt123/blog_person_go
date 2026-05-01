package service

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/internal/repository"
	"blog/pkg/logger"

	"go.uber.org/zap"
)

type uploadCleanService struct {
	articleRepo   repository.ArticleRepository
	userRepo      repository.UserRepository
	uploadAbsPath string
	uploadDirs    []string
}

func NewUploadCleanService(articleRepo repository.ArticleRepository, userRepo repository.UserRepository, uploadPath string) UploadCleanService {
	absPath, _ := filepath.Abs(uploadPath)
	return &uploadCleanService{
		articleRepo:   articleRepo,
		userRepo:      userRepo,
		uploadAbsPath: absPath,
		uploadDirs: []string{
			filepath.Join(absPath, "articles"),
			filepath.Join(absPath, "avatars"),
			filepath.Join(absPath, "covers"),
		},
	}
}

func (s *uploadCleanService) CleanUnusedImages() (int, error) {
	deleted := 0

	allUsedUrls, err := s.collectAllUsedImageUrls()
	if err != nil {
		logger.Error("收集已使用图片链接失败", zap.Error(err))
		return 0, err
	}

	for _, dir := range s.uploadDirs {
		n, err := s.cleanDirectory(dir, allUsedUrls)
		if err != nil {
			logger.Error("清理目录失败", zap.String("dir", dir), zap.Error(err))
			continue
		}
		deleted += n
	}

	logger.Info("图片清理完成",
		zap.Int("deletedCount", deleted),
	)

	return deleted, nil
}

func (s *uploadCleanService) collectAllUsedImageUrls() (map[string]bool, error) {
	usedUrls := make(map[string]bool)

	articles, err := s.articleRepo.ListAll()
	if err != nil {
		return nil, err
	}

	for _, a := range articles {
		s.extractUrls(a.Content, usedUrls)
		s.extractUrls(a.CoverImage, usedUrls)
	}

	users, err := s.userRepo.ListAll()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		s.extractUrls(u.Avatar, usedUrls)
	}

	return usedUrls, nil
}

func (s *uploadCleanService) extractUrls(text string, urls map[string]bool) {
	if text == "" {
		return
	}

	// 支持 /uploads/... 路径
	s.extractUrlPattern(text, "/uploads/", urls)

	// 支持自定义 URL 路径
	s.extractUrlPattern(text, "/static/uploads/", urls)
}

func (s *uploadCleanService) extractUrlPattern(text string, pattern string, urls map[string]bool) {
	idx := 0
	for idx < len(text) {
		pos := strings.Index(text[idx:], pattern)
		if pos == -1 {
			break
		}

		start := idx + pos
		end := start + len(pattern)

		for end < len(text) {
			c := text[end]
			if c == ' ' || c == '"' || c == '\'' || c == ')' || c == '(' || c == '\n' || c == '\r' {
				break
			}
			end++
		}

		url := text[start:end]
		urls[url] = true
		idx = end
	}
}

func (s *uploadCleanService) cleanDirectory(dir string, usedUrls map[string]bool) (int, error) {
	deleted := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !s.isImageFile(path) {
			return nil
		}

		if s.isFileOld(path) {
			relPath, err := filepath.Rel(s.uploadAbsPath, path)
			if err != nil {
				return nil
			}
			url := "/uploads/" + filepath.ToSlash(relPath)
			if !usedUrls[url] {
				logger.Warn("发现未使用图片", zap.String("path", path), zap.String("url", url))
			}
		}

		return nil
	})

	return deleted, err
}

func (s *uploadCleanService) isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}

func (s *uploadCleanService) isFileOld(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) > 24*time.Hour
}

func (s *uploadCleanService) StartScheduledClean(stopChan <-chan struct{}) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	logger.Info("定时图片清理任务启动（每 24 小时执行一次）")

	for {
		select {
		case <-stopChan:
			logger.Info("定时图片清理任务停止")
			return
		case <-ticker.C:
			if _, err := s.CleanUnusedImages(); err != nil {
				logger.Error("定时图片清理失败", zap.Error(err))
			}
		}
	}
}

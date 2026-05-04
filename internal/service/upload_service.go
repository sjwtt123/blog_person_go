package service

import (
	"blog/pkg/logger"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/pkg/utils"
)

type uploadService struct {
	basePath string
}

func NewUploadService(basePath string) UploadService {
	return &uploadService{
		basePath: basePath,
	}
}

func (s *uploadService) UploadFile(file *multipart.FileHeader, uploadType UploadType, basePath string) (string, error) {
	subDir, maxSize, allowedExists, err := s.getUploadConfig(uploadType)
	if err != nil {
		return "", err
	}

	if err := s.validateFile(file, maxSize, allowedExists); err != nil {
		return "", err
	}

	savePath, err := s.saveFile(file, basePath, subDir)
	if err != nil {
		return "", err
	}

	return "/" + basePath + "/" + filepath.ToSlash(filepath.Join(subDir, filepath.Base(savePath))), nil
}

func (s *uploadService) getUploadConfig(uploadType UploadType) (string, int64, []string, error) {
	switch uploadType {
	case UploadTypeAvatar:
		return "avatars", 2 * 1024 * 1024, []string{".jpg", ".jpeg", ".png", ".webp"}, nil
	case UploadTypeCover:
		return "covers", 5 * 1024 * 1024, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}, nil
	case UploadTypeArticle:
		return "articles", 5 * 1024 * 1024, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}, nil
	default:
		return "", 0, nil, fmt.Errorf("图片类型错误")
	}
}

func (s *uploadService) validateFile(file *multipart.FileHeader, maxSize int64, allowedExts []string) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !utils.Contains(allowedExts, ext) {
		return fmt.Errorf("不支持的文件格式")
	}

	if file.Size > maxSize {
		return fmt.Errorf("文件过大")
	}

	return nil
}

func (s *uploadService) saveFile(file *multipart.FileHeader, basePath, subDir string) (string, error) {
	uploadDir := filepath.Join(basePath, subDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	randomStr, err := utils.GenerateRandomString(8)
	if err != nil {
		randomStr = time.Now().Format("150405")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := time.Now().Format("20060102150405") + "_" + randomStr + ext
	savePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			logger.Errorf("关闭文件夹失败%v", err)
		}
	}(src)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			logger.Errorf("关闭文件夹失败%v", err)
		}
	}(dst)

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return savePath, nil
}

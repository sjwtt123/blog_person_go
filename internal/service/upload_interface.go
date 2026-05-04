package service

import "mime/multipart"

type UploadType string

const (
	UploadTypeAvatar  UploadType = "avatar"
	UploadTypeCover   UploadType = "cover"
	UploadTypeArticle UploadType = "article"
)

type UploadService interface {
	UploadFile(file *multipart.FileHeader, uploadType UploadType, basePath string) (string, error)
}

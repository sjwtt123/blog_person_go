package service

// UploadCleanService 上传文件清理服务接口
type UploadCleanService interface {
	CleanUnusedImages() (deleted int, err error)
	StartScheduledClean(stopChan <-chan struct{})
}

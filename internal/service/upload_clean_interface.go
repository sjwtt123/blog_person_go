package service

// UploadCleanService 上传文件清理服务接口
type UploadCleanService interface {
	StartScheduledClean(stopChan <-chan struct{})
}

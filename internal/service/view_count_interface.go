package service

type ViewCountService interface {
	IncrementViewCount(articleID uint) error

	SyncViewCounts() error

	GetTotalViewCount(articleID uint) uint

	StartScheduledSync(stopChan <-chan struct{})
}

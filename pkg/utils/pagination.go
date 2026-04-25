package utils

// DefaultPage 默认页码
const DefaultPage = 1

// DefaultSize 默认每页大小
const DefaultSize = 10

// MaxSize 最大每页大小
const MaxSize = 100

// NormalizePage 标准化页码和大小，返回标准化后的值
func NormalizePage(page, size int) (int, int) {
	if page < DefaultPage {
		page = DefaultPage
	}
	if size < DefaultSize {
		size = DefaultSize
	} else if size > MaxSize {
		size = MaxSize
	}
	return page, size
}

// CalculateOffset 计算偏移量
func CalculateOffset(page, size int) int {
	return (page - 1) * size
}

// NormalizeAndOffset 标准化页码并计算偏移量
func NormalizeAndOffset(page, size int) (normPage int, normSize int, offset int) {
	normPage, normSize = NormalizePage(page, size)
	offset = CalculateOffset(normPage, normSize)
	return
}

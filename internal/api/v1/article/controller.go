package article

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/logger"
	"blog/pkg/response"
	"blog/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadType string

const (
	UploadTypeAvatar  UploadType = "avatar"
	UploadTypeCover   UploadType = "cover"
	UploadTypeArticle UploadType = "article"
)

type Controller struct {
	articleService     service.ArticleService
	userService        service.UserService
	viewCountService   service.ViewCountService
	uploadCleanService service.UploadCleanService
}

func NewController(articleService service.ArticleService, userService service.UserService, viewCountService service.ViewCountService, uploadCleanService service.UploadCleanService) *Controller {
	return &Controller{
		articleService:     articleService,
		userService:        userService,
		viewCountService:   viewCountService,
		uploadCleanService: uploadCleanService,
	}
}

func (ctrl *Controller) ListPublic(c *gin.Context) {
	var req request.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	data, err := ctrl.articleService.ListPublic(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}
	data, err := ctrl.articleService.GetDetail(uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) ListAdmin(c *gin.Context) {
	var req request.AdminArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	data, err := ctrl.articleService.ListAdmin(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req request.ArticleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	authorID := middleware.GetUserID(c)
	err := ctrl.articleService.Create(authorID, &req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, "文章创建成功")
}

func (ctrl *Controller) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}
	var req request.ArticleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err = ctrl.articleService.Update(uint(id), &req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}
	err = ctrl.articleService.Delete(uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (ctrl *Controller) GetTimelineItems(c *gin.Context) {
	data, err := ctrl.articleService.GetTimelineItems()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "文件上传失败")
		return
	}

	uploadType := UploadType(c.PostForm("type"))
	if uploadType == "" {
		response.BadRequest(c, "图片类型错误")
		return
	}

	var uploadDir string
	var maxSize int64
	var allowedExits []string

	switch uploadType {
	case UploadTypeAvatar:
		uploadDir = "uploads/avatars/" + time.Now().Format("2006/01")
		maxSize = 2 * 1024 * 1024
		allowedExits = []string{".jpg", ".jpeg", ".png", ".webp"}
	case UploadTypeCover:
		uploadDir = "uploads/covers/" + time.Now().Format("2006/01")
		maxSize = 5 * 1024 * 1024
		allowedExits = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	case UploadTypeArticle:
		uploadDir = "uploads/articles/" + time.Now().Format("2006/01")
		maxSize = 5 * 1024 * 1024
		allowedExits = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	default:
		response.BadRequest(c, "图片类型错误")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !utils.Contains(allowedExits, ext) {
		response.BadRequest(c, "不支持的文件格式")
		return
	}

	if file.Size > maxSize {
		response.BadRequest(c, "文件过大")
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "创建目录失败")
		return
	}

	randomStr, err := utils.GenerateRandomString(8)
	if err != nil {
		randomStr = time.Now().Format("150405")
	}
	filename := time.Now().Format("20060102150405") + "_" + randomStr + ext
	join := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, join); err != nil {
		response.InternalError(c, "保存文件失败")
		return
	}

	url := "/" + filepath.ToSlash(join)

	response.Success(c, gin.H{
		"url": url,
	})
}

// IncrementViewCount 增加文章访问量并返回最新值
func (ctrl *Controller) IncrementViewCount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}

	// 1. 先增加 Redis 中的计数
	if err := ctrl.viewCountService.IncrementViewCount(uint(id)); err != nil {
		logger.Error("Failed to increment view count in Redis", zap.Error(err))
	}

	// 2. 同步 Redis 和 MySQL，返回最新值
	total := ctrl.viewCountService.GetTotalViewCount(uint(id))

	response.Success(c, gin.H{
		"view_count": total,
	})
}

// CleanUnusedImages 清理未使用的图片
func (ctrl *Controller) CleanUnusedImages(c *gin.Context) {
	deleted, err := ctrl.uploadCleanService.CleanUnusedImages()
	if err != nil {
		response.InternalError(c, "清理失败")
		return
	}

	response.Success(c, gin.H{
		"deleted": deleted,
		"message": "清理完成",
	})
}

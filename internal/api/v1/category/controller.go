package category

import (
	"strconv"

	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	categoryService service.CategoryService
	userService     service.UserService
}

func NewController(categoryService service.CategoryService, userService service.UserService) *Controller {
	return &Controller{
		categoryService: categoryService,
		userService:     userService,
	}
}

func (ctrl *Controller) ListAll(c *gin.Context) {
	list, err := ctrl.categoryService.ListAll()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req request.CategoryUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := ctrl.categoryService.Create(&req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (ctrl *Controller) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}
	var req request.CategoryUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := ctrl.categoryService.Update(uint(id), &req); err != nil {
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
	if err := ctrl.categoryService.Delete(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

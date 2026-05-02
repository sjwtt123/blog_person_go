# 后端开发规范

本文档定义本项目的后端开发规范，所有后端代码必须严格遵守。

## 1. 架构分层

```
Controller (api/v1/) → Service → Repository → Database
```

**核心规则：**
- Controller：处理 HTTP 请求/响应、参数校验、调用 Service
- Service：业务逻辑、权限判断、组合多个 Repository 操作
- Repository：数据库 CRUD 操作，不包含业务逻辑
- **禁止跨层调用**：Controller 不能直接调 Repository，Service 不能调 Controller

## 2. 目录结构

每个模块（如 article、comment、category）遵循以下结构：

```
internal/
├── api/v1/{module}/
│   ├── controller.go    # HTTP 处理 + RegisterRoutes
│   └── router.go        # 可选，路由注册可放在 controller.go
├── service/
│   ├── {module}_interface.go  # 接口定义
│   └── {module}_service.go    # 接口实现
├── repository/
│   ├── {module}_interface.go  # 接口定义
│   └── {module}_repository.go # 接口实现
└── model/
    ├── entity/
    │   └── {module}.go    # 数据库实体
    └── dto/
        ├── request/
        │   └── {module}.go  # 请求 DTO
        └── response/
            └── {module}.go  # 响应 DTO
```

## 3. 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件名 | 小写 + 下划线 | `user_service.go` |
| 包名 | 小写单词 | `service`, `repository` |
| 接口名 | 能力名 + Repository/Service | `UserRepository`, `UserService` |
| 实现结构体 | 小写接口名去掉 interface | `userService`, `userRepository` |
| Controller 结构体 | `Controller` | `Controller` |
| DTO 结构体 | `XxxRequest`, `XxxResponse` | `CreateArticleRequest` |
| 方法名 | 动词开头，大驼峰 | `ListPublic`, `GetDetail` |
| 常量 | 大驼峰 | `DefaultPage`, `MaxSize` |
| 变量 | 小驼峰 | `articleID`, `userID` |

## 4. Controller 规范

### 4.1 结构体定义

```go
type Controller struct {
    xxxService service.XxxService
    userService    service.UserService  // 权限判断需要
}

func NewController(xxxService service.XxxService, userService service.UserService) *Controller {
    return &Controller{
        xxxService: xxxService,
        userService:    userService,
    }
}
```

### 4.2 路由注册

- 公开接口直接注册在 `r` 上
- 需要登录的接口使用 `r.Group("/all", middleware.Auth())` 分组
- 需要管理员权限的接口使用 `r.Group("/admin", middleware.Auth(), middleware.Admin(ctrl.userService))` 分组

```go
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
    // 公开接口
    r.GET("/articles", ctrl.ListPublic)
    r.GET("/articles/:id", ctrl.GetDetail)

    // 需要登录
    allGroup := r.Group("/all", middleware.Auth())
    {
        allGroup.POST("/articles", ctrl.Create)
        allGroup.PUT("/articles/:id", ctrl.Update)
        allGroup.DELETE("/articles/:id", ctrl.Delete)
    }
}
```

### 4.3 参数获取

**路径参数：**
```go
id, err := strconv.ParseUint(c.Param("id"), 10, 64)
if err != nil || id == 0 {
    response.BadRequest(c, "id 参数错误")
    return
}
```

**查询参数（分页、筛选等）使用 DTO + ShouldBindQuery：**
```go
var req request.ArticleListRequest
if err := c.ShouldBindQuery(&req); err != nil {
    response.BadRequest(c, err.Error())
    return
}
```

**请求体使用 ShouldBindJSON：**
```go
var req request.CreateCommentRequest
if err := c.ShouldBindJSON(&req); err != nil {
    response.BadRequest(c, err.Error())
    return
}
```

**获取当前用户 ID：**
```go
userID := middleware.GetUserID(c)
```

### 4.4 错误处理

```go
// 业务错误统一用 BizError
data, err := ctrl.xxxService.DoSomething()
if err != nil {
    response.BizError(c, err)
    return
}

// 成功响应
response.Success(c, data)
```

### 4.5 成功响应

- 返回列表/详情：`response.Success(c, data)`
- 创建/更新/删除成功：`response.Success(c, nil)` 或 `response.Success(c, "操作成功")`
- 分页数据：Service 层返回 `*response.PageResponse`，Controller 直接 `response.Success(c, data)`

## 5. Service 规范

### 5.1 接口定义

```go
type XxxService interface {
    Create(userID uint, req *request.CreateXxxRequest) (*response.XxxResponse, error)
    GetByID(id uint) (*response.XxxResponse, error)
    List(req *request.XxxListRequest) (*response.PageResponse, error)
    Update(id uint, req *request.UpdateXxxRequest) error
    Delete(id uint) error
}
```

### 5.2 分页处理

**必须使用 `utils.NormalizeAndOffset` 和 `response.NewPageResponse`：**

```go
func (s *xxxService) List(req *request.XxxListRequest) (*response.PageResponse, error) {
    page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

    list, total, err := s.xxxRepo.List(offset, size, req.Keyword)
    if err != nil {
        return nil, err
    }

    return response.NewPageResponse(list, total, page, size), nil
}
```

### 5.3 权限控制

- 在 Service 层做权限判断（如：只能修改自己的数据）
- 通过 `userRepo.FindByID(userID)` 获取用户角色判断权限

### 5.4 错误返回

使用 `bizerrors.New(code, message)` 返回业务错误：

```go
if article == nil {
    return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
}
```

## 6. Repository 规范

### 6.1 接口定义

```go
type XxxRepository interface {
    FindByID(id uint) (*entity.Xxx, error)
    List(offset, limit int) ([]*entity.Xxx, int64, error)
    Create(xxx *entity.Xxx) error
    Update(xxx *entity.Xxx) error
    Delete(id uint) error
}
```

### 6.2 FindByID 处理未找到

```go
func (r *xxxRepository) FindByID(id uint) (*entity.Xxx, error) {
    var xxx entity.Xxx
    err := r.db.First(&xxx, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil  // 返回 nil, nil 表示未找到
        }
        return nil, err
    }
    return &xxx, nil
}
```

### 6.3 列表查询返回 total

```go
func (r *xxxRepository) List(offset, limit int) ([]*entity.Xxx, int64, error) {
    var list []*entity.Xxx
    var total int64

    query := r.db.Model(&entity.Xxx{})
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
        return nil, 0, err
    }

    return list, total, nil
}
```

## 7. DTO 规范

### 7.1 请求 DTO

- 放在 `internal/model/dto/request/`
- 使用 struct tag 定义 JSON 绑定和校验规则
- 查询参数使用 `form:"xxx"` tag
- 请求体使用 `json:"xxx"` tag + `binding:"required"` 等校验

```go
type CreateCommentRequest struct {
    Content string `json:"content" binding:"required,min=1,max=1000"`
    ParentID uint  `json:"parent_id"`
}

type CommentListRequest struct {
    Page int `form:"page"`
    Size int `form:"size"`
}
```

### 7.2 响应 DTO

- 放在 `internal/model/dto/response/`
- 使用 `json:"xxx"` tag
- 分页响应使用 `response.PageResponse`

```go
type PageResponse struct {
    List      interface{} `json:"list"`
    Total     int64       `json:"total"`
    Page      int         `json:"page"`
    Size      int         `json:"size"`
    TotalPage int         `json:"total_page"`
}
```

## 8. 实体规范

```go
type Comment struct {
    BaseEntity          // 继承 ID, CreatedAt, UpdatedAt
    ArticleID uint   `gorm:"not null;index;comment:文章ID" json:"article_id"`
    UserID    uint   `gorm:"not null;index;comment:用户ID" json:"user_id"`
    Content   string `gorm:"type:varchar(1000);not null;comment:评论内容" json:"content"`
    Status    int    `gorm:"type:tinyint;default:1;comment:状态" json:"status"`

    Article *Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
    User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Comment) TableName() string {
    return "comments"
}
```

## 9. 路由注册流程

1. 在 `internal/api/v1/{module}/controller.go` 中定义 `RegisterRoutes` 方法
2. 在 `internal/api/v1/{module}/router.go` 中可选定义路由分组（也可在 controller.go 中）
3. 在 `internal/api/router.go` 的 `Router` 结构体中添加 controller 字段
4. 在 `NewRouter` 中创建 controller 实例
5. 在 `Setup` 方法中调用 `ctrl.RegisterRoutes(v1)`
6. 在 `internal/app/app.go` 的 `initDependencies` 中创建 Repository 和 Service 并注入

**示例**:
```go
// router.go
type Router struct {
    articleCtrl    *article.Controller
    commentCtrl    *comment.Controller
    likeCtrl       *like.Controller
    // ...
}

func NewRouter(...) *Router {
    return &Router{
        articleCtrl: article.NewController(articleSvc),
        commentCtrl: comment.NewController(commentSvc, userSvc),
        likeCtrl:    like.NewController(likeSvc, userSvc),
    }
}

func (r *Router) Setup(engine *gin.Engine) {
    v1 := engine.Group("/api/v1")
    r.articleCtrl.RegisterRoutes(v1)
    r.commentCtrl.RegisterRoutes(v1)
    r.likeCtrl.RegisterRoutes(v1)
}
```

## 10. 错误码规范

错误码定义在 `pkg/errors/code.go`：

- `0` - 成功
- `400-599` - HTTP 级别错误
- `1000-1999` - 用户相关
- `2000-2999` - 参数相关
- `3000-3999` - 资源相关（文章、分类、标签、评论等）
- `3000-3099` - 点赞相关

新增错误码需添加到 `code.go`。

**常用错误码**:
- `1001` - 用户不存在
- `1002` - 用户已存在
- `1003` - 用户名或密码错误
- `1004` - 用户已被禁用
- `1005` - 无效的令牌
- `1006` - 令牌已过期
- `2001` - 文章不存在
- `2002` - 文章已存在
- `2003` - 分类不存在
- `2004` - 分类已存在
- `2005` - 标签不存在
- `2006` - 标签已存在
- `2007` - 评论不存在
- `2008` - 禁止评论
- `3001` - 已点赞
- `3002` - 未点赞
- `4001` - 文件上传失败
- `4002` - 不支持的文件格式
- `4003` - 文件过大

## 11. 依赖注入

所有 Service 和 Repository 通过构造函数注入，不使用全局变量：

```go
// app.go
commentRepo := repository.NewCommentRepository(a.mysqlDB)
commentSvc := service.NewCommentService(commentRepo, userRepo, articleRepo)
a.router = api.NewRouter(..., commentSvc)
```

## 12. 日志

- 使用 `blog/pkg/logger` 包
- 关键操作记录日志：`logger.Info("xxx", zap.Any("key", value))`
- 错误记录：`logger.Error("xxx", zap.Error(err))`

## 13. 编译检查

每次修改后端代码后必须执行：
```bash
go build -o bin/server.exe ./cmd/server
```
确保编译通过后再提交。

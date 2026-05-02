# 架构设计文档

## 概述

CloudQue 采用标准的 MVC 三层架构设计，实现了清晰的分层结构和依赖注入。

## 架构分层

```
┌─────────────────────────────────────────┐
│          API Layer (Controller)          │  HTTP 请求处理
├─────────────────────────────────────────┤
│         Service Layer (Business)         │  业务逻辑处理
├─────────────────────────────────────────┤
│      Repository Layer (Data Access)      │  数据访问
├─────────────────────────────────────────┤
│         Database (MySQL + Redis)         │  数据存储
└─────────────────────────────────────────┘
```

## 层次说明

### 1. API 层 (Controller)

**位置**: `internal/api/`

**职责**:
- 处理 HTTP 请求和响应
- 参数绑定和验证
- 调用 Service 层处理业务
- 返回统一格式的响应

**示例**:
```go
func (ctrl *Controller) GetProfile(c *gin.Context) {
    userID := middleware.GetUserID(c)
    user, err := ctrl.userService.GetUserByID(userID)
    // ...
}
```

### 2. Service 层

**位置**: `internal/service/`

**职责**:
- 实现核心业务逻辑
- 事务管理
- 调用 Repository 层进行数据操作
- 业务规则验证

**示例**:
```go
func (s *userService) Register(req *request.RegisterRequest) error {
    // 业务验证
    exists, _ := s.userRepo.ExistsByUsername(req.Username)
    if exists {
        return errors.ErrUserAlreadyExists
    }
    // 业务处理
    // ...
}
```

### 3. Repository 层

**位置**: `internal/repository/`

**职责**:
- 封装数据访问逻辑
- CRUD 操作
- 数据查询和转换

**示例**:
```go
func (r *userRepository) FindByID(id uint) (*entity.User, error) {
    var user entity.User
    err := r.db.First(&user, id).Error
    // ...
}
```

## 设计模式

### 依赖注入

所有层之间通过构造函数注入依赖，便于测试和扩展。

```go
// Controller 注入 Service
func NewController(userService service.UserService) *Controller {
    return &Controller{userService: userService}
}

// Service 注入 Repository
func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}
```

### 接口隔离

Service 和 Repository 层使用接口定义，支持 Mock 测试和多种实现。

```go
// 定义接口
type UserRepository interface {
    FindByID(id uint) (*entity.User, error)
    Create(user *entity.User) error
}

// 实现接口
type userRepository struct {
    db *gorm.DB
}
```

## 核心组件

### 1. 配置管理

**位置**: `pkg/config/`

- 使用 Viper 加载 YAML 配置
- 支持环境变量覆盖
- 支持多环境配置（dev/prod）

### 2. 日志系统

**位置**: `pkg/logger/`

- 使用 Zap 结构化日志
- 集成 Lumberjack 实现日志轮转
- 支持控制台和文件输出

### 3. 数据库管理

**位置**: `pkg/database/`

- MySQL 连接池配置
- GORM 配置和初始化
- Redis 连接管理

### 4. 认证系统

**位置**: `pkg/jwt/`, `internal/middleware/auth.go`

- JWT Token 生成和验证
- 认证中间件
- 从 Token 提取用户信息

### 5. 错误处理

**位置**: `pkg/errors/`

- 自定义错误类型
- 统一错误码管理
- 业务错误和系统错误分离

### 6. 中间件系统

**位置**: `internal/middleware/`

- **Logger**: 记录请求日志
- **Recovery**: Panic 恢复
- **CORS**: 跨域处理
- **Auth**: JWT 认证（需要登录）
- **OptionalAuth**: 可选认证（有 token 则解析，无 token 则跳过）
- **Admin**: 管理员权限验证

## 核心模块说明

### 点赞模块

点赞功能采用分离设计，将公开数据和用户数据分开处理：

#### 1. 公开点赞数接口

- **路径**: `GET /api/v1/likes/article/:id/count`
- **认证**: 无需认证
- **功能**: 获取文章的点赞总数
- **特点**: 高性能，无需解析 token

#### 2. 用户点赞状态接口

- **路径**: `GET /api/v1/likes/article/:id/status`
- **认证**: 需要认证
- **功能**: 判断当前用户是否已点赞
- **特点**: 需要解析 token 获取 userID

#### 3. 点赞/取消点赞接口

- **路径**: `POST /api/v1/likes/article/:id` 和 `DELETE /api/v1/likes/article/:id`
- **认证**: 需要认证
- **功能**: 点赞或取消点赞
- **特点**: 使用数据库事务保证数据一致性

#### 4. 用户点赞列表接口

- **路径**: `GET /api/v1/likes/my`
- **认证**: 需要认证
- **功能**: 获取用户点赞过的文章列表
- **特点**: 分页查询，支持文章详情关联

#### 核心业务逻辑

1. **点赞操作**:
   - 检查用户是否已点赞（通过 likes 表唯一索引）
   - 如未点赞，插入点赞记录并增加文章 like_count
   - 使用数据库事务保证原子性

2. **取消点赞**:
   - 软删除点赞记录（设置 deleted_at）
   - 减少文章 like_count
   - 使用数据库事务保证原子性

3. **防重复点赞**:
   - 数据库层面：`(article_id, user_id)` 联合唯一索引
   - 应用层面：点赞前检查是否已点赞

## 数据模型

### Entity (实体)

**位置**: `internal/model/entity/`

对应数据库表结构，使用 GORM 标签配置：

```go
type User struct {
    BaseEntity
    Username string `gorm:"type:varchar(50);uniqueIndex"`
    Password string `gorm:"type:varchar(255)"`
    // ...
}
```

### 数据库表结构

#### 1. users - 用户表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| username | varchar(50) | 用户名（唯一） |
| password | varchar(255) | 密码（bcrypt 加密） |
| email | varchar(100) | 邮箱（唯一） |
| nickname | varchar(50) | 昵称 |
| avatar | varchar(255) | 头像 URL |
| role | varchar(20) | 角色（admin/user） |
| status | tinyint | 状态（1 正常，2 禁用） |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

#### 2. articles - 文章表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| title | varchar(255) | 文章标题 |
| slug | varchar(255) | 文章别名（URL 友好） |
| summary | text | 文章摘要 |
| content | longtext | 文章内容（Markdown） |
| cover_image | varchar(255) | 封面图片 URL |
| author_id | uint | 作者 ID |
| category_id | uint | 分类 ID |
| status | tinyint | 状态（1 发布，2 草稿，3 隐藏） |
| view_count | int | 浏览量 |
| like_count | int | 点赞数 |
| comment_count | int | 评论数 |
| published_at | datetime | 发布时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

#### 3. categories - 分类表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | varchar(50) | 分类名称 |
| slug | varchar(50) | 分类别名（唯一） |
| description | varchar(255) | 分类描述 |
| sort_order | int | 排序权重 |
| post_count | int | 文章数量 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

#### 4. tags - 标签表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | varchar(50) | 标签名称 |
| slug | varchar(50) | 标签别名（唯一） |
| description | varchar(255) | 标签描述 |
| post_count | int | 文章数量 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

#### 5. article_tags - 文章标签关联表

| 字段 | 类型 | 说明 |
|------|------|------|
| article_id | uint | 文章 ID |
| tag_id | uint | 标签 ID |

#### 6. comments - 评论表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| article_id | uint | 文章 ID |
| user_id | uint | 用户 ID |
| parent_id | uint | 父评论 ID（0 表示顶级评论） |
| reply_to_id | uint | 回复的评论 ID |
| reply_to_name | varchar(50) | 回复的用户名 |
| content | varchar(1000) | 评论内容 |
| status | tinyint | 状态（1 正常，2 隐藏） |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

#### 7. likes - 点赞表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| article_id | uint | 文章 ID |
| user_id | uint | 用户 ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

**索引设计**:
- 联合唯一索引：`(article_id, user_id)` - 防止重复点赞
- 普通索引：`article_id` - 快速查询文章点赞记录
- 普通索引：`user_id` - 快速查询用户点赞记录

### DTO (数据传输对象)

**位置**: `internal/model/dto/`

- **Request**: API 请求数据结构
- **Response**: API 响应数据结构

## 请求流程

```
HTTP Request
    ↓
[Middleware] → Recovery, Logger, CORS, Auth
    ↓
[Controller] → 参数绑定和验证
    ↓
[Service] → 业务逻辑处理
    ↓
[Repository] → 数据库操作
    ↓
Database
    ↓
[Response] → 统一响应格式
    ↓
HTTP Response
```

## 数据库设计

### 表结构

- **users**: 用户表
  - 基础字段：id, created_at, updated_at, deleted_at
  - 业务字段：username, password, email, nickname, avatar, status

### 索引设计

- 唯一索引：username, email
- 普通索引：deleted_at

## 安全设计

1. **密码加密**: 使用 bcrypt 加密存储
2. **JWT 认证**: Token 有效期控制
3. **CORS 配置**: 限制跨域访问
4. **参数验证**: 使用 validator 标签
5. **SQL 注入防护**: 使用 GORM 参数化查询

## 性能优化

1. **连接池**: MySQL 和 Redis 连接池配置
2. **日志轮转**: 避免日志文件过大
3. **优雅关闭**: 等待现有请求完成
4. **跳过默认事务**: GORM 配置提升性能

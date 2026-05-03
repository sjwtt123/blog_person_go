# MyBlog API - 个人博客系统

基于 Go 语言和 Gin 框架构建的个人博客系统后端项目，采用标准三层 MVC 架构，提供完整的博客管理、用户认证、评论互动等功能。

## 项目简介

是一个功能完整的个人博客系统后端，采用标准的 MVC 分层设计，支持文章管理、分类标签、评论点赞、用户认证等核心功能，可作为个人博客项目的后端基础。

### 架构特点

- **应用启动器模式**：通过 `internal/app/app.go` 统一管理应用初始化流程，main.go 仅需 20 行代码
- **清晰的分层架构**：Controller → Service → Repository，职责明确
- **依赖注入**：所有层次通过构造函数注入，便于测试和扩展
- **接口隔离**：Service 和 Repository 层使用接口定义，支持多种实现

## 技术栈

- **Web 框架**: Gin v1.11.0
- **ORM**: GORM v1.31.1
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **配置管理**: Viper
- **日志**: Zap + Lumberjack
- **认证**: JWT
- **密码加密**: bcrypt
- **图片上传**: 本地存储/云存储

## 项目结构

```
blog-go/
├── cmd/
│   └── server/
│       └── main.go                    # 主程序入口（20行代码）
├── internal/                           # 私有应用代码
│   ├── app/                           # 应用启动器
│   │   └── app.go                     # 应用初始化、依赖注入、优雅关闭
│   ├── api/                           # API 层
│   │   ├── v1/                        # API v1 版本
│   │   │   ├── auth/                  # 认证模块（注册/登录/刷新Token）
│   │   │   ├── user/                  # 用户模块（个人信息/关于页/用户管理）
│   │   │   ├── article/               # 文章模块（CRUD/上传/时间轴/访问量）
│   │   │   ├── category/              # 分类模块
│   │   │   ├── tag/                   # 标签模块
│   │   │   └── comment/               # 评论模块
│   │   └── routeres.go                # 路由注册
│   ├── service/                       # 业务逻辑层（接口+实现）
│   │   ├── *_interface.go             # Service接口定义
│   │   └── *_service.go               # Service实现
│   ├── repository/                    # 数据访问层（接口+实现）
│   │   ├── *_interface.go             # Repository接口定义
│   │   └── *_repository.go            # Repository实现
│   ├── model/                         # 数据模型
│   │   ├── entity/                    # 数据库实体（User/Article/Category等）
│   │   └── dto/                       # 数据传输对象
│   │       ├── request/               # 请求参数DTO
│   │       └── response/              # 响应数据DTO
│   └── middleware/                    # 中间件（Auth/Admin/CORS/Logger/Recovery）
├── pkg/                               # 公共工具包
│   ├── config/                        # 配置管理（Viper）
│   ├── logger/                        # 日志系统（Zap + Lumberjack）
│   ├── database/                      # 数据库管理（MySQL + Redis）
│   ├── jwt/                           # JWT工具（生成/解析/刷新）
│   ├── response/                      # 统一响应格式
│   ├── errors/                        # 业务错误处理（BizError）
│   ├── util/                         # 工具函数（slug生成）
│   └── utils/                         # 工具函数（分页/字符串/时间）
├── configs/
│   └── config.yaml.example            # 配置文件示例
├── scripts/
│   └── migrate.sql                    # 数据库迁移脚本
├── docs/                              # 项目文档
├── uploads/                           # 上传文件存储
│   ├── articles/                      # 文章图片
│   ├── covers/                        # 封面图片
│   └── avatars/                       # 用户头像
├── Makefile                           # 构建命令
├── go.mod                             # Go 模块定义
└── README.md                          # 项目说明
```

### 架构设计说明

**后台管理路由设计：**
- 采用**分散式管理**，各业务模块内部包含自己的 admin 路由组
- 不使用独立的 `admin/` 目录，而是通过中间件 `middleware.Admin()` 控制权限
- Admin路由统一使用 `/api/v1/admin/*` 前缀，例如：
  - `/api/v1/admin/articles` - 文章管理
  - `/api/v1/admin/users` - 用户管理
  - `/api/v1/admin/categories` - 分类管理
  - `/api/v1/admin/comments` - 评论管理
## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 安装步骤

1. 克隆项目
```bash
git clone https://cloudque.git
cd cloudque
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库
```bash
# config.yaml.example重命名 config.yaml
# 编辑 configs/config.yaml，修改数据库连接信息
# 或使用环境变量覆盖
export MYSQL_PASSWORD=your_password
```

4. 创建数据库
```bash
mysql -u root -p < scripts/migrate.sql
```

5. 运行项目
```bash
go run cmd/server/main.go
```

或使用 Makefile：
```bash
make run
```

### 配置说明

配置文件位于 `configs/config.yaml`，主要配置项：

- `app`: 应用配置（名称、版本、端口）
- `database`: 数据库配置（MySQL、Redis）
- `jwt`: JWT 认证配置
- `log`: 日志配置
- `cors`: 跨域配置

支持通过环境变量覆盖敏感配置：
- `MYSQL_PASSWORD`
- `REDIS_PASSWORD`
- `JWT_SECRET`

### API 端点

#### 健康检查
```
GET /api/v1/health
```

#### 用户认证（公开）
```
POST /api/v1/auth/register  # 用户注册
POST /api/v1/auth/login     # 用户登录
POST /api/v1/auth/refresh   # 刷新 Token
```

#### 关于页面（公开）
```
GET /api/v1/about           # 获取关于页信息
```

#### 用户管理（需认证）
```
GET  /api/v1/user/profile    # 获取用户信息
PUT  /api/v1/user/profile    # 更新用户信息
POST /api/v1/user/password   # 修改密码
POST /api/v1/user/upload/image  # 上传图片（文章/封面/头像）
```

#### 文章管理
**公开接口：**
```
GET /api/v1/articles              # 获取文章列表（分页、筛选、搜索）
GET /api/v1/articles/:id          # 获取文章详情
GET /api/v1/articles/timelines    # 获取时间轴数据
PUT /api/v1/articles/:id/view     # 增加访问量（Redis缓存）
GET /api/v1/articles/category/:id # 按分类获取文章
GET /api/v1/articles/tag/:tag     # 按标签获取文章
```

**管理员接口：**
```
POST   /api/v1/admin/articles      # 创建文章
PUT    /api/v1/admin/articles/:id  # 更新文章
DELETE /api/v1/admin/articles/:id  # 删除文章
GET    /api/v1/admin/articles      # 管理员文章列表（含草稿/隐藏）
```

#### 文章分类
**公开接口：**
```
GET /api/v1/categories             # 获取所有分类（含文章数量）
```

**管理员接口：**
```
POST   /api/v1/admin/categories     # 创建分类
PUT    /api/v1/admin/categories/:id # 更新分类
DELETE /api/v1/admin/categories/:id # 删除分类
```

> ⚠️ **待实现**: `GET /api/v1/categories/:id` - 分类详情接口

#### 标签管理
**公开接口：**
```
GET /api/v1/tags                   # 获取所有标签（含文章数量）
```

**管理员接口：**
```
POST   /api/v1/admin/tags          # 创建标签
PUT    /api/v1/admin/tags/:id      # 更新标签
DELETE /api/v1/admin/tags/:id      # 删除标签
```

> ⚠️ **待实现**: `GET /api/v1/tags/:id` - 标签详情接口

#### 评论管理
**公开接口：**
```
GET /api/v1/articles/:id/comments  # 获取文章评论列表
```

**需认证接口：**
```
POST   /api/v1/user/articles/:id/comments  # 发布评论
PUT    /api/v1/user/comments/:id           # 更新自己的评论
DELETE /api/v1/user/comments/:id           # 删除自己的评论
```

**管理员接口：**
```
GET    /api/v1/admin/comments       # 获取所有评论（分页）
DELETE /api/v1/admin/comments/:id   # 删除任意评论
```

#### 点赞功能

**公开接口：**
```
GET /api/v1/likes/article/:id   # 获取文章点赞数和状态
```

**需认证接口：**
```
POST   /api/v1/likes/article/:id   # 点赞文章
DELETE /api/v1/likes/article/:id   # 取消点赞
GET    /api/v1/likes/user/articles # 获取用户点赞的文章列表
```

**实现说明：**
- 使用数据库事务保证点赞记录和点赞数的一致性
- 防止重复点赞和取消不存在的点赞
- 支持未登录用户查看点赞数

#### 后台管理 - 用户（需管理员权限）
```
GET    /api/v1/admin/users         # 获取用户列表
POST   /api/v1/admin/users         # 创建用户
PUT    /api/v1/admin/users/:id     # 更新用户
DELETE /api/v1/admin/users/:id     # 删除用户
PUT    /api/v1/admin/about         # 更新关于页内容
```
### API 测试

#### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com",
    "nickname": "测试用户"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

#### 获取用户信息（需要 Token）
```bash
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### 获取文章列表（分页、筛选）
```bash
# 获取所有文章（分页）
curl "http://localhost:8080/api/v1/articles?page=1&size=10"

# 根据标题搜索
curl "http://localhost:8080/api/v1/articles?keyword=Go 语言&page=1&size=10"

# 按分类筛选
curl "http://localhost:8080/api/v1/articles?category_id=1&page=1&size=10"

# 按标签筛选
curl "http://localhost:8080/api/v1/articles?tag=Go&page=1&size=10"

# 按时间范围筛选
curl "http://localhost:8080/api/v1/articles?start_date=2024-01-01&end_date=2024-12-31&page=1&size=10"
```

#### 获取文章详情
```bash
curl http://localhost:8080/api/v1/articles/1
```

#### 创建文章（需要 Token）
```bash
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Go 语言入门",
    "summary": "Go 语言基础教程",
    "content": "# Go 语言入门\n\nGo 是一种开源的编程语言...",
    "category_id": 1,
    "tags": ["Go", "编程", "教程"],
    "cover_image": "https://example.com/cover.jpg"
  }'
```

#### 发布评论（需要 Token）
```bash
curl -X POST http://localhost:8080/api/v1/user/articles/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "content": "写得很棒！学到了很多",
    "parent_id": 0
  }'
```

#### 上传图片（需要 Token）
```bash
curl -X POST http://localhost:8080/api/v1/user/upload/image \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "file=@/path/to/image.jpg"
```

#### 获取时间轴数据
```bash
curl http://localhost:8080/api/v1/articles/timelines
```

#### 增加文章访问量
```bash
curl -X PUT http://localhost:8080/api/v1/articles/1/view
```

## Makefile 命令

```bash
make build        # 编译项目
make run          # 运行项目
make test         # 运行测试
make clean        # 清理构建文件
make deps         # 下载依赖
make fmt          # 格式化代码
make vet          # 代码静态检查
make help         # 显示帮助
```

## 开发指南

详细的开发指南请查看：
- [架构设计](docs/ARCHITECTURE.md)
- [开发指南](docs/DEVELOPMENT.md)
- [API 文档](docs/API.md)

### 开发规范

#### 1. 分层架构规范
- **Controller 层**：仅负责参数绑定、校验、调用 Service、响应返回
- **Service 层**：业务逻辑处理、事务管理、权限校验
- **Repository 层**：数据库操作封装，不暴露 GORM 细节

#### 2. 函数长度规范
- 单个函数不超过 **50 行**
- 超过 50 行的函数必须拆分为职责单一的私有方法
- 公共逻辑提取为独立函数，命名以 `handle`、`validate`、`build`、`update` 等开头

#### 3. 错误处理规范
- **统一使用** `pkg/errors` 包中的 `BizError` 定义业务错误
- **禁止使用** `panic` 或 `log.Fatal` 处理业务错误
- **Controller 层**统一使用 `response.BizError(c, err)` 处理错误响应
- **Service 层**返回 `*errors.BizError` 类型，便于 Controller 识别

#### 4. 事务处理规范
- 事务在 **Repository 层**（数据层）的一个方法内完成
- 使用 `db.Transaction(func(tx *gorm.DB) error { ... })` 闭包保证原子性
- 相关数据库操作通过 `tx` 传递，确保在同一事务中
- **Service 层**仅调用 Repository 的事务方法，不直接管理事务

#### 5. 数据库操作规范
- **GORM 查询链**不超过 3 层调用
- 复杂查询封装在 Repository 内，不泄漏到 Service
- 使用 `Preload` 预加载关联，避免 N+1 查询

#### 6. 命名规范
- **文件名**：小写下划线，如 `user_service.go`
- **包名**：小写单词，如 `service`、`repository`
- **变量名**：驼峰命名，如 `userRepo`、`articleList`
- **函数名**：大驼峰公开，小驼峰私有，如 `Create`、`validateCategory`
- **接收器**：单字母缩写，如 `(s *userService)`、`(r *articleRepo)`

#### 7. DTO 复用规范
- **核心实体**在 `model/entity` 包统一定义
- **请求参数**在 `model/dto/request` 包定义
- **响应数据**在 `model/dto/response` 包定义
- 禁止在不同模块重复定义相同字段的结构体

## 核心特性

### 基础架构
- ✅ 标准 MVC 三层架构（Controller → Service → Repository）
- ✅ 应用启动器（App Launcher）- 统一管理初始化流程，main.go 仅 20 行
- ✅ 依赖注入设计（构造函数注入）
- ✅ 接口隔离原则（Service 和 Repository 都有接口定义）
- ✅ 优雅关闭（Graceful Shutdown）

### 用户系统
- ✅ JWT 认证机制（Access Token + Refresh Token）
- ✅ 用户注册/登录（用户名+密码）
- ✅ 密码加密存储（bcrypt）
- ✅ 用户信息管理（昵称、邮箱、头像）
- ✅ 管理员权限控制（基于角色的访问控制 RBAC）
- ✅ 关于页面管理

### 博客功能
- ✅ 文章发布与管理（草稿/发布/隐藏状态）
- ✅ 文章分页展示
- ✅ 多条件筛选（标题、分类、标签、状态）
- ✅ 文章搜索（关键词搜索）
- ✅ 文章分类（支持层级分类）
- ✅ 标签管理
- ✅ Markdown 格式支持
- ✅ 图片上传功能（文章图片、封面图、头像）
- ✅ Slug 友好URL生成
- ✅ **时间轴展示** - 按时间归档文章
- ✅ **访问量统计** - Redis缓存 + 定时同步MySQL（每5分钟）
- ✅ **图片清理服务** - 定时清理未使用图片（每24小时）

### 互动功能
- ✅ 评论系统（支持回复和嵌套评论）
- ✅ 评论事务处理（创建/删除时自动更新文章评论数）
- ✅ 评论仅登录用户可发
- ✅ 点赞功能（支持点赞/取消点赞/点赞数统计）

### 后台管理
- ✅ 文章增删改查（管理员可管理所有状态的文章）
- ✅ 分类管理（创建/编辑/删除）
- ✅ 标签管理
- ✅ 标签管理
- ✅ 评论管理
- ✅ 用户管理

### 技术特性
- ✅ 统一响应格式（`pkg/response`）
- ✅ 统一错误处理（`pkg/errors` - BizError）
- ✅ 分页查询支持（`pkg/utils/pagination.go`）
- ✅ 结构化日志（Zap + Lumberjack）
- ✅ 配置管理（Viper，支持环境变量覆盖）
- ✅ 数据库迁移（GORM AutoMigrate）
- ✅ 优雅关闭（Graceful Shutdown）
- ✅ CORS 跨域支持
- ✅ 中间件系统（Logger/Recovery/CORS/Auth/Admin）

---

## 数据库设计

### 主要数据表

| 表名 | 说明 | 状态 |
|------|------|------|
| **users** | 用户表 | ✅ 已使用 |
| **articles** | 文章表 | ✅ 已使用 |
| **categories** | 分类表 | ✅ 已使用 |
| **tags** | 标签表 | ✅ 已使用 |
| **article_tags** | 文章与标签关联表 | ✅ 已使用 |
| **comments** | 评论表 | ✅ 已使用 |
| **likes** | 点赞表 | ✅ 已使用 |

详细数据库设计请查看 [scripts/migrate.sql](scripts/migrate.sql)

---

## 开发进度

### ✅ 已完成功能
- [x] 用户认证系统（注册/登录/JWT/刷新Token）
- [x] 文章管理（CRUD/分页/搜索/筛选）
- [x] 分类和标签管理
- [x] 评论系统（支持回复/嵌套评论）
- [x] 图片上传功能
- [x] 访问量统计（Redis缓存+定时同步）
- [x] 时间轴展示
- [x] 关于页面管理
- [x] 后台管理界面API
- [x] Slug友好URL生成
- [x] 图片清理服务

### 🔄 进行中 / 待开发
- [ ] 评论通知功能
- [ ] RSS 订阅支持
- [ ] 站点地图生成（Sitemap）
- [ ] SEO 优化（Meta标签、结构化数据）
- [ ] 定时发布功能
- [ ] 文章版本管理/历史记录
- [ ] 第三方登录（GitHub、Google等）
- [ ] 全文检索引擎集成（Elasticsearch/Meilisearch）
- [ ] 文章导出（PDF/Word）
- [ ] 多语言支持（i18n）

---

## 注意事项

1. **图片上传**: 需要配置上传目录权限（`uploads/`），或使用云存储服务
2. **JWT 安全**: 请妥善保管 `JWT_SECRET`，定期更换，建议使用环境变量注入
3. **数据库备份**: 定期备份 MySQL 数据库，建议配置自动备份策略
4. **日志管理**: 日志文件存储在 `logs/` 目录，Lumberjack 会自动轮转和清理旧日志
5. **Redis 依赖**: 访问量统计功能依赖 Redis，请确保 Redis 服务正常运行
6. **定时任务**: 项目包含两个后台定时任务：
   - **访问量同步**：每5分钟将 Redis 中的访问量同步到 MySQL
   - **图片清理**：每24小时清理未使用的图片文件
7. **性能优化**: 大数据量时已通过 Redis 缓存访问量，可考虑添加更多缓存层
8. **评论审核**: 建议开启评论审核功能，防止垃圾评论（待实现）
9. **Slug 冲突**: 系统会自动生成文章 Slug，如遇冲突会追加随机后缀

---

## 常见问题

### 1. 如何修改服务端口？
编辑 `configs/config.yaml` 文件中的 `app.port` 配置项

### 2. 如何更换数据库？
修改 `configs/config.yaml` 中的数据库连接信息，确保 MySQL 服务已启动

### 3. 上传文件存储在哪里？
默认存储在 `uploads/articles/` 目录，可在配置文件中修改

### 4. 如何设置管理员？
在数据库中直接修改 `users` 表的 `role` 字段为 `admin`

### 5. JWT Token 有效期多久？
默认 24 小时，可在 `configs/config.yaml` 的 `jwt.expire_hours` 中修改

---

## 许可证

MIT License

---

## 联系方式

如有问题，请提交 Issue 或联系开发者。
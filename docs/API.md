# API 文档 - MyBlog 个人博客系统

## 基础信息

- **Base URL**: `http://localhost:8080`
- **API 版本**: v1
- **Content-Type**: `application/json`

## 统一响应格式

所有 API 返回统一的响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- `code`: 状态码，0 表示成功
- `message`: 响应消息
- `data`: 响应数据

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 1001 | 用户不存在 |
| 1002 | 用户已存在 |
| 1003 | 用户名或密码错误 |
| 1004 | 用户已被禁用 |
| 1005 | 无效的令牌 |
| 1006 | 令牌已过期 |
| 2001 | 文章不存在 |
| 2002 | 文章已存在 |
| 2003 | 分类不存在 |
| 2004 | 分类已存在 |
| 2005 | 标签不存在 |
| 2006 | 标签已存在 |
| 2007 | 评论不存在 |
| 2008 | 禁止评论 |
| 3001 | 已点赞 |
| 3002 | 未点赞 |
| 4001 | 文件上传失败 |
| 4002 | 不支持的文件格式 |
| 4003 | 文件过大 |

## 认证方式

使用 JWT Bearer Token 认证：

```
Authorization: Bearer {token}
```

## API 接口

### 健康检查

#### 检查服务状态

```http
GET /api/v1/health
```

**响应示例**:
```json
{
  "status": "ok",
  "message": "CloudQue API is running"
}
```

---

### 用户认证

#### 用户注册

```http
POST /api/v1/auth/register
```

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com",
  "nickname": "测试用户"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| password | string | 是 | 密码（6-50字符） |
| email | string | 是 | 邮箱地址 |
| nickname | string | 否 | 昵称（最多50字符） |

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

#### 用户登录

```http
POST /api/v1/auth/login
```

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### 刷新 Token

```http
POST /api/v1/auth/refresh
```

**请求参数**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 用户管理

以下接口需要认证，请在请求头中携带 Token。

#### 获取当前用户信息

```http
GET /api/v1/user/profile
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 更新用户信息

```http
PUT /api/v1/user/profile
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/avatar.jpg"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| nickname | string | 否 | 昵称 |
| avatar | string | 否 | 头像 URL |

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

#### 修改密码

```http
POST /api/v1/user/password
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码（6-50字符） |

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

#### 获取用户列表（分页）

```http
GET /api/v1/user/list?page=1&size=10
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码，从 1 开始 |
| size | int | 是 | 每页大小，1-100 |

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "nickname": "测试用户",
        "avatar": "",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10,
    "total_page": 10
  }
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 用户数据列表 |
| total | int64 | 总记录数 |
| page | int | 当前页码 |
| size | int | 每页大小 |
| total_page | int | 总页数 |

---

### 文章接口

#### 获取文章列表（分页、筛选）

```http
GET /api/v1/articles?page=1&size=10&keyword=Go&category_id=1&tag=Go&start_date=2024-01-01&end_date=2024-12-31
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码，从 1 开始 |
| size | int | 是 | 每页大小，1-100 |
| keyword | string | 否 | 标题关键词 |
| category_id | int | 否 | 分类 ID |
| tag | string | 否 | 标签名称 |
| start_date | string | 否 | 开始日期（YYYY-MM-DD） |
| end_date | string | 否 | 结束日期（YYYY-MM-DD） |

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Go 语言入门",
        "slug": "go-intro",
        "summary": "Go 语言基础教程",
        "cover_image": "https://example.com/cover.jpg",
        "view_count": 100,
        "like_count": 20,
        "comment_count": 5,
        "author": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "category": {
          "id": 1,
          "name": "技术",
          "slug": "tech"
        },
        "tags": [
          {"id": 1, "name": "Go", "slug": "go"},
          {"id": 2, "name": "编程", "slug": "programming"}
        ],
        "published_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 10,
    "total_page": 5
  }
}
```

#### 获取文章详情

```http
GET /api/v1/articles/:id
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "Go 语言入门",
    "slug": "go-intro",
    "summary": "Go 语言基础教程",
    "content": "# Go 语言入门\n\nGo 是一种开源的编程语言...",
    "html_content": "<h1>Go 语言入门</h1><p>Go 是一种开源的编程语言...</p>",
    "cover_image": "https://example.com/cover.jpg",
    "view_count": 100,
    "like_count": 20,
    "comment_count": 5,
    "author": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "avatar": ""
    },
    "category": {
      "id": 1,
      "name": "技术",
      "slug": "tech",
      "description": "技术类文章"
    },
    "tags": [
      {"id": 1, "name": "Go", "slug": "go"},
      {"id": 2, "name": "编程", "slug": "programming"}
    ],
    "published_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 创建文章（需认证）

```http
POST /api/v1/articles
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "title": "Go 语言入门",
  "slug": "go-intro",
  "summary": "Go 语言基础教程",
  "content": "# Go 语言入门\n\nGo 是一种开源的编程语言...",
  "category_id": 1,
  "tags": ["Go", "编程", "教程"],
  "cover_image": "https://example.com/cover.jpg",
  "status": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 文章标题（1-255 字符） |
| slug | string | 否 | 文章别名（URL 友好） |
| summary | string | 否 | 文章摘要 |
| content | string | 是 | 文章内容（Markdown 格式） |
| category_id | int | 是 | 分类 ID |
| tags | array | 否 | 标签名称数组 |
| cover_image | string | 否 | 封面图片 URL |
| status | int | 否 | 状态：1 发布，2 草稿，3 隐藏（默认 1） |

**响应示例**:
```json
{
  "code": 0,
  "message": "文章创建成功",
  "data": {
    "id": 1,
    "title": "Go 语言入门",
    "slug": "go-intro"
  }
}
```

#### 更新文章（需认证）

```http
PUT /api/v1/articles/:id
Authorization: Bearer {token}
```

**请求参数**: 同创建文章

#### 删除文章（需认证）

```http
DELETE /api/v1/articles/:id
Authorization: Bearer {token}
```

---

### 分类接口

#### 获取所有分类

```http
GET /api/v1/categories
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "技术",
        "slug": "tech",
        "description": "技术类文章",
        "post_count": 10,
        "sort_order": 1
      },
      {
        "id": 2,
        "name": "生活",
        "slug": "life",
        "description": "生活随笔",
        "post_count": 5,
        "sort_order": 2
      }
    ]
  }
}
```

#### 创建分类（需认证）

```http
POST /api/v1/categories
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "name": "技术",
  "slug": "tech",
  "description": "技术类文章",
  "sort_order": 1,
  "parent_id": 0
}
```

---

### 标签接口

#### 获取所有标签

```http
GET /api/v1/tags
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Go",
        "slug": "go",
        "description": "Go 语言相关",
        "post_count": 5
      },
      {
        "id": 2,
        "name": "编程",
        "slug": "programming",
        "description": "编程技术",
        "post_count": 10
      }
    ]
  }
}
```

---

### 评论接口

#### 获取文章评论列表

```http
GET /api/v1/comments/article/:id?page=1&size=20
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "article_id": 1,
        "user_id": 2,
        "parent_id": 0,
        "content": "写得很棒！学到了很多",
        "user": {
          "id": 2,
          "username": "user1",
          "nickname": "用户 1",
          "avatar": ""
        },
        "reply_count": 2,
        "like_count": 5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "size": 20,
    "total_page": 1
  }
}
```

#### 发布评论（需认证）

```http
POST /api/v1/comments
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "article_id": 1,
  "content": "写得很棒！学到了很多",
  "parent_id": 0
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章 ID |
| content | string | 是 | 评论内容（1-1000 字符） |
| parent_id | int | 否 | 父评论 ID（回复评论时使用） |

---

### 点赞接口

#### 点赞文章（需认证，每人限一次）

```http
POST /api/v1/likes/article/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "点赞成功",
  "data": {
    "like_count": 21
  }
}
```

#### 取消点赞（需认证）

```http
DELETE /api/v1/likes/article/:id
Authorization: Bearer {token}
```

#### 获取文章点赞数

```http
GET /api/v1/likes/article/:id
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "like_count": 21,
    "user_liked": true
  }
}
```

---

### 文件上传接口

#### 上传图片（需认证）

```http
POST /api/v1/upload/image
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:
- `file`: 图片文件（支持 jpg、jpeg、png、gif 格式，最大 5MB）

**响应示例**:
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "url": "https://example.com/uploads/articles/2024/01/image.jpg",
    "filename": "image.jpg",
    "size": 102400
  }
}
```

---

### 后台管理接口（需管理员权限）

#### 获取文章列表（管理端）

```http
GET /api/v1/admin/articles?page=1&size=10&status=1
Authorization: Bearer {token}
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| status | int | 否 | 文章状态筛选 |
| keyword | string | 否 | 搜索关键词 |

#### 删除评论（管理员）

```http
DELETE /api/v1/admin/comments/:id
Authorization: Bearer {token}
```

---

## 使用示例

### cURL

```bash
# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","email":"test@example.com"}'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 获取用户信息
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# 获取文章列表（分页）
curl "http://localhost:8080/api/v1/articles?page=1&size=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### JavaScript (Fetch)

```javascript
// 用户登录
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    username: 'test',
    password: '123456',
  }),
});

const data = await response.json();
const token = data.data.token;

// 获取用户信息
const profile = await fetch('http://localhost:8080/api/v1/user/profile', {
  headers: {
    'Authorization': `Bearer ${token}`,
  },
});

const profileData = await profile.json();

// 获取用户列表（分页）
const users = await fetch('http://localhost:8080/api/v1/user/list?page=1&size=10', {
  headers: {
    'Authorization': `Bearer ${token}`,
  },
});

const usersData = await users.json();
console.log(usersData.data.list); // 用户列表
console.log(usersData.data.total); // 总记录数
console.log(usersData.data.total_page); // 总页数
```

### Python (requests)

```python
import requests

# 用户登录
response = requests.post('http://localhost:8080/api/v1/auth/login', json={
    'username': 'test',
    'password': '123456',
})

data = response.json()
token = data['data']['token']

# 获取用户信息
profile = requests.get('http://localhost:8080/api/v1/user/profile', headers={
    'Authorization': f'Bearer {token}',
})

profile_data = profile.json()

# 获取用户列表（分页）
users = requests.get('http://localhost:8080/api/v1/user/list', params={
    'page': 1,
    'size': 10
}, headers={
    'Authorization': f'Bearer {token}',
})

users_data = users.json()
print(users_data['data']['list'])  # 用户列表
print(users_data['data']['total'])  # 总记录数
print(users_data['data']['total_page'])  # 总页数
```

---

## 注意事项

1. **Token 有效期**: Token 默认有效期为 24 小时
2. **密码安全**: 密码使用 bcrypt 加密存储，服务端无法查看明文密码
3. **请求频率**: 建议客户端实现请求频率限制
4. **错误处理**: 请根据错误码进行相应的错误处理
5. **时区**: 所有时间使用 UTC 时区，格式为 ISO 8601

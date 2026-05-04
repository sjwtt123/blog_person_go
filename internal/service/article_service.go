package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/utils"
	"strings"
	"time"
)

type articleService struct {
	articleRepo  repository.ArticleRepository
	tagRepo      repository.TagRepository
	categoryRepo repository.CategoryRepository
}

func NewArticleService(
	articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
	categoryRepo repository.CategoryRepository,
) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
	}
}

// ListPublic 获取公开文章列表（仅返回已发布文章）
func (s *articleService) ListPublic(req *request.ArticleListRequest) (*response.PageResponse, error) {
	// 计算分页参数
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	// 构建查询过滤器
	filter := repository.ArticleListFilter{
		Offset:     offset,
		Limit:      size,
		Keyword:    strings.TrimSpace(req.Keyword),
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		OnlyPublic: true,
	}

	// 查询文章列表和总数
	list, total, err := s.articleRepo.List(filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]response.ArticleListItemResponse, 0, len(list))
	for _, a := range list {
		items = append(items, toArticleListItem(a))
	}
	return response.NewPageResponse(items, total, page, size), nil
}

// GetDetail 获取文章详情
func (s *articleService) GetDetail(id uint) (*response.ArticleDetailResponse, error) {
	// 根据 ID 查询文章
	a, err := s.articleRepo.FindArticleByID(id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}
	// 转换为详情响应格式
	return toArticleDetail(a), nil
}

// ListAdmin 获取管理员文章列表（包含所有状态）
func (s *articleService) ListAdmin(req *request.AdminArticleListRequest) (*response.PageResponse, error) {
	// 计算分页参数
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	// 构建查询过滤器（支持按状态筛选）
	filter := repository.ArticleListFilter{
		Offset:  offset,
		Limit:   size,
		Keyword: strings.TrimSpace(req.Keyword),
		Status:  req.Status,
	}

	// 查询文章列表和总数
	list, total, err := s.articleRepo.List(filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]response.ArticleListItemResponse, 0, len(list))
	for _, a := range list {
		items = append(items, toArticleListItem(a))
	}
	return response.NewPageResponse(items, total, page, size), nil
}

// Create 创建文章
func (s *articleService) Create(authorID uint, req *request.ArticleUpsertRequest) error {
	// 校验分类是否存在
	categoryID, err := s.validateCategory(req.CategoryID)
	if err != nil {
		return err
	}

	// 规范化文章状态
	status := normalizeArticleStatus(req.Status)
	// 根据状态决定发布时间
	publishedAt := s.determinePublishedAt(status)
	// 确定文章 Slug
	slug := s.determineSlug(req.Slug, req.Title)
	// 确保标签存在
	tags, err := s.ensureTags(req.Tags)
	if err != nil {
		return err
	}

	// 构建文章实体
	a := s.buildArticle(authorID, req, categoryID, status, publishedAt, slug, tags)
	// 级联创建文章及关联
	return s.articleRepo.CreateWithCascade(a)
}

// validateCategory 校验分类是否存在
func (s *articleService) validateCategory(categoryID *uint) (*uint, error) {
	if categoryID == nil {
		return nil, nil
	}
	// 根据 ID 查询分类
	cat, err := s.categoryRepo.FindByID(*categoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "分类不存在")
	}
	return categoryID, nil
}

// determinePublishedAt 根据状态决定发布时间
func (s *articleService) determinePublishedAt(status int) *time.Time {
	if status == 1 {
		now := time.Now()
		return &now
	}
	return nil
}

// determineSlug 确定文章 Slug，未提供时从标题生成
func (s *articleService) determineSlug(slug, title string) string {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		slug = utils.Slugify(title)
	}
	return slug
}

// buildArticle 构建文章实体
func (s *articleService) buildArticle(authorID uint, req *request.ArticleUpsertRequest,
	categoryID *uint, status int, publishedAt *time.Time, slug string, tags []*entity.Tag) *entity.Article {
	return &entity.Article{
		Title:       strings.TrimSpace(req.Title),
		Slug:        slug,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		CategoryID:  categoryID,
		AuthorID:    authorID,
		Status:      status,
		PublishedAt: publishedAt,
		Tags:        tags,
	}
}

// Update 更新文章
func (s *articleService) Update(id uint, req *request.ArticleUpsertRequest) error {
	// 查询文章是否存在
	a, err := s.articleRepo.FindArticleByID(id)
	if err != nil {
		return err
	}
	if a == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}

	// 校验分类是否存在
	categoryID, err := s.validateCategory(req.CategoryID)
	if err != nil {
		return err
	}

	// 规范化文章状态
	status := normalizeArticleStatus(req.Status)
	// 更新发布时间
	s.updatePublishedAt(a, status)

	// 确定文章 Slug
	slug := s.determineSlug(req.Slug, req.Title)
	// 确保标签存在
	tags, err := s.ensureTags(req.Tags)
	if err != nil {
		return err
	}

	// 保存旧的分类和标签，用于级联更新
	oldCategoryID := a.CategoryID
	oldTags := a.Tags

	// 更新文章字段
	s.updateArticleFields(a, req, categoryID, status, slug, tags)
	// 级联更新文章及关联
	return s.articleRepo.UpdateWithCascade(a, oldCategoryID, oldTags)
}

// updatePublishedAt 更新发布时间：发布时设置，取消发布时清空
func (s *articleService) updatePublishedAt(a *entity.Article, status int) {
	if status == 1 && a.PublishedAt == nil {
		now := time.Now()
		a.PublishedAt = &now
	}
	if status != 1 {
		a.PublishedAt = nil
	}
}

// updateArticleFields 更新文章字段
func (s *articleService) updateArticleFields(a *entity.Article, req *request.ArticleUpsertRequest,
	categoryID *uint, status int, slug string, tags []*entity.Tag) {
	a.Title = strings.TrimSpace(req.Title)
	a.Slug = slug
	a.Summary = req.Summary
	a.Content = req.Content
	a.CoverImage = req.CoverImage
	a.Status = status
	a.Tags = tags
	a.CategoryID = categoryID
}

// Delete 删除文章
func (s *articleService) Delete(id uint) error {
	// 查询文章是否存在
	a, err := s.articleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if a == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}

	// 级联删除文章及关联数据
	return s.articleRepo.DeleteWithCascade(id)
}

// ensureTags 确保标签存在，不存在则创建
func (s *articleService) ensureTags(names []string) ([]*entity.Tag, error) {
	// 去重和清理空标签名
	clean := make([]string, 0, len(names))
	seen := map[string]struct{}{}
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		clean = append(clean, n)
	}

	if len(clean) == 0 {
		return nil, nil
	}

	// 查找或创建标签
	out := make([]*entity.Tag, 0, len(clean))
	for _, name := range clean {
		t, err := s.tagRepo.FindTagByName(name)
		if err != nil {
			return nil, err
		}
		if t == nil {
			// 标签不存在，创建新标签
			tag := &entity.Tag{Name: name, Slug: utils.Slugify(name)}
			if err := s.tagRepo.Create(tag); err != nil {
				return nil, err
			}
			out = append(out, tag)
			continue
		}
		out = append(out, t)
	}
	return out, nil
}

// normalizeArticleStatus 规范化文章状态：1=已发布, 2=草稿, 3=隐藏
func normalizeArticleStatus(v int) int {
	if v <= 0 {
		return 2
	} else if v > 3 {
		return 2
	}
	return v
}

// toArticleListItem 将文章实体转换为列表响应项
func toArticleListItem(a *entity.Article) response.ArticleListItemResponse {
	// 转换作者信息
	var author *response.AuthorResponse
	if a.Author != nil {
		author = &response.AuthorResponse{
			ID:       a.Author.ID,
			Username: a.Author.Username,
			Avatar:   a.Author.Avatar,
		}
	}
	// 转换分类信息
	var category *response.CategoryResponse
	if a.Category != nil {
		category = &response.CategoryResponse{
			ID:          a.Category.ID,
			Name:        a.Category.Name,
			Slug:        a.Category.Slug,
			Description: a.Category.Description,
			PostCount:   a.Category.PostCount,
			SortOrder:   a.Category.SortOrder,
			ParentID:    a.Category.ParentID,
		}
	}
	// 转换标签列表
	tags := make([]response.TagResponse, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, response.TagResponse{
			ID:          t.ID,
			Name:        t.Name,
			Slug:        t.Slug,
			Description: t.Description,
			PostCount:   t.PostCount,
		})
	}
	return response.ArticleListItemResponse{
		ID:           a.ID,
		Title:        a.Title,
		Slug:         a.Slug,
		Summary:      a.Summary,
		CoverImage:   a.CoverImage,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		CommentCount: a.CommentCount,
		Author:       author,
		Category:     category,
		Tags:         tags,
		PublishedAt:  a.PublishedAt,
		UpdatedAt:    a.UpdatedAt,
		CreatedAt:    a.CreatedAt,
		Status:       a.Status,
	}
}

// toArticleDetail 将文章实体转换为详情响应
func toArticleDetail(a *entity.Article) *response.ArticleDetailResponse {
	// 复用列表转换逻辑
	item := toArticleListItem(a)
	return &response.ArticleDetailResponse{
		ID:           item.ID,
		Title:        item.Title,
		Slug:         item.Slug,
		Summary:      item.Summary,
		Content:      a.Content,
		CoverImage:   item.CoverImage,
		ViewCount:    item.ViewCount,
		LikeCount:    item.LikeCount,
		CommentCount: item.CommentCount,
		Author:       item.Author,
		Category:     item.Category,
		Tags:         item.Tags,
		PublishedAt:  item.PublishedAt,
		UpdatedAt:    a.UpdatedAt,
		Status:       item.Status,
	}
}

// GetTimelineItems 获取时间轴文章列表
func (s *articleService) GetTimelineItems() ([]*response.TimelineItem, error) {
	// 查询时间轴文章数据
	articles, err := s.articleRepo.ListTimeline()
	if err != nil {
		return nil, err
	}

	// 转换为时间轴响应格式
	items := make([]*response.TimelineItem, 0, len(articles))
	for _, a := range articles {
		var category *response.CategoryResponse
		if a.Category != nil {
			category = &response.CategoryResponse{
				ID:   a.Category.ID,
				Name: a.Category.Name,
				Slug: a.Category.Slug,
			}
		}
		items = append(items, &response.TimelineItem{
			ID:         a.ID,
			Title:      a.Title,
			Slug:       a.Slug,
			CoverImage: a.CoverImage,
			ViewCount:  a.ViewCount,
			Category:   category,
			CreatedAt:  a.CreatedAt,
		})
	}
	return items, nil
}

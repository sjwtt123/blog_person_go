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

	"gorm.io/gorm"
)

type articleService struct {
	articleRepo  repository.ArticleRepository
	tagRepo      repository.TagRepository
	categoryRepo repository.CategoryRepository
	db           *gorm.DB
}

func NewArticleService(
	articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
	categoryRepo repository.CategoryRepository,
	db *gorm.DB,
) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
		db:           db,
	}
}

func (s *articleService) ListPublic(req *request.ArticleListRequest) (*response.PageResponse, error) {
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	filter := repository.ArticleListFilter{
		Offset:     offset,
		Limit:      size,
		Keyword:    strings.TrimSpace(req.Keyword),
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		OnlyPublic: true,
	}

	list, total, err := s.articleRepo.List(filter)
	if err != nil {
		return nil, err
	}

	items := make([]response.ArticleListItemResponse, 0, len(list))
	for _, a := range list {
		items = append(items, toArticleListItem(a))
	}
	return response.NewPageResponse(items, total, page, size), nil
}

func (s *articleService) GetDetail(id uint) (*response.ArticleDetailResponse, error) {
	a, err := s.articleRepo.FindArticleByID(id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}
	return toArticleDetail(a), nil
}

func (s *articleService) ListAdmin(req *request.AdminArticleListRequest) (*response.PageResponse, error) {
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	filter := repository.ArticleListFilter{
		Offset:  offset,
		Limit:   size,
		Keyword: strings.TrimSpace(req.Keyword),
		Status:  req.Status,
	}

	list, total, err := s.articleRepo.List(filter)
	if err != nil {
		return nil, err
	}

	items := make([]response.ArticleListItemResponse, 0, len(list))
	for _, a := range list {
		items = append(items, toArticleListItem(a))
	}
	return response.NewPageResponse(items, total, page, size), nil
}

func (s *articleService) Create(authorID uint, req *request.ArticleUpsertRequest) error {
	categoryID, err := s.validateCategory(req.CategoryID)
	if err != nil {
		return err
	}

	status := normalizeArticleStatus(req.Status)
	publishedAt := s.determinePublishedAt(status)
	slug := s.determineSlug(req.Slug, req.Title)
	tags, err := s.ensureTags(req.Tags)
	if err != nil {
		return err
	}

	a := s.buildArticle(authorID, req, categoryID, status, publishedAt, slug, tags)
	return s.createWithTransaction(a, categoryID)
}

func (s *articleService) validateCategory(categoryID *uint) (*uint, error) {
	if categoryID == nil {
		return nil, nil
	}
	cat, err := s.categoryRepo.FindByID(*categoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "分类不存在")
	}
	return categoryID, nil
}

func (s *articleService) determinePublishedAt(status int) *time.Time {
	if status == 1 {
		now := time.Now()
		return &now
	}
	return nil
}

func (s *articleService) determineSlug(slug, title string) string {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		slug = slugify(title)
	}
	return slug
}

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

func (s *articleService) createWithTransaction(a *entity.Article, categoryID *uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.CreateInTx(tx, a); err != nil {
			return err
		}

		if categoryID != nil {
			if err := s.incrementCategoryCount(tx, *categoryID); err != nil {
				return err
			}
		}

		for _, tag := range a.Tags {
			if err := s.incrementTagCount(tx, tag.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *articleService) incrementCategoryCount(tx *gorm.DB, categoryID uint) error {
	cat, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return err
	}
	if cat == nil {
		return nil
	}
	return s.categoryRepo.UpdatePostCountInTx(tx, categoryID, cat.PostCount+1)
}

func (s *articleService) incrementTagCount(tx *gorm.DB, tagID uint) error {
	tag, err := s.tagRepo.FindTagByID(tagID)
	if err != nil {
		return err
	}
	if tag == nil {
		return nil
	}
	return s.tagRepo.UpdatePostCountInTx(tx, tagID, tag.PostCount+1)
}

func (s *articleService) Update(id uint, req *request.ArticleUpsertRequest) error {
	a, err := s.articleRepo.FindArticleByID(id)
	if err != nil {
		return err
	}
	if a == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}

	categoryID, err := s.validateCategory(req.CategoryID)
	if err != nil {
		return err
	}

	status := normalizeArticleStatus(req.Status)
	s.updatePublishedAt(a, status)

	slug := s.determineSlug(req.Slug, req.Title)
	tags, err := s.ensureTags(req.Tags)
	if err != nil {
		return err
	}

	s.updateArticleFields(a, req, categoryID, status, slug, tags)
	return s.updateWithTransaction(a, categoryID)
}

func (s *articleService) updatePublishedAt(a *entity.Article, status int) {
	if status == 1 && a.PublishedAt == nil {
		now := time.Now()
		a.PublishedAt = &now
	}
	if status != 1 {
		a.PublishedAt = nil
	}
}

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

func (s *articleService) updateWithTransaction(a *entity.Article, categoryID *uint) error {
	oldArticle, err := s.articleRepo.FindArticleByID(a.ID)
	if err != nil {
		return err
	}
	oldCategoryID := oldArticle.CategoryID
	oldTags := oldArticle.Tags

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.Update(a); err != nil {
			return err
		}

		if err := s.adjustCategoryCounts(tx, oldCategoryID, categoryID); err != nil {
			return err
		}

		return s.adjustTagCounts(tx, oldTags, a.Tags)
	})
}

func (s *articleService) adjustCategoryCounts(tx *gorm.DB, oldCategoryID, newCategoryID *uint) error {
	// 如果分类未改变，无需调整
	if oldCategoryID != nil && newCategoryID != nil && *oldCategoryID == *newCategoryID {
		return nil
	}

	// 减少旧分类计数
	if oldCategoryID != nil {
		if err := s.decrementCategoryCount(tx, *oldCategoryID); err != nil {
			return err
		}
	}

	// 增加新分类计数
	if newCategoryID != nil {
		if err := s.incrementCategoryCount(tx, *newCategoryID); err != nil {
			return err
		}
	}
	return nil
}

func (s *articleService) adjustTagCounts(tx *gorm.DB, oldTags, newTags []*entity.Tag) error {
	oldTagMap := make(map[uint]bool)
	for _, t := range oldTags {
		oldTagMap[t.ID] = true
	}

	newTagMap := make(map[uint]bool)
	for _, t := range newTags {
		newTagMap[t.ID] = true
	}

	// 减少被移除的标签计数
	for _, t := range oldTags {
		if !newTagMap[t.ID] {
			if err := s.decrementTagCount(tx, t.ID); err != nil {
				return err
			}
		}
	}

	// 增加新添加的标签计数
	for _, t := range newTags {
		if !oldTagMap[t.ID] {
			if err := s.incrementTagCount(tx, t.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *articleService) decrementCategoryCount(tx *gorm.DB, categoryID uint) error {
	cat, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return err
	}
	if cat != nil && cat.PostCount > 0 {
		return s.categoryRepo.UpdatePostCountInTx(tx, categoryID, cat.PostCount-1)
	}
	return nil
}

func (s *articleService) decrementTagCount(tx *gorm.DB, tagID uint) error {
	tag, err := s.tagRepo.FindTagByID(tagID)
	if err != nil {
		return err
	}
	if tag != nil && tag.PostCount > 0 {
		return s.tagRepo.UpdatePostCountInTx(tx, tagID, tag.PostCount-1)
	}
	return nil
}

func (s *articleService) Delete(id uint) error {
	a, err := s.articleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if a == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.DeleteInTx(tx, id); err != nil {
			return err
		}

		if a.CategoryID != nil {
			cat, err := s.categoryRepo.FindByID(*a.CategoryID)
			if err != nil {
				return err
			}
			if cat != nil && cat.PostCount > 0 {
				if err := s.categoryRepo.UpdatePostCountInTx(tx, *a.CategoryID, cat.PostCount-1); err != nil {
					return err
				}
			}
		}

		for _, tag := range a.Tags {
			if err := s.decrementTagCount(tx, tag.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *articleService) ensureTags(names []string) ([]*entity.Tag, error) {
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

	out := make([]*entity.Tag, 0, len(clean))
	for _, name := range clean {
		t, err := s.tagRepo.FindTagByName(name)
		if err != nil {
			return nil, err
		}
		if t == nil {
			tag := &entity.Tag{Name: name, Slug: slugify(name)}
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

func normalizeArticleStatus(v int) int {
	if v <= 0 {
		return 2
	} else if v > 3 {
		return 1
	}
	return v
}

func toArticleListItem(a *entity.Article) response.ArticleListItemResponse {
	var author *response.AuthorResponse
	if a.Author != nil {
		author = &response.AuthorResponse{
			ID:       a.Author.ID,
			Username: a.Author.Username,
			Avatar:   a.Author.Avatar,
		}
	}
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

func toArticleDetail(a *entity.Article) *response.ArticleDetailResponse {
	item := toArticleListItem(a)
	return &response.ArticleDetailResponse{
		ID:           item.ID,
		Title:        item.Title,
		Slug:         item.Slug,
		Summary:      item.Summary,
		Content:      a.Content,
		CoverImage:   item.CoverImage,
		ViewCount:    item.ViewCount,
		CommentCount: item.CommentCount,
		Author:       item.Author,
		Category:     item.Category,
		Tags:         item.Tags,
		PublishedAt:  item.PublishedAt,
		UpdatedAt:    a.UpdatedAt,
		Status:       item.Status,
	}
}

func (s *articleService) GetTimelineItems() ([]*response.TimelineItem, error) {
	articles, err := s.articleRepo.ListTimeline()
	if err != nil {
		return nil, err
	}

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

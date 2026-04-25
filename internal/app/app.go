package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog/internal/api"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/internal/service"
	"blog/pkg/config"
	"blog/pkg/database"
	"blog/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	cfg               *config.Config
	mysqlDB           *gorm.DB
	redis             *redis.Client
	router            *api.Router
	server            *http.Server
	viewCountStopCh   chan struct{}
	uploadCleanStopCh chan struct{}
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// Initialize 初始化应用
func (a *App) Initialize() error {
	// 1. 加载配置
	if err := a.initConfig(); err != nil {
		return err
	}

	// 2. 初始化日志
	if err := a.initLogger(); err != nil {
		return err
	}

	// 3. 初始化数据库
	if err := a.initDatabase(); err != nil {
		return err
	}

	// 4. 初始化依赖
	a.initDependencies()

	// 5. 初始化路由
	a.initRouter()

	// 6. 初始化服务器
	a.initServer()

	return nil
}

// initConfig 加载配置
func (a *App) initConfig() error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	a.cfg = cfg
	return nil
}

// initLogger 初始化日志
func (a *App) initLogger() error {
	if err := logger.Init(&a.cfg.Log); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 打印启动横幅
	logger.Info("=========================================")
	logger.Info(fmt.Sprintf("欢迎使用 %s", a.cfg.App.Name))
	logger.Info(fmt.Sprintf("版本: %s", a.cfg.App.Version))
	logger.Info(fmt.Sprintf("模式: %s", a.cfg.App.Mode))
	logger.Info("配置加载成功")
	logger.Info("=========================================")

	return nil
}

// initDatabase 初始化数据库
func (a *App) initDatabase() error {
	// 初始化 MySQL
	mysqlDB, err := database.InitMySQL(&a.cfg.Database.MySQL)
	if err != nil {
		return fmt.Errorf("MySQL 初始化失败: %w", err)
	}
	a.mysqlDB = mysqlDB

	// 自动迁移数据库表
	logger.Info("开始数据库迁移...")
	if err := a.mysqlDB.AutoMigrate(
		&entity.User{},
		&entity.Article{},
		&entity.Category{},
		&entity.Tag{},
		&entity.Comment{},
	); err != nil {
		logger.Warn("数据库迁移警告", zap.Error(err))
	} else {
		logger.Info("数据库迁移完成")
	}

	// 初始化 Redis（可选）
	rs, err := database.InitRedis(&a.cfg.Database.Redis)
	if err != nil {
		logger.Warn("Redis 初始化失败，将不影响核心功能", zap.Error(err))
	}
	a.redis = rs

	return nil
}

// initDependencies 初始化依赖注入
func (a *App) initDependencies() {
	// 创建 Repository
	userRepo := repository.NewUserRepository(a.mysqlDB)
	articleRepo := repository.NewArticleRepository(a.mysqlDB, a.redis)
	categoryRepo := repository.NewCategoryRepository(a.mysqlDB)
	tagRepo := repository.NewTagRepository(a.mysqlDB)
	commentRepo := repository.NewCommentRepository(a.mysqlDB)
	viewCountRepo := repository.NewViewCountRepository(a.redis)

	// 创建 Service
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, userSvc)
	articleSvc := service.NewArticleService(articleRepo, tagRepo, categoryRepo, a.mysqlDB)
	categorySvc := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	viewCountSvc := service.NewViewCountService(articleRepo, viewCountRepo)
	commentSvc := service.NewCommentService(commentRepo, articleRepo, a.mysqlDB)
	uploadCleanSvc := service.NewUploadCleanService(articleRepo, userRepo)

	// 创建停止通道，用于优雅关闭定时任务
	a.viewCountStopCh = make(chan struct{})
	a.uploadCleanStopCh = make(chan struct{})

	// 启动定时同步（每 5 分钟同步一次 Redis 到 MySQL）
	go viewCountSvc.StartScheduledSync(a.viewCountStopCh)

	// 启动定时清理（每 24 小时清理一次未使用图片）
	go uploadCleanSvc.StartScheduledClean(a.uploadCleanStopCh)

	// 创建 Router
	a.router = api.NewRouter(userSvc, authSvc, articleSvc, categorySvc, tagService, viewCountSvc, commentSvc, uploadCleanSvc)
}

// initRouter 初始化路由
func (a *App) initRouter() {
	// 设置 Gin 模式
	gin.SetMode(a.cfg.App.Mode)
}

// initServer 初始化 HTTP 服务器
func (a *App) initServer() {
	engine := gin.New()

	// 注册路由
	a.router.Setup(engine)

	// 创建 HTTP 服务器
	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:        engine,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// Run 运行应用
func (a *App) Run() {
	// 启动 HTTP 服务器
	go func() {
		logger.Info("HTTP 服务器启动",
			zap.String("addr", a.server.Addr),
			zap.String("mode", a.cfg.App.Mode),
		)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	a.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 关闭定时同步任务，触发最后一次同步
	if a.viewCountStopCh != nil {
		close(a.viewCountStopCh)
		logger.Info("已发送停止信号给同步任务")
	}

	// 关闭图片清理任务
	if a.uploadCleanStopCh != nil {
		close(a.uploadCleanStopCh)
		logger.Info("已发送停止信号给图片清理任务")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	// 关闭数据库连接
	_ = database.CloseMySQL()
	_ = database.CloseRedis()

	// 同步日志
	_ = logger.Sync()

	logger.Info("服务器已关闭")
	logger.Info("=========================================")
}

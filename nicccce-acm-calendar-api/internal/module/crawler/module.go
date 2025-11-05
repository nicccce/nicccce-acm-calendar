package crawler

import (
	"time"

	"github.com/gin-gonic/gin"
)

type ModuleCrawler struct {
	service   *CrawlerService
	scheduler *Scheduler
	limiter   *RateLimiter
}

func (m *ModuleCrawler) GetName() string {
	return "crawler"
}

func (m *ModuleCrawler) Init() {
	// 初始化爬虫
	InitCrawlers()
	
	// 创建服务实例
	m.service = NewCrawlerService()
	m.scheduler = NewScheduler(m.service)
	m.limiter = NewRateLimiter()

	// 启动定时任务
	if err := m.scheduler.Start(); err != nil {
		panic("Failed to start scheduler: " + err.Error())
	}

	// 初始化时立即更新比赛状态
	if err := UpdateContestStatus(); err != nil {
		panic("Failed to update contest status: " + err.Error())
	}
}

func (m *ModuleCrawler) InitRouter(r *gin.RouterGroup) {
	// 比赛相关API
	contestGroup := r.Group("/contests")
	{
		contestGroup.GET("", m.GetContests)
		contestGroup.GET("/:id", m.GetContestByID)
		contestGroup.GET("/platform/:platform", m.GetContestsByPlatform)
		contestGroup.GET("/status/:status", m.GetContestsByStatus)
	}

	// 刷新相关API（需要速率限制）
	refreshGroup := r.Group("/refresh")
	refreshGroup.Use(m.limiter.GinMiddleware(5, time.Minute)) // 每分钟5次
	{
		refreshGroup.POST("", m.RefreshAllPlatforms)
		refreshGroup.POST("/:platform", m.RefreshSinglePlatform)
		refreshGroup.GET("/status", m.GetRefreshStatus)
		refreshGroup.GET("/limit", m.GetRateLimitInfo)
	}

	// 管理API
	adminGroup := r.Group("/admin/contests")
	{
		adminGroup.GET("/stats", m.GetContestStats)
		adminGroup.GET("/logs", m.GetRefreshLogs)
		adminGroup.DELETE("/:id", m.DeleteContest)
	}
}

// GetService 获取爬虫服务实例
func (m *ModuleCrawler) GetService() *CrawlerService {
	return m.service
}

// GetScheduler 获取调度器实例
func (m *ModuleCrawler) GetScheduler() *Scheduler {
	return m.scheduler
}

// GetLimiter 获取限流器实例
func (m *ModuleCrawler) GetLimiter() *RateLimiter {
	return m.limiter
}
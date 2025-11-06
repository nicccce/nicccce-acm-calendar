package crawler

import (
	"fmt"
	"strconv"
	"time"

	"nicccce-acm-calendar-api/internal/global/database"
	"nicccce-acm-calendar-api/internal/global/response"
	"nicccce-acm-calendar-api/internal/model"

	"github.com/gin-gonic/gin"
)

// GetContests 获取比赛列表
func (m *ModuleCrawler) GetContests(c *gin.Context) {
	var contests []model.Contest
	query := database.DB.Order("start_time ASC")

	// 时间过滤：只显示未来30天内的比赛
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	
	if startTime == "" {
		// 默认显示从现在开始30天内的比赛
		startTime = time.Now().Format("2006-01-02")
	}
	if endTime == "" {
		endTime = time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	}

	query = query.Where("start_time >= ? AND start_time <= ?", startTime, endTime)

	// 平台过滤
	if platform := c.Query("platform"); platform != "" {
		query = query.Where("platform = ?", platform)
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&contests).Error; err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	// 转换为DTO
	var contestDtos []model.ContestDto
	for _, contest := range contests {
		contestDto := contest.ToDto()
		contestDto.TimeRemaining = getTimeRemaining(contest.StartTime, contest.EndTime)
		contestDtos = append(contestDtos, contestDto)
	}

	response.Success(c, contestDtos)
}

// GetContestByID 根据ID获取比赛
func (m *ModuleCrawler) GetContestByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.ErrInvalidRequest)
		return
	}

	var contest model.Contest
	if err := database.DB.First(&contest, id).Error; err != nil {
		response.Fail(c, response.ErrNotFound)
		return
	}

	contestDto := contest.ToDto()
	contestDto.TimeRemaining = getTimeRemaining(contest.StartTime, contest.EndTime)

	response.Success(c, contestDto)
}

// GetContestsByPlatform 根据平台获取比赛
func (m *ModuleCrawler) GetContestsByPlatform(c *gin.Context) {
	platform := c.Param("platform")
	
	var contests []model.Contest
	if err := database.DB.
		Where("platform = ?", platform).
		Where("start_time >= ?", time.Now()).
		Order("start_time ASC").
		Find(&contests).Error; err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	var contestDtos []model.ContestDto
	for _, contest := range contests {
		contestDto := contest.ToDto()
		contestDto.TimeRemaining = getTimeRemaining(contest.StartTime, contest.EndTime)
		contestDtos = append(contestDtos, contestDto)
	}

	response.Success(c, contestDtos)
}

// GetContestsByStatus 根据状态获取比赛
func (m *ModuleCrawler) GetContestsByStatus(c *gin.Context) {
	status := c.Param("status")
	
	var contests []model.Contest
	query := database.DB.Where("status = ?", status)
	
	// 对于已结束的比赛，限制数量
	if status == "finished" {
		query = query.Order("end_time DESC").Limit(50)
	} else {
		query = query.Order("start_time ASC")
	}

	if err := query.Find(&contests).Error; err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	var contestDtos []model.ContestDto
	for _, contest := range contests {
		contestDto := contest.ToDto()
		contestDto.TimeRemaining = getTimeRemaining(contest.StartTime, contest.EndTime)
		contestDtos = append(contestDtos, contestDto)
	}

	response.Success(c, contestDtos)
}

// RefreshAllPlatforms 刷新所有平台
func (m *ModuleCrawler) RefreshAllPlatforms(c *gin.Context) {
	// 检查速率限制
	userID := uint(1) // 暂时使用固定用户ID，后续可以集成用户系统
	allowed, remaining, err := m.limiter.CheckRefreshLimit(userID, "all")
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	if !allowed {
		response.Fail(c, response.ErrServerInternal.WithTips("Rate limit exceeded", fmt.Sprintf("retry_after: %f", remaining.Seconds())))
		return
	}

	results, err := m.service.RefreshAllPlatforms(c.Request.Context())
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, results)
}

// RefreshSinglePlatform 刷新单个平台
func (m *ModuleCrawler) RefreshSinglePlatform(c *gin.Context) {
	platform := c.Param("platform")
	
	// 检查速率限制
	userID := uint(1)
	allowed, remaining, err := m.limiter.CheckRefreshLimit(userID, platform)
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	if !allowed {
		response.Fail(c, response.ErrServerInternal.WithTips("Rate limit exceeded", fmt.Sprintf("retry_after: %f", remaining.Seconds())))
		return
	}

	result, err := m.service.RefreshSinglePlatform(c.Request.Context(), platform)
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, result)
}

// GetRefreshStatus 获取刷新状态
func (m *ModuleCrawler) GetRefreshStatus(c *gin.Context) {
	logs, err := m.service.GetRecentRefreshLogs(10)
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, logs)
}

// GetRateLimitInfo 获取速率限制信息
func (m *ModuleCrawler) GetRateLimitInfo(c *gin.Context) {
	userID := uint(1)
	platform := c.Query("platform")
	if platform == "" {
		platform = "all"
	}

	current, limit, window, err := m.limiter.GetRefreshRateLimitInfo(userID, platform)
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, gin.H{
		"current": current,
		"limit":   limit,
		"window":  window.String(),
		"platform": platform,
	})
}

// GetContestStats 获取比赛统计信息
func (m *ModuleCrawler) GetContestStats(c *gin.Context) {
	var stats []struct {
		Platform string `json:"platform"`
		Status   string `json:"status"`
		Count    int    `json:"count"`
	}

	if err := database.DB.Model(&model.Contest{}).
		Select("platform, status, COUNT(*) as count").
		Where("start_time >= ?", time.Now().AddDate(0, -3, 0)).
		Group("platform, status").
		Order("platform, status").
		Find(&stats).Error; err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, stats)
}

// GetRefreshLogs 获取刷新日志
func (m *ModuleCrawler) GetRefreshLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	logs, err := m.service.GetRecentRefreshLogs(limit)
	if err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, logs)
}

// DeleteContest 删除比赛
func (m *ModuleCrawler) DeleteContest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.ErrInvalidRequest)
		return
	}

	if err := database.DB.Delete(&model.Contest{}, id).Error; err != nil {
		response.Fail(c, response.ErrServerInternal)
		return
	}

	response.Success(c, gin.H{"message": "Contest deleted successfully"})
}

// getTimeRemaining 计算剩余时间
func getTimeRemaining(startTime, endTime time.Time) string {
	now := time.Now()
	
	if now.Before(startTime) {
		// 比赛未开始
		duration := startTime.Sub(now)
		if duration.Hours() > 24 {
			return fmt.Sprintf("%.0f天后开始", duration.Hours()/24)
		}
		return fmt.Sprintf("%.0f小时后开始", duration.Hours())
	} else if now.After(endTime) {
		// 比赛已结束
		return "已结束"
	} else {
		// 比赛进行中
		duration := endTime.Sub(now)
		if duration.Hours() > 24 {
			return fmt.Sprintf("%.0f天后结束", duration.Hours()/24)
		}
		return fmt.Sprintf("%.0f小时后结束", duration.Hours())
	}
}
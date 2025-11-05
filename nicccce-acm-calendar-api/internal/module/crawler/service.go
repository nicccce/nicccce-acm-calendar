package crawler

import (
	"context"
	"fmt"
	"nicccce-acm-calendar-api/internal/global/database"
	"nicccce-acm-calendar-api/internal/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

type CrawlerService struct {
	mu sync.RWMutex
}

func NewCrawlerService() *CrawlerService {
	return &CrawlerService{}
}

// RefreshAllPlatforms 刷新所有平台的比赛数据
func (s *CrawlerService) RefreshAllPlatforms(ctx context.Context) (map[string]*RefreshResult, error) {
	results := make(map[string]*RefreshResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, crawler := range GetAllCrawlers() {
		wg.Add(1)
		go func(crawlerName string, c Crawler) {
			defer wg.Done()

			startTime := time.Now()
			result := &RefreshResult{
				Platform: crawlerName,
				StartTime: startTime,
			}

			defer func() {
				result.Duration = time.Since(startTime).Milliseconds()
				mu.Lock()
				results[crawlerName] = result
				mu.Unlock()
				
				// 记录刷新日志
				s.logRefreshResult(result)
			}()

			contests, err := c.Crawl(ctx)
			if err != nil {
				result.Status = "failed"
				result.Message = err.Error()
				return
			}

			// 保存到数据库
			newCount, updatedCount, err := s.saveContests(ctx, contests, crawlerName)
			if err != nil {
				result.Status = "failed"
				result.Message = err.Error()
				return
			}

			result.Status = "success"
			result.NewCount = newCount
			result.UpdatedCount = updatedCount
			result.Message = fmt.Sprintf("成功获取%d场比赛，新增%d场，更新%d场", len(contests), newCount, updatedCount)
		}(name, crawler)
	}

	wg.Wait()
	return results, nil
}

// RefreshSinglePlatform 刷新单个平台的比赛数据
func (s *CrawlerService) RefreshSinglePlatform(ctx context.Context, platform string) (*RefreshResult, error) {
	crawler, exists := GetCrawler(platform)
	if !exists {
		return nil, fmt.Errorf("crawler for platform %s not found", platform)
	}

	startTime := time.Now()
	result := &RefreshResult{
		Platform:  platform,
		StartTime: startTime,
	}

	defer func() {
		result.Duration = time.Since(startTime).Milliseconds()
		s.logRefreshResult(result)
	}()

	contests, err := crawler.Crawl(ctx)
	if err != nil {
		result.Status = "failed"
		result.Message = err.Error()
		return result, err
	}

	newCount, updatedCount, err := s.saveContests(ctx, contests, platform)
	if err != nil {
		result.Status = "failed"
		result.Message = err.Error()
		return result, err
	}

	result.Status = "success"
	result.NewCount = newCount
	result.UpdatedCount = updatedCount
	result.Message = fmt.Sprintf("成功获取%d场比赛，新增%d场，更新%d场", len(contests), newCount, updatedCount)

	return result, nil
}

// saveContests 保存比赛数据到数据库
func (s *CrawlerService) saveContests(ctx context.Context, contests []*model.Contest, platform string) (int, int, error) {
	var newCount, updatedCount int

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, contest := range contests {
			contest.Platform = platform

			var existingContest model.Contest
			result := tx.Where("source_id = ?", contest.SourceID).First(&existingContest)

			if result.Error == gorm.ErrRecordNotFound {
				// 新比赛
				if err := tx.Create(contest).Error; err != nil {
					return err
				}
				newCount++
			} else if result.Error != nil {
				return result.Error
			} else {
				// 更新现有比赛
				existingContest.Name = contest.Name
				existingContest.StartTime = contest.StartTime
				existingContest.EndTime = contest.EndTime
				existingContest.DurationSeconds = contest.DurationSeconds
				existingContest.ContestURL = contest.ContestURL
				existingContest.Status = contest.Status
				existingContest.LastUpdated = contest.LastUpdated

				if err := tx.Save(&existingContest).Error; err != nil {
					return err
				}
				updatedCount++
			}
		}
		return nil
	})

	if err != nil {
		return 0, 0, err
	}

	return newCount, updatedCount, nil
}

// logRefreshResult 记录刷新结果到数据库
func (s *CrawlerService) logRefreshResult(result *RefreshResult) {
	log := &model.ContestRefreshLog{
		Platform:    result.Platform,
		Status:      result.Status,
		Message:     result.Message,
		NewCount:    result.NewCount,
		UpdatedCount: result.UpdatedCount,
		Duration:    result.Duration,
	}

	if err := database.DB.Create(log).Error; err != nil {
		fmt.Printf("Failed to log refresh result: %v\n", err)
	}
}

// GetRecentRefreshLogs 获取最近的刷新日志
func (s *CrawlerService) GetRecentRefreshLogs(limit int) ([]model.ContestRefreshLog, error) {
	var logs []model.ContestRefreshLog
	err := database.DB.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

type RefreshResult struct {
	Platform     string
	Status       string
	Message      string
	NewCount     int
	UpdatedCount int
	Duration     int64
	StartTime    time.Time
}
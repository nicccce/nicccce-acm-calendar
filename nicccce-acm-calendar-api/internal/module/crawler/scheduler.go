package crawler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron          *cron.Cron
	crawlerService *CrawlerService
	mu            sync.RWMutex
	jobs          map[string]cron.EntryID
}

func NewScheduler(crawlerService *CrawlerService) *Scheduler {
	return &Scheduler{
		cron:          cron.New(cron.WithSeconds()),
		crawlerService: crawlerService,
		jobs:          make(map[string]cron.EntryID),
	}
}

// Start 启动定时任务调度器
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 添加默认的定时任务
	// 每天凌晨2点刷新所有平台
	if _, err := s.cron.AddFunc("0 0 2 * * *", s.refreshAllPlatformsJob); err != nil {
		return fmt.Errorf("failed to add daily refresh job: %w", err)
	}

	// 每小时更新比赛状态
	if _, err := s.cron.AddFunc("0 0 * * * *", s.updateContestStatusJob); err != nil {
		return fmt.Errorf("failed to add status update job: %w", err)
	}

	s.cron.Start()
	return nil
}

// Stop 停止定时任务调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cron != nil {
		s.cron.Stop()
	}
}

// AddPlatformRefreshJob 添加特定平台的定时刷新任务
func (s *Scheduler) AddPlatformRefreshJob(platform, schedule string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	crawler, exists := GetCrawler(platform)
	if !exists {
		return fmt.Errorf("crawler for platform %s not found", platform)
	}

	job := func() {
		ctx := context.Background()
		s.crawlerService.RefreshSinglePlatform(ctx, platform)
	}

	entryID, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return fmt.Errorf("failed to add refresh job for %s: %w", platform, err)
	}

	s.jobs[platform] = entryID
	return nil
}

// RemovePlatformRefreshJob 移除特定平台的定时刷新任务
func (s *Scheduler) RemovePlatformRefreshJob(platform string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.jobs[platform]
	if !exists {
		return fmt.Errorf("no job found for platform %s", platform)
	}

	s.cron.Remove(entryID)
	delete(s.jobs, platform)
	return nil
}

// GetScheduledJobs 获取所有定时任务信息
func (s *Scheduler) GetScheduledJobs() []JobInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var jobs []JobInfo
	entries := s.cron.Entries()

	for _, entry := range entries {
		jobInfo := JobInfo{
			ID:       int(entry.ID),
			Schedule: entry.Schedule.String(),
			Next:     entry.Next,
			Prev:     entry.Prev,
		}
		jobs = append(jobs, jobInfo)
	}

	return jobs
}

// refreshAllPlatformsJob 刷新所有平台的定时任务
func (s *Scheduler) refreshAllPlatformsJob() {
	ctx := context.Background()
	results, err := s.crawlerService.RefreshAllPlatforms(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh all platforms: %v\n", err)
		return
	}

	for platform, result := range results {
		fmt.Printf("Platform %s refresh result: %s\n", platform, result.Message)
	}
}

// updateContestStatusJob 更新比赛状态的定时任务
func (s *Scheduler) updateContestStatusJob() {
	if err := UpdateContestStatus(); err != nil {
		fmt.Printf("Failed to update contest status: %v\n", err)
		return
	}
	fmt.Println("Contest status updated successfully")
}

// ManualRefresh 手动触发刷新
func (s *Scheduler) ManualRefresh(platform string) (*RefreshResult, error) {
	ctx := context.Background()
	
	if platform == "all" {
		results, err := s.crawlerService.RefreshAllPlatforms(ctx)
		if err != nil {
			return nil, err
		}
		// 返回第一个平台的结果（通常用于显示）
		for _, result := range results {
			return result, nil
		}
		return nil, fmt.Errorf("no platforms refreshed")
	} else {
		return s.crawlerService.RefreshSinglePlatform(ctx, platform)
	}
}

type JobInfo struct {
	ID       int
	Schedule string
	Next     time.Time
	Prev     time.Time
}
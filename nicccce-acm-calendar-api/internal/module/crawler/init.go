package crawler

import (
	"context"
	"nicccce-acm-calendar-api/internal/global/database"
	"nicccce-acm-calendar-api/internal/model"
	"time"
)

type Crawler interface {
	Name() string
	Crawl(ctx context.Context) ([]*model.Contest, error)
}

var crawlers = make(map[string]Crawler)

func RegisterCrawler(c Crawler) {
	crawlers[c.Name()] = c
}

func GetCrawler(name string) (Crawler, bool) {
	c, exists := crawlers[name]
	return c, exists
}

func GetAllCrawlers() map[string]Crawler {
	return crawlers
}

func InitCrawlers() {
	// 注册所有爬虫
	RegisterCrawler(&CodeforcesCrawler{})
	RegisterCrawler(&AtCoderCrawler{})
	RegisterCrawler(&LeetCodeCrawler{})
	RegisterCrawler(&NowCoderCrawler{})
	RegisterCrawler(&LuoguCrawler{})
}

// UpdateContestStatus 更新比赛状态
func UpdateContestStatus() error {
	now := time.Now()
	
	// 更新进行中的比赛
	if err := database.DB.Model(&model.Contest{}).
		Where("start_time <= ? AND end_time >= ?", now, now).
		Update("status", "running").Error; err != nil {
		return err
	}
	
	// 更新已结束的比赛
	if err := database.DB.Model(&model.Contest{}).
		Where("end_time < ?", now).
		Update("status", "finished").Error; err != nil {
		return err
	}
	
	// 更新即将开始的比赛
	if err := database.DB.Model(&model.Contest{}).
		Where("start_time > ?", now).
		Update("status", "upcoming").Error; err != nil {
		return err
	}
	
	return nil
}
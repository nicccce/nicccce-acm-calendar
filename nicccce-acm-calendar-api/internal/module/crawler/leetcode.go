package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"nicccce-acm-calendar-api/internal/model"
	"time"

	"github.com/go-resty/resty/v2"
)

type LeetCodeCrawler struct {
	client *resty.Client
}

func (c *LeetCodeCrawler) Name() string {
	return "LeetCode"
}

func (c *LeetCodeCrawler) Crawl(ctx context.Context) ([]*model.Contest, error) {
	if c.client == nil {
		c.client = resty.New()
	}

	// 使用第三方API获取LeetCode比赛信息
	url := "https://algcontest.rainng.com/contests"

	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch LeetCode contests: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("LeetCode API returned status code: %d", resp.StatusCode())
	}

	var contestsData []struct {
		Oj              string `json:"oj"`
		Name            string `json:"name"`
		Link            string `json:"link"`
		StartTimeStamp  int64  `json:"startTimeStamp"`
		EndTimeStamp    int64  `json:"endTimeStamp"`
	}

	if err := json.Unmarshal(resp.Body(), &contestsData); err != nil {
		return nil, fmt.Errorf("failed to parse LeetCode contests: %w", err)
	}

	var contests []*model.Contest
	now := time.Now()

	for _, contestData := range contestsData {
		// 只处理LeetCode比赛
		if contestData.Oj != "LeetCode" {
			continue
		}

		startTime := time.Unix(contestData.StartTimeStamp, 0)
		endTime := time.Unix(contestData.EndTimeStamp, 0)
		duration := contestData.EndTimeStamp - contestData.StartTimeStamp

		// 跳过已经结束很久的比赛
		if endTime.Before(now.AddDate(0, -3, 0)) {
			continue
		}

		// 确定比赛状态
		status := "upcoming"
		if now.After(startTime) && now.Before(endTime) {
			status = "running"
		} else if now.After(endTime) {
			status = "finished"
		}

		contests = append(contests, &model.Contest{
			Name:            contestData.Name,
			Platform:        "LeetCode",
			StartTime:       startTime,
			EndTime:         endTime,
			DurationSeconds: duration,
			ContestURL:      contestData.Link,
			Status:          status,
			SourceID:        fmt.Sprintf("leetcode-%s", contestData.Name),
			LastUpdated:     now,
		})
	}

	return contests, nil
}

// 备用方法：直接解析LeetCode官网（如果需要的话）
func (c *LeetCodeCrawler) crawlFromWebsite(ctx context.Context) ([]*model.Contest, error) {
	// LeetCode官网比赛页面
	url := "https://leetcode.com/contest/"

	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch LeetCode website: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("LeetCode website returned status code: %d", resp.StatusCode())
	}

	// LeetCode官网需要解析JavaScript渲染的内容，比较复杂
	// 这里暂时使用第三方API，后续可以考虑使用headless browser

	return nil, fmt.Errorf("direct website crawling not implemented yet")
}
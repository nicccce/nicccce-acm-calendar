package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"nicccce-acm-calendar-api/internal/model"
	"time"

	"github.com/go-resty/resty/v2"
)

type LuoguCrawler struct {
	client *resty.Client
}

func (c *LuoguCrawler) Name() string {
	return "Luogu"
}

func (c *LuoguCrawler) Crawl(ctx context.Context) ([]*model.Contest, error) {
	if c.client == nil {
		c.client = resty.New()
	}

	// 洛谷API接口（需要设置User-Agent）
	url := "https://www.luogu.com.cn/contest/list?page=1&_contentOnly=1"

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch Luogu contests: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Luogu returned status code: %d", resp.StatusCode())
	}

	var response struct {
		CurrentData struct {
			Contests struct {
				Result []struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					StartTime int64  `json:"startTime"`
					EndTime   int64  `json:"endTime"`
					Status    string `json:"status"`
				} `json:"result"`
			} `json:"contests"`
		} `json:"currentData"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to parse Luogu response: %w", err)
	}

	var contests []*model.Contest
	now := time.Now()

	for _, contestData := range response.CurrentData.Contests.Result {
		startTime := time.Unix(contestData.StartTime, 0)
		endTime := time.Unix(contestData.EndTime, 0)
		duration := contestData.EndTime - contestData.StartTime

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
			Platform:        "洛谷",
			StartTime:       startTime,
			EndTime:         endTime,
			DurationSeconds: duration,
			ContestURL:      fmt.Sprintf("https://www.luogu.com.cn/contest/%d", contestData.ID),
			Status:          status,
			SourceID:        fmt.Sprintf("luogu-%d", contestData.ID),
			LastUpdated:     now,
		})
	}

	return contests, nil
}

// 备用方法：爬取多页比赛
func (c *LuoguCrawler) crawlMultiplePages(ctx context.Context) ([]*model.Contest, error) {
	var allContests []*model.Contest

	// 爬取前3页比赛
	for page := 1; page <= 3; page++ {
		url := fmt.Sprintf("https://www.luogu.com.cn/contest/list?page=%d&_contentOnly=1", page)

		resp, err := c.client.R().
			SetContext(ctx).
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36").
			Get(url)

		if err != nil {
			fmt.Printf("Failed to fetch Luogu page %d: %v\n", page, err)
			continue
		}

		if resp.StatusCode() != 200 {
			fmt.Printf("Luogu page %d returned status code: %d\n", page, resp.StatusCode())
			continue
		}

		var response struct {
			CurrentData struct {
				Contests struct {
					Result []struct {
						ID        int    `json:"id"`
						Name      string `json:"name"`
						StartTime int64  `json:"startTime"`
						EndTime   int64  `json:"endTime"`
					} `json:"result"`
				} `json:"contests"`
			} `json:"currentData"`
		}

		if err := json.Unmarshal(resp.Body(), &response); err != nil {
			fmt.Printf("Failed to parse Luogu page %d: %v\n", page, err)
			continue
		}

		now := time.Now()
		for _, contestData := range response.CurrentData.Contests.Result {
			startTime := time.Unix(contestData.StartTime, 0)
			endTime := time.Unix(contestData.EndTime, 0)
			duration := contestData.EndTime - contestData.StartTime

			// 跳过已经结束很久的比赛
			if endTime.Before(now.AddDate(0, -3, 0)) {
				continue
			}

			status := "upcoming"
			if now.After(startTime) && now.Before(endTime) {
				status = "running"
			} else if now.After(endTime) {
				status = "finished"
			}

			allContests = append(allContests, &model.Contest{
				Name:            contestData.Name,
				Platform:        "洛谷",
				StartTime:       startTime,
				EndTime:         endTime,
				DurationSeconds: duration,
				ContestURL:      fmt.Sprintf("https://www.luogu.com.cn/contest/%d", contestData.ID),
				Status:          status,
				SourceID:        fmt.Sprintf("luogu-%d", contestData.ID),
				LastUpdated:     now,
			})
		}
	}

	return allContests, nil
}
package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"nicccce-acm-calendar-api/internal/model"
	"time"

	"github.com/go-resty/resty/v2"
)

type CodeforcesCrawler struct {
	client *resty.Client
}

func (c *CodeforcesCrawler) Name() string {
	return "Codeforces"
}

func (c *CodeforcesCrawler) Crawl(ctx context.Context) ([]*model.Contest, error) {
	if c.client == nil {
		c.client = resty.New()
	}

	// Codeforces API endpoint
	url := "https://codeforces.com/api/contest.list"

	var response struct {
		Status  string `json:"status"`
		Comment string `json:"comment,omitempty"`
		Result  []struct {
			ID                  int    `json:"id"`
			Name                string `json:"name"`
			Type                string `json:"type"`
			Phase               string `json:"phase"`
			Frozen             bool   `json:"frozen"`
			DurationSeconds    int64  `json:"durationSeconds"`
			StartTimeSeconds   int64  `json:"startTimeSeconds,omitempty"`
			RelativeTimeSeconds int64  `json:"relativeTimeSeconds,omitempty"`
		} `json:"result"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch Codeforces contests: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Codeforces API returned status code: %d", resp.StatusCode())
	}

	if response.Status != "OK" {
		return nil, fmt.Errorf("Codeforces API error: %s", response.Comment)
	}

	var contests []*model.Contest
	now := time.Now()

	for _, contest := range response.Result {
		// 只处理即将开始和正在进行的比赛
		if contest.Phase != "BEFORE" && contest.Phase != "CODING" {
			continue
		}

		// 跳过已经结束很久的比赛
		if contest.StartTimeSeconds > 0 && time.Unix(contest.StartTimeSeconds, 0).Before(now.AddDate(0, -6, 0)) {
			continue
		}

		startTime := time.Unix(contest.StartTimeSeconds, 0)
		endTime := startTime.Add(time.Duration(contest.DurationSeconds) * time.Second)

		// 确定比赛状态
		status := "upcoming"
		if now.After(startTime) && now.Before(endTime) {
			status = "running"
		} else if now.After(endTime) {
			status = "finished"
		}

		contests = append(contests, &model.Contest{
			Name:            contest.Name,
			Platform:        "Codeforces",
			StartTime:       startTime,
			EndTime:         endTime,
			DurationSeconds: contest.DurationSeconds,
			ContestURL:      fmt.Sprintf("https://codeforces.com/contest/%d", contest.ID),
			Status:          status,
			SourceID:        fmt.Sprintf("codeforces-%d", contest.ID),
			LastUpdated:     now,
		})
	}

	return contests, nil
}

// CodeforcesResponse 用于解析API响应
type CodeforcesResponse struct {
	Status  string          `json:"status"`
	Comment string          `json:"comment,omitempty"`
	Result  json.RawMessage `json:"result"`
}
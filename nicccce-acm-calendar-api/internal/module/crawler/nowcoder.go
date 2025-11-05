package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"nicccce-acm-calendar-api/internal/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type NowCoderCrawler struct {
	client *resty.Client
}

func (c *NowCoderCrawler) Name() string {
	return "NowCoder"
}

func (c *NowCoderCrawler) Crawl(ctx context.Context) ([]*model.Contest, error) {
	if c.client == nil {
		c.client = resty.New()
	}

	// 牛客竞赛页面（学校竞赛）
	url := "https://ac.nowcoder.com/acm/contest/vip-index?topCategoryFilter=14"

	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch NowCoder contests: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("NowCoder returned status code: %d", resp.StatusCode())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var contests []*model.Contest
	now := time.Now()

	// 解析比赛列表
	doc.Find(".platform-item.js-item").Each(func(i int, s *goquery.Selection) {
		contest, err := c.parseContestItem(s, now)
		if err != nil {
			fmt.Printf("Failed to parse NowCoder contest item: %v\n", err)
			return
		}
		if contest != nil {
			contests = append(contests, contest)
		}
	})

	return contests, nil
}

func (c *NowCoderCrawler) parseContestItem(s *goquery.Selection, now time.Time) (*model.Contest, error) {
	// 获取JSON数据
	dataJSON, exists := s.Attr("data-json")
	if !exists {
		return nil, fmt.Errorf("data-json attribute not found")
	}

	// HTML解码并解析JSON
	decodedJSON := html.UnescapeString(dataJSON)
	var contestData struct {
		ContestName      string `json:"contestName"`
		ContestStartTime int64  `json:"contestStartTime"`
		ContestEndTime   int64  `json:"contestEndTime"`
		ContestDuration  int64  `json:"contestDuration"`
	}

	if err := json.Unmarshal([]byte(decodedJSON), &contestData); err != nil {
		return nil, fmt.Errorf("failed to parse contest JSON: %w", err)
	}

	// 获取比赛链接
	link := s.Find("a").First()
	if link.Length() == 0 {
		return nil, fmt.Errorf("contest link not found")
	}

	contestURL, exists := link.Attr("href")
	if !exists {
		return nil, fmt.Errorf("contest href not found")
	}
	contestURL = "https://ac.nowcoder.com" + contestURL

	// 转换时间（毫秒到秒）
	startTime := time.Unix(contestData.ContestStartTime/1000, 0)
	endTime := time.Unix(contestData.ContestEndTime/1000, 0)
	duration := contestData.ContestDuration / 1000

	// 跳过已经结束很久的比赛
	if endTime.Before(now.AddDate(0, -3, 0)) {
		return nil, nil
	}

	// 确定比赛状态
	status := "upcoming"
	if now.After(startTime) && now.Before(endTime) {
		status = "running"
	} else if now.After(endTime) {
		status = "finished"
	}

	return &model.Contest{
		Name:            contestData.ContestName,
		Platform:        "牛客",
		StartTime:       startTime,
		EndTime:         endTime,
		DurationSeconds: duration,
		ContestURL:      contestURL,
		Status:          status,
		SourceID:        fmt.Sprintf("nowcoder-%s", strings.ToLower(strings.ReplaceAll(contestData.ContestName, " ", "-"))),
		LastUpdated:     now,
	}, nil
}

// 爬取牛客主机竞赛
func (c *NowCoderCrawler) crawlHostContests(ctx context.Context) ([]*model.Contest, error) {
	// 主机竞赛页面
	url := "https://ac.nowcoder.com/acm/contest/vip-index?topCategoryFilter=13"

	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch NowCoder host contests: %w", err)
	}

	// 解析逻辑与学校竞赛类似
	// 这里可以复用 parseContestItem 方法

	return nil, fmt.Errorf("host contests crawling not fully implemented yet")
}
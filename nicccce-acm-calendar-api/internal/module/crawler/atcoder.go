package crawler

import (
	"context"
	"fmt"
	"nicccce-acm-calendar-api/internal/model"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type AtCoderCrawler struct {
	client *resty.Client
}

func (c *AtCoderCrawler) Name() string {
	return "AtCoder"
}

func (c *AtCoderCrawler) Crawl(ctx context.Context) ([]*model.Contest, error) {
	if c.client == nil {
		c.client = resty.New()
	}

	url := "https://atcoder.jp/contests/"

	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch AtCoder contests: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("AtCoder returned status code: %d", resp.StatusCode())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var contests []*model.Contest
	now := time.Now()

	// 解析即将开始的比赛
	doc.Find("#contest-table-upcoming .table-default tbody tr").Each(func(i int, s *goquery.Selection) {
		contest, err := c.parseContestRow(s, now)
		if err != nil {
			fmt.Printf("Failed to parse AtCoder contest row: %v\n", err)
			return
		}
		contests = append(contests, contest)
	})

	// 解析正在进行的比赛
	doc.Find("#contest-table-active .table-default tbody tr").Each(func(i int, s *goquery.Selection) {
		contest, err := c.parseContestRow(s, now)
		if err != nil {
			fmt.Printf("Failed to parse AtCoder contest row: %v\n", err)
			return
		}
		if contest != nil {
			contest.Status = "running"
			contests = append(contests, contest)
		}
	})

	return contests, nil
}

func (c *AtCoderCrawler) parseContestRow(s *goquery.Selection, now time.Time) (*model.Contest, error) {
	// 获取比赛时间
	timeText := s.Find("time").First().Text()
	startTime, err := parseAtCoderTime(timeText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	// 获取比赛链接和名称
	link := s.Find("td:nth-child(2) a").First()
	if link.Length() == 0 {
		return nil, fmt.Errorf("contest link not found")
	}

	contestURL, exists := link.Attr("href")
	if !exists {
		return nil, fmt.Errorf("contest href not found")
	}
	contestURL = "https://atcoder.jp" + contestURL

	contestName := link.Text()

	// 获取比赛时长
	durationText := s.Find("td:nth-child(3)").Text()
	duration, err := parseAtCoderDuration(durationText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	endTime := startTime.Add(time.Duration(duration) * time.Second)

	// 确定比赛状态
	status := "upcoming"
	if now.After(startTime) && now.Before(endTime) {
		status = "running"
	} else if now.After(endTime) {
		status = "finished"
	}

	return &model.Contest{
		Name:            contestName,
		Platform:        "AtCoder",
		StartTime:       startTime,
		EndTime:         endTime,
		DurationSeconds: duration,
		ContestURL:      contestURL,
		Status:          status,
		SourceID:        fmt.Sprintf("atcoder-%s", strings.ToLower(strings.ReplaceAll(contestName, " ", "-"))),
		LastUpdated:     now,
	}, nil
}

func parseAtCoderTime(timeStr string) (time.Time, error) {
	// AtCoder时间格式: 2024-01-01 21:00:00+0900
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})([+-]\d{4})`)
	matches := re.FindStringSubmatch(timeStr)
	if len(matches) < 3 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	// 解析时间（忽略时区，因为AtCoder使用日本时间）
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	// 转换为UTC
	return t.UTC(), nil
}

func parseAtCoderDuration(durationStr string) (int64, error) {
	// 格式: 1:40 或 100:00
	parts := strings.Split(durationStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse hours: %w", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse minutes: %w", err)
	}

	return int64(hours*3600 + minutes*60), nil
}
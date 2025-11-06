package model

import (
	"time"
)

type Contest struct {
	Model
	Name            string    `gorm:"size:255;not null;comment:比赛名称"`
	Platform        string    `gorm:"size:50;not null;index;comment:比赛平台"`
	StartTime       time.Time `gorm:"not null;index;comment:开始时间"`
	EndTime         time.Time `gorm:"not null;index;comment:结束时间"`
	DurationSeconds int64     `gorm:"not null;comment:持续时间(秒)"`
	ContestURL      string    `gorm:"size:500;not null;comment:比赛链接"`
	Status          string    `gorm:"size:20;default:'upcoming';index;comment:比赛状态(upcoming/running/finished)"`
	SourceID        string    `gorm:"size:100;index;comment:原始平台ID"`
	LastUpdated     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:最后更新时间"`
}

type ContestPlatform struct {
	Model
	Name           string `gorm:"size:50;not null;uniqueIndex;comment:平台名称"`
	DisplayName    string `gorm:"size:50;not null;comment:显示名称"`
	APIURL         string `gorm:"size:500;comment:API地址"`
	IsActive       bool   `gorm:"default:true;comment:是否激活"`
	UpdateInterval int    `gorm:"default:3600;comment:更新间隔(秒)"`
}

type ContestRefreshLog struct {
	Model
	Platform     string `gorm:"size:50;not null;index;comment:平台名称"`
	Status       string `gorm:"size:20;not null;comment:刷新状态(success/failed)"`
	Message      string `gorm:"size:500;comment:刷新消息"`
	NewCount     int    `gorm:"default:0;comment:新增比赛数量"`
	UpdatedCount int    `gorm:"default:0;comment:更新比赛数量"`
	Duration     int64  `gorm:"comment:耗时(毫秒)"`
}

// ContestDto 用于API返回
type ContestDto struct {
	Dto
	Name            string    `json:"name"`
	Platform        string    `json:"platform"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	DurationSeconds int64     `json:"duration_seconds"`
	ContestURL      string    `json:"contest_url"`
	Status          string    `json:"status"`
	TimeRemaining   string    `json:"time_remaining,omitempty"`
}

func (c *Contest) ToDto() ContestDto {
	return ContestDto{
		Dto: Dto{
			ID:         c.ID,
			CreateTime: c.CreateTime(),
			UpdateTime: c.UpdateTime(),
		},
		Name:            c.Name,
		Platform:        c.Platform,
		StartTime:       c.StartTime,
		EndTime:         c.EndTime,
		DurationSeconds: c.DurationSeconds,
		ContestURL:      c.ContestURL,
		Status:          c.Status,
	}
}

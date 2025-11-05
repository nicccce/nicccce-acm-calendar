package crawler

import (
	"fmt"
	"nicccce-acm-calendar-api/internal/global/redis"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	mu sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

// CheckRateLimit 检查速率限制
func (r *RateLimiter) CheckRateLimit(key string, limit int, window time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// 使用Redis sorted set实现滑动窗口限流
	// 移除过期的时间戳
	if err := redis.RDB.ZRemRangeByScore(key, "0", fmt.Sprintf("%d", windowStart)).Err(); err != nil {
		return false, err
	}

	// 获取当前窗口内的请求数量
	count, err := redis.RDB.ZCard(key).Result()
	if err != nil {
		return false, err
	}

	if count >= int64(limit) {
		return false, nil
	}

	// 添加当前请求时间戳
	if err := redis.RDB.ZAdd(key, redis.Z{
		Score:  float64(now),
		Member: now,
	}).Err(); err != nil {
		return false, err
	}

	// 设置key的过期时间
	if err := redis.RDB.Expire(key, window+time.Second).Err(); err != nil {
		return false, err
	}

	return true, nil
}

// GinMiddleware Gin中间件用于API速率限制
func (r *RateLimiter) GinMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用IP地址作为限流key
		clientIP := c.ClientIP()
		key := fmt.Sprintf("ratelimit:api:%s", clientIP)

		allowed, err := r.CheckRateLimit(key, limit, window)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(429, gin.H{
				"error": "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Please try again in %s", window.String()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckRefreshLimit 检查刷新操作的速率限制
func (r *RateLimiter) CheckRefreshLimit(userID uint, platform string) (bool, time.Duration, error) {
	key := fmt.Sprintf("ratelimit:refresh:%d:%s", userID, platform)
	limit := 10 // 每10分钟10次
	window := 10 * time.Minute

	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// 移除过期的时间戳
	if err := redis.RDB.ZRemRangeByScore(key, "0", fmt.Sprintf("%d", windowStart)).Err(); err != nil {
		return false, 0, err
	}

	// 获取当前窗口内的请求数量
	count, err := redis.RDB.ZCard(key).Result()
	if err != nil {
		return false, 0, err
	}

	if count >= int64(limit) {
		// 获取最早的时间戳来计算剩余时间
		oldest, err := redis.RDB.ZRangeWithScores(key, 0, 0).Result()
		if err != nil {
			return false, 0, err
		}

		if len(oldest) > 0 {
			oldestTime := time.Unix(int64(oldest[0].Score), 0)
			nextAvailable := oldestTime.Add(window)
			remaining := time.Until(nextAvailable)
			return false, remaining, nil
		}
		return false, window, nil
	}

	// 添加当前请求时间戳
	if err := redis.RDB.ZAdd(key, redis.Z{
		Score:  float64(now),
		Member: now,
	}).Err(); err != nil {
		return false, 0, err
	}

	// 设置key的过期时间
	if err := redis.RDB.Expire(key, window+time.Second).Err(); err != nil {
		return false, 0, err
	}

	return true, 0, nil
}

// GetRefreshRateLimitInfo 获取刷新速率限制信息
func (r *RateLimiter) GetRefreshRateLimitInfo(userID uint, platform string) (int, int, time.Duration, error) {
	key := fmt.Sprintf("ratelimit:refresh:%d:%s", userID, platform)
	limit := 10
	window := 10 * time.Minute

	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// 移除过期的时间戳
	if err := redis.RDB.ZRemRangeByScore(key, "0", fmt.Sprintf("%d", windowStart)).Err(); err != nil {
		return 0, 0, 0, err
	}

	// 获取当前计数
	count, err := redis.RDB.ZCard(key).Result()
	if err != nil {
		return 0, 0, 0, err
	}

	return int(count), limit, window, nil
}
# ACM竞赛日历API接口文档

## 1. 概述

ACM竞赛日历API提供了一个统一的接口来获取各大编程竞赛平台的比赛信息，包括Codeforces、AtCoder、LeetCode等。通过这个API，您可以轻松地获取即将到来的比赛信息，并将其集成到您的应用中。

### 1.1 基础URL

```
http://localhost:8080/admin/api/
```

### 1.2 默认响应格式

所有API响应都遵循以下JSON格式：

```json
{
  "code": 200,
  "msg": "Success",
  "data": {},
  "timestamp": 1700123456
}
```

字段说明：
- `code`: 业务状态码，200表示成功，其他值表示错误
- `msg`: 响应消息，描述操作结果
- `data`: 实际数据内容
- `timestamp`: 响应时间戳

### 1.3 状态码说明

| 状态码 | 描述 |
|-------|------|
| 200 | 成功 |
| 400 | 无效的请求 |
| 401 | 权限不足 |
| 403 | 禁止访问 |
| 404 | 目标不存在 |
| 409 | 目标已存在 |
| 500 | 服务器内部错误 |

## 2. 比赛相关接口

### 2.1 获取比赛列表

#### 接口地址
```
GET /contests
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| start_time | string | 否 | 开始时间 (格式: YYYY-MM-DD)，默认为当前日期 |
| end_time | string | 否 | 结束时间 (格式: YYYY-MM-DD)，默认为30天后的日期 |
| platform | string | 否 | 平台筛选 (codeforces, atcoder, leetcode等) |
| status | string | 否 | 状态筛选 (upcoming, running, finished) |

#### 响应数据
```json
[
  {
    "id": 1,
    "create_time": "2023-11-15T10:00:00Z",
    "update_time": "2023-11-15T10:00:00Z",
    "name": "Codeforces Round #800",
    "platform": "codeforces",
    "start_time": "2023-11-20T18:00:00Z",
    "end_time": "2023-11-20T20:00:00Z",
    "duration_seconds": 7200,
    "contest_url": "https://codeforces.com/contests/1800",
    "status": "upcoming",
    "time_remaining": "5天后开始"
  }
]
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/contests?platform=codeforces&status=upcoming"
```

### 2.2 根据ID获取比赛详情

#### 接口地址
```
GET /contests/{id}
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| id | integer | 是 | 比赛ID |

#### 响应数据
```json
{
  "id": 1,
  "create_time": "2023-11-15T10:00:00Z",
  "update_time": "2023-11-15T10:00:00Z",
  "name": "Codeforces Round #800",
  "platform": "codeforces",
  "start_time": "2023-11-20T18:00:00Z",
  "end_time": "2023-11-20T20:00:00Z",
  "duration_seconds": 7200,
  "contest_url": "https://codeforces.com/contests/1800",
  "status": "upcoming",
  "time_remaining": "5天后开始"
}
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/contests/1"
```

### 2.3 根据平台获取比赛列表

#### 接口地址
```
GET /contests/platform/{platform}
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| platform | string | 是 | 平台名称 (codeforces, atcoder, leetcode等) |

#### 响应数据
```json
[
  {
    "id": 1,
    "create_time": "2023-11-15T10:00:00Z",
    "update_time": "2023-11-15T10:00:00Z",
    "name": "Codeforces Round #800",
    "platform": "codeforces",
    "start_time": "2023-11-20T18:00:00Z",
    "end_time": "2023-11-20T20:00:00Z",
    "duration_seconds": 7200,
    "contest_url": "https://codeforces.com/contests/1800",
    "status": "upcoming",
    "time_remaining": "5天后开始"
  }
]
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/contests/platform/codeforces"
```

### 2.4 根据状态获取比赛列表

#### 接口地址
```
GET /contests/status/{status}
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| status | string | 是 | 比赛状态 (upcoming, running, finished) |

#### 响应数据
```json
[
  {
    "id": 1,
    "create_time": "2023-11-15T10:00:00Z",
    "update_time": "2023-11-15T10:00:00Z",
    "name": "Codeforces Round #800",
    "platform": "codeforces",
    "start_time": "2023-11-20T18:00:00Z",
    "end_time": "2023-11-20T20:00:00Z",
    "duration_seconds": 7200,
    "contest_url": "https://codeforces.com/contests/1800",
    "status": "upcoming",
    "time_remaining": "5天后开始"
  }
]
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/contests/status/upcoming"
```

## 3. 数据刷新接口

### 3.1 刷新所有平台数据

#### 接口地址
```
POST /refresh
```

#### 请求参数
无

#### 响应数据
```json
{
  "message": "Refresh completed successfully"
}
```

#### 示例请求
```bash
curl -X POST "http://localhost:8080/admin/api/refresh"
```

### 3.2 刷新单个平台数据

#### 接口地址
```
POST /refresh/{platform}
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| platform | string | 是 | 平台名称 (codeforces, atcoder, leetcode等) |

#### 响应数据
```json
{
  "message": "Refresh completed successfully"
}
```

#### 示例请求
```bash
curl -X POST "http://localhost:8080/admin/api/refresh/codeforces"
```

### 3.3 获取刷新状态

#### 接口地址
```
GET /refresh/status
```

#### 请求参数
无

#### 响应数据
```json
[
  {
    "id": 1,
    "create_time": "2023-11-15T10:00:00Z",
    "update_time": "2023-11-15T10:00:00Z",
    "platform": "codeforces",
    "status": "success",
    "message": "Refreshed 10 contests",
    "new_count": 5,
    "updated_count": 5,
    "duration": 1200
  }
]
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/refresh/status"
```

### 3.4 获取速率限制信息

#### 接口地址
```
GET /refresh/limit
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| platform | string | 否 | 平台名称，默认为"all" |

#### 响应数据
```json
{
  "current": 2,
  "limit": 5,
  "window": "1m0s",
  "platform": "all"
}
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/refresh/limit?platform=codeforces"
```

## 4. 管理接口

### 4.1 获取比赛统计数据

#### 接口地址
```
GET /admin/contests/stats
```

#### 请求参数
无

#### 响应数据
```json
[
  {
    "platform": "codeforces",
    "status": "upcoming",
    "count": 10
  }
]
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/admin/contests/stats"
```

### 4.2 获取刷新日志

#### 接口地址
```
GET /admin/contests/logs
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| limit | integer | 否 | 返回日志条数，默认为50 |

#### 响应数据
```json
[
  {
    "id": 1,
    "create_time": "2023-11-15T10:00:00Z",
    "update_time": "2023-11-15T10:00:00Z",
    "platform": "codeforces",
    "status": "success",
    "message": "Refreshed 10 contests",
    "new_count": 5,
    "updated_count": 5,
    "duration": 1200
  }
]
```

#### 示例请求
```bash
curl "http://localhost:8080/admin/api/admin/contests/logs?limit=10"
```

### 4.3 删除比赛

#### 接口地址
```
DELETE /admin/contests/{id}
```

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| id | integer | 是 | 比赛ID |

#### 响应数据
```json
{
  "message": "Contest deleted successfully"
}
```

#### 示例请求
```bash
curl -X DELETE "http://localhost:8080/admin/api/admin/contests/1"
```

## 5. 数据模型

### 5.1 比赛数据模型

| 字段名 | 类型 | 描述 |
|--------|------|------|
| id | integer | 比赛唯一标识符 |
| create_time | datetime | 创建时间 |
| update_time | datetime | 更新时间 |
| name | string | 比赛名称 |
| platform | string | 比赛平台 |
| start_time | datetime | 开始时间 |
| end_time | datetime | 结束时间 |
| duration_seconds | integer | 持续时间(秒) |
| contest_url | string | 比赛链接 |
| status | string | 比赛状态(upcoming/running/finished) |
| time_remaining | string | 剩余时间(仅在响应中提供) |

### 5.2 刷新日志模型

| 字段名 | 类型 | 描述 |
|--------|------|------|
| id | integer | 日志唯一标识符 |
| create_time | datetime | 创建时间 |
| update_time | datetime | 更新时间 |
| platform | string | 平台名称 |
| status | string | 刷新状态(success/failed) |
| message | string | 刷新消息 |
| new_count | integer | 新增比赛数量 |
| updated_count | integer | 更新比赛数量 |
| duration | integer | 耗时(毫秒) |

## 6. 支持的平台

目前API支持以下编程竞赛平台：

1. Codeforces
2. AtCoder
3. LeetCode
4. NowCoder
5. Luogu

## 7. 错误处理

当请求发生错误时，API会返回相应的错误码和错误信息：

```json
{
  "code": 404,
  "msg": "目标不存在",
  "timestamp": 1700123456
}
```

常见错误及解决方案：

1. `400 Invalid Request`: 请求参数不正确，请检查参数格式
2. `404 Not Found`: 请求的资源不存在，请检查URL路径
3. `500 Internal Server Error`: 服务器内部错误，请稍后再试

## 8. 速率限制

为了防止API被滥用，系统对某些接口实施了速率限制：

- `/refresh` 接口：每分钟最多5次请求
- `/refresh/{platform}` 接口：每分钟最多5次请求

当超过速率限制时，API会返回400错误和重试提示。
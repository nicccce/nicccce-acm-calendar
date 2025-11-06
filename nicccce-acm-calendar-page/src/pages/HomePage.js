import React, { useState, useEffect } from 'react';
import { Card, Typography, Button, Space, Alert, Row, Col, Tag, Tooltip } from 'antd';
import { ReloadOutlined, ClockCircleOutlined, CalendarOutlined, LinkOutlined } from '@ant-design/icons';
import apiService from '../services/api';

const { Title, Paragraph } = Typography;

const HomePage = () => {
  const [contests, setContests] = useState([]);
  const [filteredContests, setFilteredContests] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [activeFilter, setActiveFilter] = useState('all');
  const [lastUpdate, setLastUpdate] = useState(null);

  // 获取比赛数据
  const fetchContests = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiService.getContests();
      // 确保数据是数组格式
      const data = Array.isArray(response) ? response : (response.data || []);
      setContests(data);
      setFilteredContests(data);
      setLastUpdate(new Date());
    } catch (err) {
      console.error('获取比赛数据失败:', err);
      setError('无法连接到后端服务器，请确保后端服务正在运行');
    } finally {
      setLoading(false);
    }
  };

  // 刷新数据
  const refreshData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      await apiService.refreshAllPlatforms();
      await fetchContests();
    } catch (err) {
      console.error('刷新数据失败:', err);
      setError('刷新数据失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  // 筛选比赛
  const filterContests = (platform) => {
    setActiveFilter(platform);
    if (platform === 'all') {
      setFilteredContests(contests);
    } else {
      const filtered = contests.filter(contest => contest.platform === platform);
      setFilteredContests(filtered);
    }
  };

  // 格式化持续时间
  const formatDuration = (seconds) => {
    if (!seconds) return '未知';
    
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    
    if (hours > 0) {
      return `${hours}小时${minutes > 0 ? `${minutes}分钟` : ''}`;
    } else {
      return `${minutes}分钟`;
    }
  };

  // 获取比赛状态颜色
  const getStatusColor = (status) => {
    switch (status) {
      case 'upcoming': return 'orange';
      case 'running': return 'green';
      case 'finished': return 'gray';
      default: return 'blue';
    }
  };

  // 获取平台颜色
  const getPlatformColor = (platform) => {
    switch (platform) {
      case 'Codeforces': return '#4CAF50';
      case 'AtCoder': return '#FF6B6B';
      case 'LeetCode': return 'hsl(167, 100%, 54%)';
      case 'NowCoder':
      case '牛客': return '#4285F4';
      case 'Luogu':
      case '洛谷': return '#FF8C00';
      default: return '#1890ff';
    }
  };

  // 组件挂载时获取数据
  useEffect(() => {
    fetchContests();
  }, []);

  // 获取所有平台
  const getAllPlatforms = () => {
    // 确保 contests 是数组并且有数据
    if (!Array.isArray(contests) || contests.length === 0) {
      return [];
    }
    const platforms = [...new Set(contests.map(contest => contest.platform))];
    return platforms;
  };

  return (
    <div>
      <Title level={2}>🏆 ACM 竞赛日历</Title>
      <Paragraph>
        一站式获取近10日各大平台竞赛信息
      </Paragraph>
      
      {/* 控制区域 */}
      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 10 }}>
          <Space>
            <Button
              type="primary"
              icon={<ReloadOutlined />}
              onClick={refreshData}
              loading={loading}
            >
              刷新竞赛信息
            </Button>
          </Space>
          <div style={{ color: '#666', fontSize: '0.9rem' }}>
            {lastUpdate && `最后更新: ${lastUpdate.toLocaleString('zh-CN')}`}
          </div>
        </div>
        
        {/* 平台筛选器 */}
        <div style={{ marginTop: 15, display: 'flex', flexWrap: 'wrap', gap: 8 }}>
          <Button
            type={activeFilter === 'all' ? 'primary' : 'default'}
            onClick={() => filterContests('all')}
          >
            全部
          </Button>
          {getAllPlatforms().map(platform => (
            <Button
              key={platform}
              type={activeFilter === platform ? 'primary' : 'default'}
              onClick={() => filterContests(platform)}
            >
              {platform}
            </Button>
          ))}
        </div>
      </Card>
      
      {/* 错误提示 */}
      {error && (
        <Alert
          message="获取数据失败"
          description={error}
          type="error"
          showIcon
          style={{ marginBottom: 24 }}
        />
      )}
      
      {/* 比赛卡片网格 */}
      <Row gutter={[16, 16]}>
        {filteredContests.length === 0 ? (
          <Col span={24}>
            <Card>
              <div style={{ textAlign: 'center', padding: 40, color: '#666' }}>
                {loading ? '正在加载竞赛信息...' : '暂无相关竞赛信息'}
              </div>
            </Card>
          </Col>
        ) : (
          filteredContests.map(contest => (
            <Col xs={24} sm={12} md={8} lg={6} key={contest.id}>
              <Card
                size="small"
                style={{
                  height: '100%',
                  borderLeft: `4px solid ${getPlatformColor(contest.platform)}`,
                  position: 'relative'
                }}
                bodyStyle={{
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'space-between',
                  height: '100%'
                }}
              >
                {/* 平台标签 */}
                <Tag
                  color={getPlatformColor(contest.platform)}
                  style={{
                    position: 'absolute',
                    top: 8,
                    right: 8,
                    fontWeight: 'bold'
                  }}
                >
                  {contest.platform}
                </Tag>
                
                {/* 比赛名称 */}
                <div style={{
                  fontSize: '1rem',
                  fontWeight: 'bold',
                  marginBottom: 10,
                  color: '#333',
                  lineHeight: 1.3,
                  display: '-webkit-box',
                  WebkitLineClamp: 2,
                  WebkitBoxOrient: 'vertical',
                  overflow: 'hidden'
                }}>
                  {contest.name}
                </div>
                
                {/* 比赛信息 */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                  <div style={{ display: 'flex', alignItems: 'flex-start', gap: 6, color: '#666', fontSize: '0.8rem' }}>
                    <CalendarOutlined style={{ fontSize: '0.9rem', minWidth: 16 }} />
                    <span>
                      {new Date(contest.start_time).toLocaleString('zh-CN', {
                        month: 'numeric',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}
                    </span>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'flex-start', gap: 6, color: '#666', fontSize: '0.8rem' }}>
                    <ClockCircleOutlined style={{ fontSize: '0.9rem', minWidth: 16 }} />
                    <span>{formatDuration(contest.duration_seconds)}</span>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: '0.8rem' }}>
                    <div style={{
                      width: 8,
                      height: 8,
                      borderRadius: '50%',
                      backgroundColor: getStatusColor(contest.status),
                      flexShrink: 0
                    }} />
                    <span>{contest.status === 'upcoming' ? '即将开始' : contest.status === 'running' ? '进行中' : '已结束'}</span>
                  </div>
                  <div style={{
                    fontWeight: 'bold',
                    color: contest.status === 'upcoming' ? '#e74c3c' : contest.status === 'running' ? '#27ae60' : '#666',
                    fontSize: '0.8rem',
                    marginTop: 4
                  }}>
                    {contest.time_remaining || '未知'}
                  </div>
                </div>
                
                {/* 查看详情按钮 */}
                <a
                  href={contest.contest_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  style={{
                    display: 'inline-block',
                    background: 'linear-gradient(45deg, #667eea, #764ba2)',
                    color: 'white',
                    textDecoration: 'none',
                    padding: '6px 12px',
                    borderRadius: '20px',
                    marginTop: 10,
                    fontWeight: 'bold',
                    textAlign: 'center',
                    fontSize: '0.85rem',
                    width: '100%',
                    height: 40,
                    lineHeight: '28px'
                  }}
                >
                  查看详情
                </a>
              </Card>
            </Col>
          ))
        )}
      </Row>
    </div>
  );
};

export default HomePage;
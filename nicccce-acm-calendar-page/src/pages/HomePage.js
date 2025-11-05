import React, { useState, useEffect } from 'react';
import { Card, Typography, Button, Space, Alert } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import apiService from '../services/api';

const { Title, Paragraph } = Typography;

const HomePage = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // 刷新数据
  const refreshData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      await apiService.refreshAllPlatforms();
    } catch (err) {
      console.error('刷新数据失败:', err);
      setError('刷新数据失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>ACM 竞赛日历</Title>
      <Paragraph>
        欢迎使用 ACM 竞赛日历应用。在这里您可以查看即将到来的编程竞赛信息。
      </Paragraph>
      
      {/* 控制区域 */}
      <Card style={{ marginBottom: 24 }}>
        <Space>
          <Button 
            type="primary" 
            icon={<ReloadOutlined />} 
            onClick={refreshData}
            loading={loading}
          >
            刷新数据
          </Button>
        </Space>
      </Card>
      
      {/* 错误提示 */}
      {error && (
        <Alert 
          message="获取数据失败" 
          description={error} 
          type="error" 
            showIcon 
          />
        )}
    </div>
  );
};

export default HomePage;
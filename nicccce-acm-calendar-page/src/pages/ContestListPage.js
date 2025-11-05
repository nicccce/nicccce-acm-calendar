import React, { useState, useEffect } from 'react';
import { Table, Card, Button, Space, Alert } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';

const ContestListPage = () => {
  const [contests, setContests] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // 获取比赛数据
  const fetchContests = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/contests');
      const data = await response.json();
      setContests(data || []);
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
      const response = await fetch('/api/refresh', {
        method: 'POST'
      });
      await fetchContests();
    } catch (err) {
      console.error('刷新数据失败:', err);
      setError('刷新数据失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  // 组件挂载时获取数据
  useEffect(() => {
    fetchContests();
  }, []);

  return (
    <div>
      <h2>竞赛列表</h2>
      
      {/* 控制区域 */}
      <Card style={{ marginBottom: 16 }}>
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
        />)
      }
      {/* 比赛列表 */}
      <Card>
        <Table
          dataSource={contests}
          columns={[
            {
              title: '比赛名称',
              dataIndex: 'name',
              key: 'name'
            },
            {
              title: '平台',
              dataIndex: 'platform',
              key: 'platform'
            }
          ]}
          rowKey="id"
          pagination={false}
        />
      </Card>
    </div>
  );
};

export default ContestListPage;
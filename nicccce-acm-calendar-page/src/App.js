import React from 'react';
import { Layout, Menu, theme } from 'antd';
import {
  UploadOutlined,
  UserOutlined,
  VideoCameraOutlined,
} from '@ant-design/icons';
import './App.css';

const { Header, Sider, Content } = Layout;

function App() {
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  return (
    <Layout className="layout">
      <Sider trigger={null} collapsible collapsed={false}>
        <div className="demo-logo-vertical" />
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['1']}
          items={[
            {
              key: '1',
              icon: <UserOutlined />,
              label: '首页',
            },
            {
              key: '2',
              icon: <VideoCameraOutlined />,
              label: '比赛日历',
            },
            {
              key: '3',
              icon: <UploadOutlined />,
              label: '个人中心',
            },
          ]}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: 0,
            background: colorBgContainer,
          }}
        >
          <div style={{ paddingLeft: 20, fontSize: 20, fontWeight: 'bold' }}>
            ACM 竞赛日历
          </div>
        </Header>
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
          }}
        >
          <div style={{ textAlign: 'center', paddingTop: 50 }}>
            <h1>欢迎使用 ACM 竞赛日历</h1>
            <p>在这里您可以查看即将到来的编程竞赛信息</p>
          </div>
        </Content>
      </Layout>
    </Layout>
  );
}

export default App;

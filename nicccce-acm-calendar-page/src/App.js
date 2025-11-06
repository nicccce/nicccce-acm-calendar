import React from 'react';
import { Layout, theme } from 'antd';
import { Outlet } from 'react-router-dom';
import './App.css';

const { Header, Content } = Layout;

function App() {
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  return (
    <Layout className="layout">
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
          {/* 子路由将在这里渲染 */}
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}

export default App;

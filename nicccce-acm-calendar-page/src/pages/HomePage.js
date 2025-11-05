import React from 'react';
import { Card, Col, Row, Typography, Space } from 'antd';

const { Title, Paragraph } = Typography;

const HomePage = () => {
  return (
    <div>
      <Title level={2}>ACM 竞赛日历</Title>
      <Paragraph>
        欢迎使用 ACM 竞赛日历应用。在这里您可以查看即将到来的编程竞赛信息。
      </Paragraph>
      
      <Row gutter={[16, 16]}>
        <Col span={8}>
          <Card title="即将开始" bordered={false}>
            <p>查看即将开始的竞赛</p>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="热门竞赛" bordered={false}>
            <p>查看最受欢迎的竞赛</p>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="我的收藏" bordered={false}>
            <p>查看您收藏的竞赛</p>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default HomePage;
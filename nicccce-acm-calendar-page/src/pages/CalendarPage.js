import React from 'react';
import { Calendar, Badge } from 'antd';

const CalendarPage = () => {
  const getListData = (value) => {
    let listData;
    switch (value.date()) {
      case 8:
        listData = [
          { type: 'warning', content: 'Codeforces Round 789' },
        ];
        break;
      case 10:
        listData = [
          { type: 'success', content: 'AtCoder Beginner Contest 256' },
        ];
        break;
      case 15:
        listData = [
          { type: 'warning', content: 'LeetCode Weekly Contest 298' },
          { type: 'success', content: 'NowCoder Contest 123' },
        ];
        break;
      default:
    }
    return listData || [];
  };

  const dateCellRender = (value) => {
    const listData = getListData(value);
    return (
      <ul className="events">
        {listData.map((item) => (
          <li key={item.content}>
            <Badge status={item.type} text={item.content} />
          </li>
        ))}
      </ul>
    );
  };

  return (
    <div>
      <h2>竞赛日历</h2>
      <Calendar dateCellRender={dateCellRender} />
    </div>
  );
};

export default CalendarPage;
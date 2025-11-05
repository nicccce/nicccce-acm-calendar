import React from 'react';
import { Routes, Route } from 'react-router-dom';
import App from './App';
import HomePage from './pages/HomePage';
import CalendarPage from './pages/CalendarPage';
import ContestListPage from './pages/ContestListPage';

const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<App />}>
        <Route index element={<HomePage />} />
        <Route path="calendar" element={<CalendarPage />} />
        <Route path="contests" element={<ContestListPage />} />
      </Route>
    </Routes>
  );
};

export default AppRoutes;
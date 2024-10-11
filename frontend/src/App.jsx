import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from './Pages/HomePage';
import LoginPage from './Pages/LoginPage';
import CanvisPage from './Pages/CanvisPage';
import AccountPage from './Pages/AccountPage';
import Logo from './Componants/Logo';

function App() {
  return (
    <Router>
      <div className="w-full h-screen">
        <div className="absolute top-6 left-6 z-10">
          <Logo />
        </div>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage isLogin={true} />} />
          <Route path="/register" element={<LoginPage isLogin={false} />} />
          <Route path="/:randomId" element={<CanvisPage />} />
          <Route path="/account" element={<AccountPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
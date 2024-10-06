import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from './Pages/HomePage';
import LoginPage from './Pages/LoginPage';
import CanisPage from './Pages/CanvisPage'
import AccountPage from './Pages/AccountPage';

function App() {
  return (
    <Router>
      <div className="w-full h-screen">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/:randomId" element={<CanisPage />} />
          <Route path="/account" element={<AccountPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
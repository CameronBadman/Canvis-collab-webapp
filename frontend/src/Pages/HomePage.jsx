import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import CustomBackground from '../Componants/CustomBackground';
import { UserCircle, LogOut } from 'lucide-react';
import { authService } from '../services/AuthService';
import { cookieService } from '../services/CookieService';

const HomePage = () => {
  const [inputValue, setInputValue] = useState('');
  const [currentUser, setCurrentUser] = useState(null);
  const [allCookies, setAllCookies] = useState({});
  const navigate = useNavigate();

  useEffect(() => {
    // Get current user from cookies or local storage
    const user = authService.getCurrentUser();
    setCurrentUser(user);
    updateCookieInfo();
  }, []);

  const updateCookieInfo = () => {
    const cookies = document.cookie.split(';').reduce((acc, cookie) => {
      const [name, value] = cookie.trim().split('=');
      acc[name] = cookieService.getCookie(name);
      return acc;
    }, {});
    setAllCookies(cookies);
  };

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (inputValue) {
      navigate(`/${inputValue}`);
    }
  };

  const handleLogout = async () => {
    // Logout and clear cookies
    await authService.logout();
    setCurrentUser(null);
    updateCookieInfo();
    navigate('/');
  };

  const handleAccountClick = async () => {
    const user = authService.getCurrentUser();
    if (user && user.jwt_token) {
      try {
        // Verify token with the backend (this can be a simple ping endpoint or something similar)
        await axios.get('/api/auth/verify', {
          headers: { Authorization: `Bearer ${user.jwt_token}` }
        });
        navigate('/account');
      } catch (error) {
        console.error('Token verification failed:', error);
        // Token verification failed, so logout the user
        await handleLogout();
      }
    } else {
      // Not logged in, redirect to homepage
      navigate('/');
    }
  };

  return (
      <div className="relative w-full h-screen overflow-hidden">
        <CustomBackground />

        {/* Account Menu */}
        <div className="absolute top-6 right-6 z-10">
          <div className="relative group">
            <UserCircle size={64} className="text-gray-800 hover:text-gray-600 cursor-pointer" />
            <div className="absolute right-0 mt-2 w-64 bg-white rounded-md shadow-lg py-2 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300">
              {currentUser ? (
                  <>
                    <div className="block px-6 py-3 text-base font-bold text-gray-900 border-b border-gray-200">
                      {currentUser.email}
                    </div>
                    <button
                        onClick={handleAccountClick}
                        className="flex items-center w-full px-6 py-3 text-base text-gray-700 hover:bg-gray-100"
                    >
                      Account
                    </button>
                    <button
                        onClick={handleLogout}
                        className="flex items-center w-full px-6 py-3 text-base text-gray-700 hover:bg-gray-100"
                    >
                      <LogOut size={20} className="mr-3" />
                      Logout
                    </button>
                  </>
              ) : (
                  <>
                    <a href="/login" className="block px-6 py-3 text-base text-gray-700 hover:bg-gray-100">Login</a>
                    <a href="/register" className="block px-6 py-3 text-base text-gray-700 hover:bg-gray-100">Register</a>
                  </>
              )}
            </div>
          </div>
        </div>

        <div className="absolute inset-0 flex flex-col items-center justify-center p-4">
          <h1 className="text-5xl sm:text-6xl md:text-7xl font-bold mb-12 text-center text-gray-800 shadow-lg">
            Welcome
          </h1>
          <form onSubmit={handleSubmit} className="w-full max-w-md px-4">
            <input
                type="text"
                value={inputValue}
                onChange={handleInputChange}
                placeholder="Please type in your canvas code"
                className="w-full px-6 py-4 text-xl border-2 border-gray-400 rounded-lg bg-white bg-opacity-70 focus:outline-none focus:ring-2 focus:ring-gray-500 mb-6 text-gray-800 placeholder-gray-400"
            />
            <button
                type="submit"
                className="w-full bg-gray-800 text-white py-4 px-6 text-xl rounded-lg border border-transparent hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 transition duration-300"
            >
              Enter
            </button>
          </form>

          {/* Cookie Debug Information */}
          <div className="mt-8 p-4 bg-white bg-opacity-70 rounded-lg shadow-lg max-w-md w-full">
            <h2 className="text-xl font-bold mb-4">Cookie Debug Information</h2>
            <pre className="whitespace-pre-wrap break-words text-sm">
            {JSON.stringify(allCookies, null, 2)}
          </pre>
          </div>
        </div>
      </div>
  );
};

export default HomePage;

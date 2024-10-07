import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import CustomBackground from '../Componants/CustomBackground';
import { UserCircle } from 'lucide-react';

const HomePage = () => {
  const [inputValue, setInputValue] = useState('');
  const navigate = useNavigate();

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (inputValue) {
      navigate(`/${inputValue}`);
    }
  };

  return (
    <div className="relative w-full h-screen overflow-hidden">
      <CustomBackground />

      {/* Account Menu */}
      <div className="absolute top-6 right-6 z-10">
        <div className="relative group">
          <UserCircle size={48} className="text-gray-800 hover:text-gray-600 cursor-pointer" />
          <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300">
            <a href="/login" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Login</a>
            <a href="/register" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Register</a>
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
      </div>
    </div>
  );
};

export default HomePage;
import React, { useState } from 'react';
import CustomBackground from '../Componants/CustomBackground';

const HomePage = () => {
  const [inputValue, setInputValue] = useState('');

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (inputValue) {
      console.log(`Navigating to /${inputValue}`);
    }
  };

  return (
    <div className="relative w-full h-screen overflow-hidden">
      <CustomBackground />
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
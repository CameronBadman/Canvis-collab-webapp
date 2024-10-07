import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import CustomBackground from '../Componants/CustomBackground';
import { UserCircle } from 'lucide-react';

const LoginPage = ({ isLogin }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    // Handle login or register logic here
    console.log(isLogin ? 'Logging in' : 'Registering', email, password);
    // After successful login/register, navigate to home
    navigate('/');
  };

  return (
    <div className="relative w-full h-screen overflow-hidden">
      <CustomBackground />

      {/* Account Menu */}
      <div className="absolute top-6 right-6 z-10">
        <div className="relative group">
          <UserCircle size={48} className="text-gray-800 hover:text-gray-600 cursor-pointer" />
          <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300">
            <Link to="/login" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Login</Link>
            <Link to="/register" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Register</Link>
          </div>
        </div>
      </div>

      <div className="absolute inset-0 flex flex-col items-center justify-center p-4">
        <h1 className="text-4xl font-bold mb-8 text-center text-gray-800">
          {isLogin ? 'Login' : 'Register'}
        </h1>
        <form onSubmit={handleSubmit} className="w-full max-w-md px-4">
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email"
            className="w-full px-6 py-4 text-xl border-2 border-gray-400 rounded-lg bg-white bg-opacity-70 focus:outline-none focus:ring-2 focus:ring-gray-500 mb-6 text-gray-800 placeholder-gray-400"
          />
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password"
            className="w-full px-6 py-4 text-xl border-2 border-gray-400 rounded-lg bg-white bg-opacity-70 focus:outline-none focus:ring-2 focus:ring-gray-500 mb-6 text-gray-800 placeholder-gray-400"
          />
          <button
            type="submit"
            className="w-full bg-gray-800 text-white py-4 px-6 text-xl rounded-lg border border-transparent hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 transition duration-300"
          >
            {isLogin ? 'Login' : 'Register'}
          </button>
        </form>
        <p className="mt-4 text-gray-600">
          {isLogin ? "Don't have an account? " : "Already have an account? "}
          <Link to={isLogin ? "/register" : "/login"} className="text-gray-800 hover:underline">
            {isLogin ? "Register" : "Login"}
          </Link>
        </p>
      </div>
    </div>
  );
};

export default LoginPage;
import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { UserCircle, Loader } from 'lucide-react';
import { authService } from '../services/AuthService';
import CustomBackground from "../Componants/CustomBackground.jsx";

const LoginPage = ({ isLogin }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      if (isLogin) {
        // Perform login and handle successful response
        const loginResponse = await authService.login(email, password);
        if (loginResponse.jwt_token) {
          // Redirect to account page upon successful login
          navigate('/account');
        }
      } else {
        // Perform registration
        const registerResponse = await authService.register(email, password);
        if (registerResponse.message) {
          // After successful registration, you might want to navigate to the login page or auto-login
          navigate('/login');
        }
      }
    } catch (err) {
      console.error('Auth error:', err);
      setError(err.response?.data?.message || err.message || 'An error occurred. Please try again.');
    } finally {
      setIsLoading(false);
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
              <Link to="/login" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Login</Link>
              <Link to="/register" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Register</Link>
            </div>
          </div>
        </div>

        <div className="absolute inset-0 flex flex-col items-center justify-center p-4">
          <h1 className="text-4xl font-bold mb-8 text-center text-gray-800">
            {isLogin ? 'Login' : 'Register'}
          </h1>
          {error && <p className="text-red-500 mb-4">{error}</p>}
          <form onSubmit={handleSubmit} className="w-full max-w-md px-4">
            <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Email"
                className="w-full px-6 py-4 text-xl border-2 border-gray-400 rounded-lg bg-white bg-opacity-70 focus:outline-none focus:ring-2 focus:ring-gray-500 mb-6 text-gray-800 placeholder-gray-400"
                disabled={isLoading}
            />
            <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Password"
                className="w-full px-6 py-4 text-xl border-2 border-gray-400 rounded-lg bg-white bg-opacity-70 focus:outline-none focus:ring-2 focus:ring-gray-500 mb-6 text-gray-800 placeholder-gray-400"
                disabled={isLoading}
            />
            <button
                type="submit"
                className="w-full bg-gray-800 text-white py-4 px-6 text-xl rounded-lg border border-transparent hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 transition duration-300 flex items-center justify-center"
                disabled={isLoading}
            >
              {isLoading ? (
                  <>
                    <Loader className="animate-spin mr-2" size={24} />
                    {isLogin ? 'Logging in...' : 'Registering...'}
                  </>
              ) : (
                  isLogin ? 'Login' : 'Register'
              )}
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

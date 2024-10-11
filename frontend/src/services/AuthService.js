// services/AuthService.js
import axios from 'axios';
import { cookieService } from './CookieService';

const API_URL = '/auth'; // This matches the nginx configuration
const USER_COOKIE_NAME = 'user_data';

// Create an axios instance with default config
const axiosInstance = axios.create({
  baseURL: 'http://localhost', // or your actual domain
});

export const authService = {
  register: async (email, password) => {
    try {
      const response = await axiosInstance.post(`${API_URL}/register`, { email, password });
      if (response.data.token) {
        cookieService.setCookie(USER_COOKIE_NAME, response.data);
      }
      return response.data;
    } catch (error) {
      console.error('Registration error:', error.response ? error.response.data : error.message);
      throw error.response ? error.response.data : error.message;
    }
  },

  login: async (email, password) => {
    try {
      const response = await axiosInstance.post(`${API_URL}/login`, { email, password });
      if (response.data.token) {
        cookieService.setCookie(USER_COOKIE_NAME, response.data);
      }
      return response.data;
    } catch (error) {
      console.error('Login error:', error.response ? error.response.data : error.message);
      throw error.response ? error.response.data : error.message;
    }
  },

  logout: async () => {
    try {
      const userData = cookieService.getCookie(USER_COOKIE_NAME);
      if (userData && userData.token) {
        await axiosInstance.post(`${API_URL}/logout`, {}, {
          headers: { Authorization: `Bearer ${userData.token}` }
        });
      }
      cookieService.deleteCookie(USER_COOKIE_NAME);
    } catch (error) {
      console.error('Logout error', error);
    }
  },

  getCurrentUser: () => {
    return cookieService.getCookie(USER_COOKIE_NAME);
  },

  isAuthenticated: () => {
    const userData = cookieService.getCookie(USER_COOKIE_NAME);
    return !!(userData && userData.token);
  }
};
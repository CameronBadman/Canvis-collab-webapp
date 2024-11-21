import axios from 'axios';
import { cookieService } from './CookieService';


const backendURL = "http://localhost:8080/"


const axiosInstance = axios.create({
  baseURL: backendURL + "account",
  timeout: 5000,
});

export const authService = {
  // Register user with email and password
  register: async (email, password) => {
    try {
      console.log(`Attempting to register user: ${email}`);
      const response = await axiosInstance.post(`/register`, { email, password });
      console.log('Registration response:', response.data);

      // Since register doesn't return any auth data, just handle the response
      return response.data;
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  },

  // Log in user with email and password
  login: async (email, password) => {
    try {
      console.log(`Attempting to log in user: ${email}`);
      const response = await axiosInstance.post(`/login`, { email, password });
      console.log('Login response:', response.data);

      // Store JWT token in cookies
      if (response.data.jwt_token) {
        cookieService.setCookie('user_data', response.data);
      }

      return response.data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  // Confirm user registration (could be email verification)
  confirm: async (confirmation_code) => {
    try {

      const response = await axiosInstance.post(`/confirm`, {confirmation_code });
      console.log('Confirmation response:', response.data);

      return response.data;
    } catch (error) {
      console.error('Confirmation error:', error);
      throw error;
    }
  },

  // Log out the user
  logout: async () => {
    try {
      const userData = cookieService.getCookie('user_data');
      if (userData && userData.jwt_token) {
        await axiosInstance.post(`/logout`, {}, {
          headers: { Authorization: `Bearer ${userData.jwt_token}` }
        });
      }
      cookieService.deleteCookie('user_data');
    } catch (error) {
      console.error('Logout error:', error);
    }
  },

  // Get current user from cookie
  getCurrentUser: () => {
    return cookieService.getCookie('user_data');
  },

  // Check if user is authenticated by checking if the token exists in cookies
  isAuthenticated: () => {
    const userData = cookieService.getCookie('user_data');
    return !!(userData && userData.jwt_token);
  }
};

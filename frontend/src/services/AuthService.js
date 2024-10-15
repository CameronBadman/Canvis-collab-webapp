import axios from 'axios';
import { cookieService } from './CookieService';

const AUTH_API_URL = '/api/auth';
const CANVAS_API_URL = '/api/canvas';
const USER_COOKIE_NAME = 'user_data';

const axiosInstance = axios.create({
  baseURL: 'http://localhost:8000', // Point to the load balancer
  timeout: 5000,
});

export const authService = {
  register: async (email, password) => {
    try {
      console.log(`Attempting to register user: ${email}`);
      const response = await axiosInstance.post(`${AUTH_API_URL}/register`, { email, password });
      console.log('Registration response:', response.data);
      if (response.data.token) {
        cookieService.setCookie(USER_COOKIE_NAME, response.data);
        
        try {
          await axiosInstance.post(`${CANVAS_API_URL}/user`, 
            { firebase_uid: response.data.uid },
            { headers: { 
                Authorization: `Bearer ${response.data.token}`,
                'X-Firebase-UID': response.data.uid
              }
            }
          );
          console.log('User created in Canvas API');
        } catch (error) {
          console.error('Registration error:', error);
          throw error;
        }
      }
      return response.data;
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  },

  login: async (email, password) => {
    try {
      console.log(`Attempting to log in user: ${email}`);
      const response = await axiosInstance.post(`${AUTH_API_URL}/login`, { email, password });
      console.log('Login response:', response.data);
      if (response.data.token) {
        cookieService.setCookie(USER_COOKIE_NAME, response.data);
      }
      return response.data;
    } catch (error) {
      console.error('Login error:', error);
      if (error.response) {
        console.error('Error response:', error.response.data);
        console.error('Error status:', error.response.status);
        console.error('Error headers:', error.response.headers);
      } else if (error.request) {
        console.error('No response received:', error.request);
      } else {
        console.error('Error setting up request:', error.message);
      }
      throw error;
    }
  },

  logout: async () => {
    try {
      const userData = cookieService.getCookie(USER_COOKIE_NAME);
      if (userData && userData.token) {
        await axiosInstance.post(`${AUTH_API_URL}/logout`, {}, {
          headers: { Authorization: `Bearer ${userData.token}` }
        });
      }
      cookieService.deleteCookie(USER_COOKIE_NAME);
    } catch (error) {
      console.error('Login error:', error);
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
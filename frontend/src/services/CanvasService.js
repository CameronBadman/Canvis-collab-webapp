import axios from 'axios';
import { cookieService } from './CookieService';

const API_URL = '/api/canvas';
const USER_COOKIE_NAME = 'user_data';

const axiosInstance = axios.create({
  baseURL: 'http://localhost:80', // Point to the load balancer
  timeout: 5000,
});

const getAuthHeader = () => {
    const userData = cookieService.getCookie(USER_COOKIE_NAME);
    return userData && userData.token
      ? {
          Authorization: `Bearer ${userData.token}`,
          'X-Firebase-UID': userData.uid
        }
      : {};
  };

export const canvasService = {
  createCanvas: async (name) => {
    try {
      const response = await axiosInstance.post(`${API_URL}/canvas`, { name }, { headers: getAuthHeader() });
      return response.data;
    } catch (error) {
      console.error('Create canvas error:', error);
      throw error;
    }
  },

  getCanvas: async (id) => {
    try {
      const response = await axiosInstance.get(`${API_URL}/canvas/${id}`, { headers: getAuthHeader() });
      return response.data;
    } catch (error) {
      console.error('Get canvas error:', error);
      throw error;
    }
  },

  updateCanvas: async (id, name) => {
    try {
      const response = await axiosInstance.put(`${API_URL}/canvas/${id}`, { name }, { headers: getAuthHeader() });
      return response.data;
    } catch (error) {
      console.error('Update canvas error:', error);
      throw error;
    }
  },

  deleteCanvas: async (id) => {
    try {
      await axiosInstance.delete(`${API_URL}/canvas/${id}`, { headers: getAuthHeader() });
    } catch (error) {
      console.error('Delete canvas error:', error);
      throw error;
    }
  },

  getUserCanvases: async () => {
    try {
      const response = await axiosInstance.get(`${API_URL}/user/canvases`, { headers: getAuthHeader() });
      return response.data;
    } catch (error) {
      console.error('Get user canvases error:', error);
      throw error;
    }
  }
};
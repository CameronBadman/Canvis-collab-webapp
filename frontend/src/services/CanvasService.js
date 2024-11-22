import axios from 'axios';
import { cookieService } from './CookieService';

const backendURL = "http://localhost:8080/";

const axiosInstance = axios.create({
  baseURL: backendURL + "canvas",
  timeout: 5000,
});

// Utility function to retrieve the authentication header
const getAuthHeader = () => {
  const userData = cookieService.getCookie('user_data');
  return userData && userData.jwt_token
      ? { Authorization: `Bearer ${userData.jwt_token}` }
      : {};
};

export const canvasService = {
  // Create a new canvas
  createCanvas: async (canvasName) => {
    try {
      console.log(`Attempting to create canvas: ${canvasName}`);
      const response = await axiosInstance.post(
          `/create`,
          { canvas_name: canvasName },
          { headers: getAuthHeader() }
      );
      console.log('Create canvas response:', response.data);
      return response.data;
    } catch (error) {
      console.error('Create canvas error:', error.response?.data || error.message);
      throw error;
    }
  },

  // Get all canvases for the current user
  getUserCanvases: async () => {
    try {
      console.log(`Fetching all canvases for the user.`);
      const response = await axiosInstance.get(`/canvases`, {
        headers: getAuthHeader(),
      });
      console.log('Get user canvases response:', response.data);
      return response.data;
    } catch (error) {
      console.error('Get user canvases error:', error.response?.data || error.message);
      throw error;
    }
  },

  // Get details of a single canvas by ID
  getCanvas: async (canvasId) => {
    try {
      console.log(`Fetching canvas with ID: ${canvasId}`);
      const response = await axiosInstance.get(`/canvases/${canvasId}`, {
        headers: getAuthHeader(),
      });
      console.log('Get canvas response:', response.data);
      return response.data;
    } catch (error) {
      console.error('Get canvas error:', error.response?.data || error.message);
      throw error;
    }
  },

  // Update a canvas by ID
  updateCanvas: async (canvasId, canvasName) => {
    try {
      console.log(`Updating canvas with ID: ${canvasId}`);
      const response = await axiosInstance.put(
          `/canvases/${canvasId}`,
          { canvas_name: canvasName },
          { headers: getAuthHeader() }
      );
      console.log('Update canvas response:', response.data);
      return response.data;
    } catch (error) {
      console.error('Update canvas error:', error.response?.data || error.message);
      throw error;
    }
  },

  // Delete a canvas by ID
  deleteCanvas: async (canvasId) => {
    try {
      console.log(`Deleting canvas with ID: ${canvasId}`);
      await axiosInstance.delete(`/canvases/${canvasId}`, {
        headers: getAuthHeader(),
      });
      console.log('Delete canvas successful.');
    } catch (error) {
      console.error('Delete canvas error:', error.response?.data || error.message);
      throw error;
    }
  },
};

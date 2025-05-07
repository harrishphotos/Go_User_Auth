import * as axiosTypes from "axios"; // Import all types
import axios from "axios"; // Import axios instance
import {
  getAccessToken,
  setAccessToken,
  removeAccessToken,
} from "../utils/tokenUtil"; // Import your utility functions

export const BASE_URL = "http://localhost:3000/api";

const axiosInstance = axios.create({
  baseURL: BASE_URL,
  withCredentials: true, // Ensures credentials (cookies) are sent with requests
});

axiosInstance.interceptors.request.use(
  (config: axiosTypes.AxiosRequestConfig) => {
    // Use the correct type here
    const token = getAccessToken();
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

axiosInstance.interceptors.response.use(
  (response) => response, // Return the response if everything is fine
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshResponse = await axios.post(
          `${BASE_URL}/auth/refresh`,
          {},
          { withCredentials: true }
        );
        const newAccessToken = refreshResponse.data.accessToken;

        // Save the new access token with the correct key
        setAccessToken(newAccessToken);
        originalRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
        return axiosInstance(originalRequest); // Retry original request with new token
      } catch (refreshError) {
        removeAccessToken(); // Handle refresh token failure (log out user)
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error); // Reject if the error is not 401
  }
);

export default axiosInstance; // Default export

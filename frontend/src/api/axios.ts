<<<<<<< Updated upstream
import * as axiosTypes from "axios"; // Import all types
import axios from "axios"; // Import axios instance
=======
import axios from "axios";
import type { AxiosInstance } from "axios";
import type { AxiosRequestConfig } from "axios";
import type { AxiosResponse } from "axios";
import type { AxiosError } from "axios";
>>>>>>> Stashed changes
import {
  getAccessToken,
  setAccessToken,
  removeAccessToken,
<<<<<<< Updated upstream
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
=======
} from "../utils/tokenUtil";

export const BASE_URL = "http://localhost:3000/api";

const axiosInstance: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

axiosInstance.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const token = getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest?._retry) {
      originalRequest._retry = true;

      try {
        const refreshResponse = await axios.post<{
          accessToken: string;
        }>(
>>>>>>> Stashed changes
          `${BASE_URL}/auth/refresh`,
          {},
          { withCredentials: true }
        );
<<<<<<< Updated upstream
        const newAccessToken = refreshResponse.data.accessToken;

        // Save the new access token with the correct key
        setAccessToken(newAccessToken);
        originalRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
        return axiosInstance(originalRequest); // Retry original request with new token
      } catch (refreshError) {
        removeAccessToken(); // Handle refresh token failure (log out user)
=======

        const { accessToken: newAccessToken } = refreshResponse.data;
        setAccessToken(newAccessToken);
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        removeAccessToken();
>>>>>>> Stashed changes
        return Promise.reject(refreshError);
      }
    }

<<<<<<< Updated upstream
    return Promise.reject(error); // Reject if the error is not 401
  }
);

export default axiosInstance; // Default export
=======
    return Promise.reject(error);
  }
);

export default axiosInstance;
>>>>>>> Stashed changes

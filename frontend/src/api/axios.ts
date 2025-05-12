import axios from "axios";

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig,
} from "axios";

import {
  getAccessToken,
  setAccessToken,
  removeAccessToken,
} from "../utils/tokenUtil";

export const BASE_URL = "http://localhost:3000";

const axiosInstance: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

axiosInstance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig & {
      _retry?: boolean;
      headers?: any;
    };

    if (error.response?.status === 401 && !originalRequest?._retry) {
      const token = getAccessToken();
      if (!token) {
        removeAccessToken();
        return Promise.reject(error); // ⛔ No retry if no token
      }

      originalRequest._retry = true;
      try {
        const refreshResponse = await axiosInstance.post<{
          accessToken: string;
        }>("/api/auth/refresh", {});
        const { accessToken: newAccessToken } = refreshResponse.data;
        setAccessToken(newAccessToken);
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        removeAccessToken();
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default axiosInstance;

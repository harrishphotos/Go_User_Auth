import axiosPkg from "axios";
const axios = axiosPkg.default;

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
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
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
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

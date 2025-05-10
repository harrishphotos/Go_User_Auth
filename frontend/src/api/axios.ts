import axios from "axios";
import type { AxiosInstance } from "axios";
import type { AxiosRequestConfig } from "axios";
import type { AxiosResponse } from "axios";
import type { AxiosError } from "axios";
import {
  getAccessToken,
  setAccessToken,
  removeAccessToken,
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
        // const refreshResponse = await axios.post<{
        //   accessToken: string;
        // }>(`${BASE_URL}/auth/refresh`, {}, { withCredentials: true });
        const refreshResponse = await axiosInstance.post<{
          accessToken: string;
        }>("/auth/refresh", {});

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

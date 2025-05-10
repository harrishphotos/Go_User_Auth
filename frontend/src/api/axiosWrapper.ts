import { useEffect } from "react";
import type {
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from "axios";
import { useAuth } from "../context/AuthContext";
import axiosInstance from "./axios";

const useAxiosWrapper = () => {
  const { accessToken, setAuth, logout } = useAuth();

  useEffect(() => {
    const publicPaths: string[] = [
      "/auth/login",
      "/auth/register",
      "/auth/forgot-password",
      "/auth/reset-password",
      "/auth/verify-email",
    ];

    // Request interceptor
    const requestInterceptor = axiosInstance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const isPublic = publicPaths.some((path) => config.url?.includes(path));

        if (!isPublic && accessToken && config.headers) {
          config.headers["Authorization"] = `Bearer ${accessToken}`;
        }

        return config;
      },
      (error: AxiosError) => Promise.reject(error)
    );

    // Response interceptor
    const responseInterceptor = axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as AxiosRequestConfig & {
          _retry?: boolean;
        };

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const refreshResponse = await axiosInstance.post<{
              accessToken: string;
            }>("/auth/refresh", {});

            const { accessToken: newAccessToken } = refreshResponse.data;

            localStorage.setItem("access_token", newAccessToken);
            setAuth({ id: "userId", email: "userEmail" }, newAccessToken);

            if (originalRequest.headers) {
              originalRequest.headers[
                "Authorization"
              ] = `Bearer ${newAccessToken}`;
            }

            return axiosInstance(originalRequest);
          } catch (refreshError) {
            logout();
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );

    return () => {
      axiosInstance.interceptors.request.eject(requestInterceptor);
      axiosInstance.interceptors.response.eject(responseInterceptor);
    };
  }, [accessToken, setAuth, logout]);

  return axiosInstance;
};

export default useAxiosWrapper;

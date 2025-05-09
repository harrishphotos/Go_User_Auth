// axiosWrapper.ts
import { useAuth } from "../context/AuthContext";
import axiosInstance, { BASE_URL } from "./axios";

const useAxiosWrapper = () => {
  const { accessToken, setAuth, logout } = useAuth();

  // Set up interceptors (only once per component mount)
  axiosInstance.interceptors.request.use(
<<<<<<< Updated upstream
    (config) => {
=======
    (config: { headers: { [x: string]: string } }) => {
>>>>>>> Stashed changes
      if (accessToken) {
        config.headers["Authorization"] = `Bearer ${accessToken}`;
      }
      return config;
    },
<<<<<<< Updated upstream
    (error) => Promise.reject(error)
  );

  axiosInstance.interceptors.response.use(
    (response) => response,
    async (error) => {
=======
    (error: any) => Promise.reject(error)
  );

  axiosInstance.interceptors.response.use(
    (response: any) => response,
    async (error: { config: any; response: { status: number } }) => {
>>>>>>> Stashed changes
      const originalRequest = error.config;
      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;

        try {
          const refreshResponse = await axiosInstance.post("/auth/refresh", {});
          const { accessToken: newAccessToken, refreshToken } =
            refreshResponse.data;

          localStorage.setItem("access_token", newAccessToken);
          localStorage.setItem("refresh_token", refreshToken);

          // Replace with real user details
          setAuth({ id: "userId", email: "userEmail" }, newAccessToken);

          originalRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
          return axiosInstance(originalRequest);
        } catch (refreshError) {
          logout();
          return Promise.reject(refreshError);
        }
      }
      return Promise.reject(error);
    }
  );

  return axiosInstance;
};

export default useAxiosWrapper;

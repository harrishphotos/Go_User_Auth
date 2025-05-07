import axiosInstance from "../api/axios"; // Import the axios instance

export const refreshAuthToken = async () => {
  try {
    const response = await axiosInstance.post("/auth/refresh");
    return response.data.accessToken; // Return the new access token
  } catch (error) {
    console.error("Error refreshing token:", error);
    throw error; // Propagate the error
  }
};

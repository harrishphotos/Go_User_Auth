// src/utils/tokenUtil.ts

export const getAccessToken = (): string | null => {
  return localStorage.getItem("access_token"); // Correct the key here
};

export const setAccessToken = (token: string) => {
  localStorage.setItem("access_token", token); // Correct the key here
};

export const removeAccessToken = () => {
  localStorage.removeItem("access_token"); // Correct the key here
};

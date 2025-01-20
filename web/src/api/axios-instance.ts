import axios, { AxiosInstance } from "axios";

const API_URL: string = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export const axiosInstance: AxiosInstance = axios.create({
  baseURL: API_URL,
});
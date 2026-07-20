import axios from "axios";
import { BASE_URL } from "../constants/api";
import { Task } from "../types/task";

const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

export const getTasks = (params?: any) =>
  api.get("/tasks", { params });

export const createTask = (task: Partial<Task>) =>
  api.post("/tasks", task);

export const updateTask = (id: number, task: Partial<Task>) =>
  api.put(`/tasks/${id}`, task);

export const deleteTask = (id: number) =>
  api.delete(`/tasks/${id}`);

export default api; 
// src/services/api.ts
import axios from 'axios';

const API_URL = 'http://localhost:4000';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Todo API functions
export const getTodos = async () => {
  const response = await api.get('/api/todos');
  return response.data;
};

export const createTodo = async (todo: { title: string; body: string }) => {
  const response = await api.post('/api/todos', todo);
  return response.data;
};

export const markTodoDone = async (id: number) => {
  const response = await api.patch(`/api/todos/${id}/done`);
  return response.data;
};

export const markTodoUndone = async (id: number) => {
  const response = await api.patch(`/api/todos/${id}/undone`);
  return response.data;
};

export const deleteTodo = async (id: number) => {
  const response = await api.delete(`/api/todos/${id}/delete`);
  return response.data;
};

export const getUsers = async () => {
  const response = await api.get('/users');
  return response.data;
};

export const createUser = async (user: any) => {
  const response = await api.post('/users', user);
  return response.data;
};

export default api;
import { api } from 'boot/axios';

export async function login(email: string, password: string) {
  const { data } = await api.post('/auth/login', { email, password });
  return data;
}

export async function register(email: string, username: string, password: string) {
  const { data } = await api.post('/auth/register', { email, username, password });
  return data;
}

export async function getMe() {
  const { data } = await api.get('/auth/me');
  return data;
}

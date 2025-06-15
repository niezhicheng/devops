import axios from 'axios';
import { Message } from '@arco-design/web-vue';

const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
});

// 获取仓库列表
export async function getRepositories() {
  const response = await request.get('/api/repositories');
  return response.data;
}

// 创建仓库
export async function createRepository(data) {
  const response = await request.post('/api/repositories', data);
  return response.data;
}

// 更新仓库
export async function updateRepository(id, data) {
  const response = await request.put(`/api/repositories/${id}`, data);
  return response.data;
}

// 删除仓库
export async function deleteRepository(id) {
  await request.delete(`/api/repositories/${id}`);
}

// 测试仓库连接
export async function testRepository(params) {
  const response = await request.post('/api/repositories/test', params);
  return response.data;
} 
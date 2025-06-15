import axios from 'axios';
import { Message } from '@arco-design/web-vue';

export interface Repository {
  id: number;
  name: string;
  platform: 'github' | 'gitlab';
  url: string;
  token: string;
  defaultBranch?: string;
  lastSyncAt?: string;
  status?: string;
}

const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
});
//
// // 请求拦截器
// request.interceptors.request.use(
//   (config) => {
//     const token = localStorage.getItem('token');
//     if (token) {
//       config.headers.Authorization = `Bearer ${token}`;
//     }
//     return config;
//   },
//   (error) => {
//     return Promise.reject(error);
//   }
// );
//
// // 响应拦截器
// request.interceptors.response.use(
//   (response) => {
//     return response.data;
//   },
//   (error) => {
//     if (error.response) {
//       switch (error.response.status) {
//         case 401:
//           // 未授权，跳转到登录页
//           window.location.href = '/login';
//           break;
//         case 403:
//           Message.error('没有权限访问');
//           break;
//         case 404:
//           Message.error('请求的资源不存在');
//           break;
//         case 500:
//           Message.error('服务器错误');
//           break;
//         default:
//           Message.error(error.response.data.message || '请求失败');
//       }
//     } else {
//       Message.error('网络错误，请检查网络连接');
//     }
//     return Promise.reject(error);
//   }
// );

// 获取仓库列表
export async function getRepositories(): Promise<Repository[]> {
  const response = await request.get('/api/repositories');
  return response.data;
}

// 创建仓库
export async function createRepository(data: Omit<Repository, 'id'>): Promise<Repository> {
  const response = await request.post('/api/repositories', data);
  return response.data;
}

// 更新仓库
export async function updateRepository(id: number, data: Partial<Repository>): Promise<Repository> {
  const response = await request.put(`/api/repositories/${id}`, data);
  return response.data;
}

// 删除仓库
export async function deleteRepository(id: number): Promise<void> {
  await request.delete(`/api/repositories/${id}`);
}

// 测试仓库连接
export interface TestRepositoryParams {
  url: string;
  token: string;
  platform: 'github' | 'gitlab';
}

export interface TestRepositoryResult {
  success: boolean;
  message?: string;
  data?: {
    username: string;
    name: string;
  };
}

export async function testRepository(params: TestRepositoryParams): Promise<TestRepositoryResult> {
  const response = await request.post('/api/repositories/test', params);
  return response.data;
}

import axios from 'axios';
import type { AxiosResponse } from 'axios';

// 配置axios默认值
axios.defaults.baseURL = 'http://localhost:8080/api';  // 设置后端API基础URL
axios.defaults.timeout = 5000;  // 设置超时时间

// 添加请求拦截器
axios.interceptors.request.use(
  (config) => {
    console.log('发送请求:', config.url, config.params || config.data);
    return config;
  },
  (error) => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 添加响应拦截器
axios.interceptors.response.use(
  (response) => {
    console.log('收到响应:', response.data);
    return response;
  },
  (error) => {
    console.log(error)
    // console.error('响应错误:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);
export interface Repository {
  id: number;
  name: string;
  platform: 'github' | 'gitlab';
  url: string;
  token: string;
  defaultBranch?: string;
  status?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Branch {
  name: string;
  isHead: boolean;
}

export interface Commit {
  hash: string;
  author: string;
  message: string;
  date: string;
}

export interface PaginationParams {
  current: number;
  pageSize: number;
}

export interface PaginatedResponse<T> {
  list: T[];
  total: number;
}

// 获取仓库列表
export function getRepositories(params: PaginationParams) {
  return axios.get<PaginatedResponse<Repository>>('/repositories', { params });
}

// 创建仓库
export function createRepository(data: Partial<Repository>) {
  return axios.post<Repository>('/repositories', data);
}

// 更新仓库
export function updateRepository(id: number, data: Partial<Repository>) {
  return axios.put<Repository>(`/repositories/${id}`, data);
}

// 删除仓库
export function deleteRepository(id: number) {
  return axios.delete(`/repositories/${id}`);
}

// 测试仓库连接
export function testRepository(data: { url: string; token: string; platform: 'github' | 'gitlab' }) {
  return axios.post<{ success: boolean; message?: string }>('/repositories/test', data);
}

// 获取仓库分支列表
export function getBranches(repoId: number) {
  return axios.get<Branch[]>(`/repositories/${repoId}/branches`);
}

// 获取提交历史
export function getCommits(repoId: number, branch?: string) {
  return axios.get<Commit[]>(`/repositories/${repoId}/commits`, {
    params: { branch },
  });
}

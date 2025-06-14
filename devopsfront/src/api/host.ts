import axios from 'axios';
import type { AxiosResponse } from 'axios';

// 配置axios默认值
axios.defaults.baseURL = 'http://localhost:8080';  // 设置后端API基础URL
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

export interface HostRecord {
  id: number;
  name: string;
  ip: string;
  port: number;
  username: string;
  password: string;
  description?: string;
  createdTime: string;
  updatedTime: string;
}

export interface HostParams {
  current: number;
  pageSize: number;
  name?: string;
  ip?: string;
}

export interface AddHostParams {
  name: string;
  ip: string;
  port: number;
  username: string;
  password: string;
  description?: string;
}

export function queryHostList(params: HostParams) {
  return axios.get<{
    list: HostRecord[];
    // total: number;
  }>('/api/host/list', { params });
}

export function addHost(data: AddHostParams) {
  return axios.post<HostRecord>('/api/host/add', data);
}

export function updateHost(id: number, data: Partial<AddHostParams>) {
  return axios.put<HostRecord>(`/api/host/${id}`, data);
}

export function deleteHost(id: number) {
  return axios.delete(`/api/host/${id}`);
}

// SFTP文件信息接口
export interface SftpFileInfo {
  name: string;
  path: string;
  size: number;
  type: 'file' | 'directory';
  modifyTime: string;
  permissions: string;
}

// 获取SFTP文件列表
export function fetchSftpFiles(hostId: number, path: string) {
  return axios.get<{ list: SftpFileInfo[] }>(`/api/host/${hostId}/sftp`, {
    params: { path },
  });
}

// 上传SFTP文件
export function uploadSftpFile(hostId: string, formData: FormData) {
  return axios.post(`/api/host/${hostId}/sftp/upload`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}

// 从SFTP下载文件
export function downloadSftpFile(hostId: number, path: string) {
  return axios.get(`/api/host/${hostId}/sftp/download`, {
    params: { path },
    responseType: 'blob',
  });
}

// 删除SFTP文件
export function deleteSftpFile(hostId: number, path: string) {
  return axios.delete(`/api/host/${hostId}/sftp`, {
    params: { path },
  });
}

// 重命名SFTP文件
export function renameSftpFile(hostId: number, oldPath: string, newPath: string) {
  return axios.put(`/api/host/${hostId}/sftp/rename`, null, {
    params: { oldPath, newPath },
  });
}

// 下载SFTP目录（压缩）
export function downloadSftpDir(hostId: number, path: string) {
  return axios.get(`/api/host/${hostId}/sftp/download-dir`, {
    params: { path },
    responseType: 'blob',
  });
}

// 压缩SFTP目录
export function compressSftpDir(hostId: string, path: string) {
  return axios.post(`/api/host/${hostId}/sftp/compress`, null, {
    params: { path }
  });
}

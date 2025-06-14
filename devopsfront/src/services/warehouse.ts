import { request } from '@/utils/request';

export interface Repository {
  id: number;
  name: string;
  url: string;
  token: string;
  defaultBranch?: string;
  lastSyncAt?: string;
  status?: string;
}

// 获取仓库列表
export async function getRepositories(): Promise<Repository[]> {
  return request('/api/repositories');
}

// 创建仓库
export async function createRepository(data: Omit<Repository, 'id'>): Promise<Repository> {
  return request('/api/repositories', {
    method: 'POST',
    data,
  });
}

// 更新仓库
export async function updateRepository(id: number, data: Partial<Repository>): Promise<Repository> {
  return request(`/api/repositories/${id}`, {
    method: 'PUT',
    data,
  });
}

// 删除仓库
export async function deleteRepository(id: number): Promise<void> {
  return request(`/api/repositories/${id}`, {
    method: 'DELETE',
  });
} 
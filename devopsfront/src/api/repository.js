import axios from 'axios';

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
    return Promise.reject(error);
  }
);

// 获取仓库列表
export function getRepositories(params) {
  return axios.get('/repositories', { params });
}

// 创建仓库
export function createRepository(data) {
  return axios.post('/repositories', data);
}

// 更新仓库
export function updateRepository(id, data) {
  return axios.put(`/repositories/${id}`, data);
}

// 删除仓库
export function deleteRepository(id) {
  return axios.delete(`/repositories/${id}`);
}

// 测试仓库连接
export function testRepository(data) {
  return axios.post('/repositories/test', data);
}

// 获取仓库分支列表
export function getBranches(repoId) {
  return axios.get(`/repositories/${repoId}/branches`);
}

// 获取提交历史
export function getCommits(repoId, branch) {
  return axios.get(`/repositories/${repoId}/commits`, {
    params: { branch },
  });
} 
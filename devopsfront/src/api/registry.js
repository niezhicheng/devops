import request from '@/utils/request';

// 获取镜像仓库列表
export function getDockerRegistries(params) {
  return request({
    url: '/api/docker-registries',
    method: 'get',
    params,
  });
}

// 创建镜像仓库
export function createDockerRegistry(data) {
  return request({
    url: '/api/docker-registries',
    method: 'post',
    data,
  });
}

// 更新镜像仓库
export function updateDockerRegistry(id, data) {
  return request({
    url: `/api/docker-registries/${id}`,
    method: 'put',
    data,
  });
}

// 删除镜像仓库
export function deleteDockerRegistry(id) {
  return request({
    url: `/api/docker-registries/${id}`,
    method: 'delete',
  });
}

// 测试镜像仓库连接
export function testDockerRegistryConnection(data) {
  return request({
    url: '/api/docker-registries/test-connection',
    method: 'post',
    data,
  });
} 
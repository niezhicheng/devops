import request from '@/utils/request'

export function getProjects(params) {
  return request({
    url: '/api/projects',
    method: 'get',
    params
  })
}

export function getProject(id) {
  return request({
    url: `/api/projects/${id}`,
    method: 'get'
  })
}

export function createProject(data) {
  return request({
    url: '/api/projects',
    method: 'post',
    data
  })
}

export function updateProject(id, data) {
  return request({
    url: `/api/projects/${id}`,
    method: 'put',
    data
  })
}

export function deleteProject(id) {
  return request({
    url: `/api/projects/${id}`,
    method: 'delete'
  })
} 
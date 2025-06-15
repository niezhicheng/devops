import Layout from '@/layout/content'
import { createRouter, createWebHashHistory } from "vue-router"

export const routes = [
  {
    path: '/404',
    component: () => import('@/views/404.vue'),
    hidden: true
  },

  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'icon-dashboard' }
      }
    ]
  },

  {
    path: '/devops',
    component: Layout,
    redirect: '/devops/host',
    meta: { title: 'DevOps管理', icon: 'icon-apps' },
    children: [
      {
        path: 'host',
        name: 'HostList',
        component: () => import('@/views/host/index.vue'),
        meta: { title: '主机管理' }
      },
      {
        path: 'project',
        name: 'ProjectList',
        component: () => import('@/views/project/index.vue'),
        meta: { title: '项目管理' }
      },
      {
        path: 'warehouse',
        name: 'WarehouseList',
        component: () => import('@/views/warehouse/index.vue'),
        meta: { title: '仓库管理' }
      },
      {
        path: 'registry',
        name: 'RegistryList',
        component: () => import('@/views/registry/index.vue'),
        meta: { title: '镜像中心' }
      }
    ]
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/404',
    hidden: true
  }
]

const buildRouter = () => createRouter({
  history: createWebHashHistory(),
  routes
})

const router = buildRouter()

export default router

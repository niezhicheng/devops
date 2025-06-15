<template>
  <div class="dashboard-container">
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>
            <div class="card-title">
              <icon-computer />
              <span>主机总数</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.hostCount }}</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>
            <div class="card-title">
              <icon-apps />
              <span>项目总数</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.projectCount }}</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>
            <div class="card-title">
              <icon-code />
              <span>仓库总数</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.warehouseCount }}</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <template #title>
            <div class="card-title">
              <icon-image />
              <span>镜像总数</span>
            </div>
          </template>
          <div class="stat-value">{{ stats.registryCount }}</div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="16">
        <a-card title="最近活动">
          <a-list :data="activities">
            <template #item="{ item }">
              <a-list-item>
                <div class="activity-item">
                  <span class="activity-time">{{ item.time }}</span>
                  <span class="activity-content">{{ item.content }}</span>
                </div>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="系统状态">
          <a-descriptions :data="systemInfo" layout="inline-horizontal" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const stats = reactive({
  hostCount: 0,
  projectCount: 0,
  warehouseCount: 0,
  registryCount: 0
})

const activities = [
  {
    time: '10:00',
    content: '添加了新的主机 192.168.1.100'
  },
  {
    time: '09:30',
    content: '创建了新项目 "Web应用"'
  },
  {
    time: '09:00',
    content: '更新了仓库配置'
  },
  {
    time: '08:30',
    content: '部署了新版本 v1.2.0'
  }
]

const systemInfo = [
  {
    label: '系统版本',
    value: 'v1.0.0'
  },
  {
    label: 'Node版本',
    value: 'v16.14.0'
  },
  {
    label: 'CPU使用率',
    value: '32%'
  },
  {
    label: '内存使用率',
    value: '45%'
  }
]
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.stat-card {
  height: 120px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: var(--color-text-1);
  margin-top: 16px;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.activity-time {
  color: var(--color-text-3);
  min-width: 60px;
}

.activity-content {
  color: var(--color-text-1);
}
</style> 
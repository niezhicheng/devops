<template>
  <div class="project-list">
    <div class="header">
      <h2>项目管理</h2>
      <a-button type="primary" @click="handleCreate">创建项目</a-button>
    </div>

    <a-card class="filter-card">
      <a-form :model="filterForm" layout="inline" class="filter-form">
        <a-form-item field="name" label="项目名称">
          <a-input
            v-model="filterForm.name"
            placeholder="请输入项目名称"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="environment" label="环境">
          <a-select
            v-model="filterForm.environment"
            placeholder="请选择环境"
            allow-clear
          >
            <a-option value="dev">开发环境</a-option>
            <a-option value="test">测试环境</a-option>
            <a-option value="prod">生产环境</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-table
      :loading="loading"
      :data="projects"
      :pagination="{
        total,
        current: page,
        pageSize,
        showTotal: true,
        showJumper: true,
        showPageSize: true,
      }"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
    >
      <template #columns>
        <a-table-column title="项目名称" data-index="name" />
        <a-table-column title="描述" data-index="description" :ellipsis="true" />
        <a-table-column title="环境" data-index="environment">
          <template #cell="{ record }">
            <a-tag :color="getEnvironmentColor(record.environment)">
              {{ getEnvironmentLabel(record.environment) }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="版本" data-index="version" />
        <a-table-column title="分支" data-index="branch" />
        <a-table-column title="构建状态" data-index="last_build_status">
          <template #cell="{ record }">
            <a-tag :color="getBuildStatusColor(record.last_build_status)">
              {{ record.last_build_status || '未构建' }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="最后构建时间" data-index="last_build_time">
          <template #cell="{ record }">
            {{ record.last_build_time ? formatDate(record.last_build_time) : '-' }}
          </template>
        </a-table-column>
        <a-table-column title="操作" fixed="right" :width="200">
          <template #cell="{ record }">
            <a-space>
              <a-button type="text" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-button type="text" @click="handleBuild(record)">
                构建
              </a-button>
              <a-button type="text" status="danger" @click="handleDelete(record)">
                删除
              </a-button>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <a-modal
      v-model:visible="dialogVisible"
      :title="dialogType === 'create' ? '创建项目' : '编辑项目'"
      @ok="handleSubmit"
      @cancel="dialogVisible = false"
      :mask-closable="false"
      :loading="submitting"
    >
      <a-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-align="right"
        :label-col-props="{ span: 6 }"
        :wrapper-col-props="{ span: 18 }"
      >
        <a-form-item field="name" label="项目名称" validate-trigger="blur">
          <a-input v-model="form.name" placeholder="请输入项目名称" />
        </a-form-item>
        <a-form-item field="description" label="项目描述">
          <a-textarea
            v-model="form.description"
            placeholder="请输入项目描述"
            :auto-size="{ minRows: 3, maxRows: 5 }"
          />
        </a-form-item>
        <a-form-item field="repository_id" label="代码仓库" validate-trigger="change">
          <a-select v-model="form.repository_id" placeholder="请选择代码仓库">
            <a-option
              v-for="repo in repositories"
              :key="repo.id"
              :value="repo.id"
              :label="repo.name"
            />
          </a-select>
        </a-form-item>
        <a-form-item field="branch" label="构建分支" validate-trigger="blur">
          <a-input v-model="form.branch" placeholder="请输入构建分支" />
        </a-form-item>
        <a-form-item field="registry_id" label="镜像仓库" validate-trigger="change">
          <a-select v-model="form.registry_id" placeholder="请选择镜像仓库">
            <a-option
              v-for="registry in registries"
              :key="registry.id"
              :value="registry.id"
              :label="registry.name"
            />
          </a-select>
        </a-form-item>
        <a-form-item field="image_name" label="镜像名称" validate-trigger="blur">
          <a-input v-model="form.image_name" placeholder="请输入镜像名称" />
        </a-form-item>
        <a-form-item field="image_tag" label="镜像标签" validate-trigger="blur">
          <a-input v-model="form.image_tag" placeholder="请输入镜像标签" />
        </a-form-item>
        <a-form-item field="environment" label="环境" validate-trigger="change">
          <a-select v-model="form.environment" placeholder="请选择环境">
            <a-option value="dev">开发环境</a-option>
            <a-option value="test">测试环境</a-option>
            <a-option value="prod">生产环境</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="version" label="版本号" validate-trigger="blur">
          <a-input v-model="form.version" placeholder="请输入版本号" />
        </a-form-item>
        <a-form-item field="build_script" label="构建脚本">
          <a-textarea
            v-model="form.build_script"
            placeholder="请输入构建脚本"
            :auto-size="{ minRows: 5, maxRows: 10 }"
          />
        </a-form-item>
        <a-form-item field="build_timeout" label="构建超时">
          <a-input-number
            v-model="form.build_timeout"
            :min="60"
            :max="7200"
            :step="60"
            placeholder="请输入构建超时时间（秒）"
          />
        </a-form-item>
        <a-form-item field="auto_build" label="自动构建">
          <a-switch v-model="form.auto_build" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { getProjects, createProject, updateProject, deleteProject } from '@/api/project'
import { getRepositories } from '@/api/repository'
import { getDockerRegistries } from '@/api/registry'
import { formatDate } from '@/utils/date'

// 数据列表
const projects = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 筛选表单
const filterForm = reactive({
  name: '',
  environment: ''
})

// 弹窗表单
const dialogVisible = ref(false)
const dialogType = ref('create')
const submitting = ref(false)
const formRef = ref(null)
const form = reactive({
  name: '',
  description: '',
  repository_id: '',
  branch: '',
  registry_id: '',
  image_name: '',
  image_tag: '',
  environment: '',
  version: '',
  build_script: '',
  build_timeout: 600,
  auto_build: false
})

// 仓库和镜像仓库列表
const repositories = ref([])
const registries = ref([])

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入项目名称' },
    { maxLength: 50, message: '项目名称不能超过50个字符' }
  ],
  repository_id: [
    { required: true, message: '请选择代码仓库' }
  ],
  branch: [
    { required: true, message: '请输入构建分支' }
  ],
  registry_id: [
    { required: true, message: '请选择镜像仓库' }
  ],
  image_name: [
    { required: true, message: '请输入镜像名称' }
  ],
  image_tag: [
    { required: true, message: '请输入镜像标签' }
  ],
  environment: [
    { required: true, message: '请选择环境' }
  ],
  version: [
    { required: true, message: '请输入版本号' }
  ]
}

// 获取项目列表
const fetchProjects = async () => {
  loading.value = true
  try {
    const res = await getProjects({
      page: page.value,
      page_size: pageSize.value,
      ...filterForm
    })
    projects.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    Message.error('获取项目列表失败')
  } finally {
    loading.value = false
  }
}

// 获取仓库列表
const fetchRepositories = async () => {
  try {
    const res = await getRepositories()
    repositories.value = res.data
  } catch (error) {
    Message.error('获取仓库列表失败')
  }
}

// 获取镜像仓库列表
const fetchRegistries = async () => {
  try {
    const res = await getDockerRegistries()
    registries.value = res.data
  } catch (error) {
    Message.error('获取镜像仓库列表失败')
  }
}

// 环境标签颜色
const getEnvironmentColor = (environment) => {
  const colors = {
    dev: 'blue',
    test: 'orange',
    prod: 'red'
  }
  return colors[environment] || 'gray'
}

// 环境标签文本
const getEnvironmentLabel = (environment) => {
  const labels = {
    dev: '开发环境',
    test: '测试环境',
    prod: '生产环境'
  }
  return labels[environment] || environment
}

// 构建状态颜色
const getBuildStatusColor = (status) => {
  const colors = {
    success: 'green',
    failed: 'red',
    building: 'blue'
  }
  return colors[status] || 'gray'
}

// 搜索
const handleSearch = () => {
  page.value = 1
  fetchProjects()
}

// 重置
const handleReset = () => {
  filterForm.name = ''
  filterForm.environment = ''
  handleSearch()
}

// 分页
const handlePageChange = (current) => {
  page.value = current
  fetchProjects()
}

// 修改每页条数
const handlePageSizeChange = (size) => {
  pageSize.value = size
  page.value = 1
  fetchProjects()
}

// 创建项目
const handleCreate = () => {
  dialogType.value = 'create'
  dialogVisible.value = true
  form.name = ''
  form.description = ''
  form.repository_id = ''
  form.branch = ''
  form.registry_id = ''
  form.image_name = ''
  form.image_tag = ''
  form.environment = ''
  form.version = ''
  form.build_script = ''
  form.build_timeout = 600
  form.auto_build = false
}

// 编辑项目
const handleEdit = (record) => {
  dialogType.value = 'edit'
  dialogVisible.value = true
  Object.assign(form, record)
}

// 删除项目
const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除项目"${record.name}"吗？`,
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteProject(record.id)
        Message.success('删除成功')
        fetchProjects()
      } catch (error) {
        Message.error('删除失败')
      }
    }
  })
}

// 构建项目
const handleBuild = (record) => {
  // TODO: 实现构建功能
  Message.info('构建功能开发中')
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    if (dialogType.value === 'create') {
      await createProject(form)
      Message.success('创建成功')
    } else {
      await updateProject(form.id, form)
      Message.success('更新成功')
    }
    dialogVisible.value = false
    fetchProjects()
  } catch (error) {
    if (error.message) {
      Message.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

// 初始化
onMounted(() => {
  fetchProjects()
  fetchRepositories()
  fetchRegistries()
})
</script>

<style scoped>
.project-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}
</style>

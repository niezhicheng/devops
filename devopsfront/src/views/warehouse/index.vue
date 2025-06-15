<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getRepositories,
  createRepository,
  updateRepository,
  deleteRepository,
  getBranches,
  getCommits,
} from '@/api/repository'

// 状态变量
const loading = ref(false)
const repositories = ref([])
const modalVisible = ref(false)
const editingId = ref(null)
const formRef = ref()
const detailVisible = ref(false)
const currentRepo = ref(null)
const branches = ref([])
const commits = ref([])
const currentBranch = ref('')
const loadingBranches = ref(false)
const loadingCommits = ref(false)

const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
})

const form = ref({
  name: '',
  platform: 'github',
  url: '',
  token: '',
})

const rules = {
  name: [{ required: true, message: '请输入仓库名称' }],
  url: [{ required: true, message: '请输入仓库地址' }],
  token: [{ required: true, message: '请输入访问令牌' }],
}

// 表格列定义
const columns = [
  {
    title: '仓库名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '平台',
    dataIndex: 'platform',
    key: 'platform',
  },
  {
    title: '仓库地址',
    dataIndex: 'url',
    key: 'url',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '操作',
    key: 'action',
    width: 280,
    fixed: 'right',
    slotName: 'action',
  },
]

// 分支表格列定义
const branchColumns = [
  {
    title: '分支名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '是否默认',
    dataIndex: 'isHead',
    key: 'isHead',
  },
]

// 提交记录表格列定义
const commitColumns = [
  {
    title: '提交信息',
    dataIndex: 'message',
    key: 'message',
  },
  {
    title: '提交者',
    dataIndex: 'author',
    key: 'author',
  },
  {
    title: '提交时间',
    dataIndex: 'date',
    key: 'date',
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
    fixed: 'right',
    slotName: 'commitAction',
  },
]

// 获取仓库列表
const fetchRepositories = async () => {
  loading.value = true
  try {
    const res = await getRepositories({
      current: pagination.value.current,
      pageSize: pagination.value.pageSize,
    })
    console.log("收到的数据", res.data.list)
    // 使用解构赋值来确保响应式更新
    repositories.value = [...res.data.list]
    pagination.value = {
      ...pagination.value,
      total: res.data.total
    }
  } catch (error) {
    Message.error('获取仓库列表失败')
  } finally {
    loading.value = false
  }
}

// 显示添加/编辑模态框
const showModal = (record) => {
  if (record) {
    editingId.value = record.id
    form.value = {
      name: record.name,
      platform: record.platform,
      url: record.url,
      token: record.token,
    }
  } else {
    editingId.value = null
    form.value = {
      name: '',
      platform: 'github',
      url: '',
      token: '',
    }
  }
  modalVisible.value = true
}

// 关闭模态框
const closeModal = () => {
  modalVisible.value = false
  formRef.value?.resetFields()
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()

    const formData = {
      name: form.value.name,
      platform: form.value.platform,
      url: form.value.url,
      token: form.value.token,
      defaultBranch: 'main',
      status: 'active',
    }

    if (editingId.value) {
      await updateRepository(editingId.value, formData)
      Message.success('仓库更新成功')
    } else {
      await createRepository(formData)
      Message.success('仓库创建成功')
    }
    closeModal()
    fetchRepositories()
  } catch (error) {
    Message.error('操作失败')
  }
}

// 删除仓库
const handleDelete = async (id) => {
  try {
    await deleteRepository(id)
    Message.success('仓库删除成功')
    await fetchRepositories()
  } catch (error) {
    Message.error('删除失败')
  }
}

// 获取分支列表
const fetchBranches = async (repositoryId) => {
  loadingBranches.value = true
  try {
    const { data } = await getBranches(repositoryId)
    branches.value = data
  } catch (error) {
    Message.error('获取分支列表失败')
  } finally {
    loadingBranches.value = false
  }
}

// 获取提交记录
const fetchCommits = async (repositoryId) => {
  loadingCommits.value = true
  try {
    const { data } = await getCommits(repositoryId)
    commits.value = data
  } catch (error) {
    Message.error('获取提交记录失败')
  } finally {
    loadingCommits.value = false
  }
}

// 显示详情
const showDetail = async (record) => {
  currentRepo.value = record
  detailVisible.value = true
  await Promise.all([
    fetchBranches(record.id),
    fetchCommits(record.id)
  ])
}

// 初始化
onMounted(() => {
  fetchRepositories()
})
</script>

<template>
  <div class="warehouse-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>代码仓库</span>
          <a-button type="primary" @click="showModal()">添加仓库</a-button>
        </div>
      </template>

      <a-table
        :columns="columns"
        :data="repositories"
        :loading="loading"
        :pagination="pagination"
        @page-change="(current) => pagination.current = current"
        @page-size-change="(pageSize) => pagination.pageSize = pageSize"
      >
        <template #action="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showModal(record)">
              编辑
            </a-button>
            <a-button type="text" size="small" @click="showDetail(record)">
              详情
            </a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record.id)">
              删除
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 添加/编辑对话框 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="editingId ? '编辑仓库' : '添加仓库'"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form
        ref="formRef"
        :model="form"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item field="name" label="仓库名称" required>
          <a-input v-model="form.name" placeholder="请输入仓库名称" />
        </a-form-item>
        <a-form-item field="platform" label="平台" required>
          <a-select v-model="form.platform" placeholder="请选择平台">
            <a-option value="github">GitHub</a-option>
            <a-option value="gitlab">GitLab</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="url" label="仓库地址" required>
          <a-input v-model="form.url" placeholder="请输入仓库地址" />
        </a-form-item>
        <a-form-item field="token" label="访问令牌" required>
          <a-input-password v-model="form.token" placeholder="请输入访问令牌" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情对话框 -->
    <a-modal
      v-model:visible="detailVisible"
      title="仓库详情"
      :width="800"
      :footer="false"
    >
      <a-tabs>
        <a-tab-pane key="branches" title="分支">
          <a-table
            :columns="branchColumns"
            :data="branches"
            :loading="loadingBranches"
          />
        </a-tab-pane>
        <a-tab-pane key="commits" title="提交记录">
          <a-table
            :columns="commitColumns"
            :data="commits"
            :loading="loadingCommits"
          >
            <template #commitAction="{ record }">
              <a-button type="text" size="small" @click="handleViewCommit(record)">
                查看
              </a-button>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<style scoped>
.warehouse-container {
  padding: 20px;
}

.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

<template>
  <div class="registry-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>镜像中心</span>
          <a-button type="primary" @click="showModal()">添加镜像仓库</a-button>
        </div>
      </template>

      <a-table
        :columns="columns"
        :data="registries"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
      >
        <template #action="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showModal(record)">
              编辑
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
      :title="currentRegistry ? '编辑镜像仓库' : '添加镜像仓库'"
      @ok="handleSubmit"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item field="name" label="仓库名称" required>
          <a-input v-model="form.name" placeholder="请输入仓库名称" />
        </a-form-item>
        <a-form-item field="type" label="仓库类型" required>
          <a-select v-model="form.type" placeholder="请选择仓库类型">
            <a-option value="public">公有仓库</a-option>
            <a-option value="private">私有仓库</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="url" label="仓库地址" required>
          <a-input v-model="form.url" placeholder="请输入仓库地址" />
        </a-form-item>
        <a-form-item field="username" label="用户名">
          <a-input v-model="form.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleTestConnection">
            测试连接
          </a-button>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { TableColumnData } from '@arco-design/web-vue';
import { getDockerRegistries, createDockerRegistry, updateDockerRegistry, deleteDockerRegistry, testDockerRegistryConnection } from '@/api/registry';

// 定义 Docker Registry 类型
interface DockerRegistry {
  id: number;
  name: string;
  type: 'public' | 'private';
  url: string;
  username: string;
  password: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

// 表格列定义
const columns: TableColumnData[] = [
  {
    title: '仓库名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '仓库类型',
    dataIndex: 'type',
    key: 'type',
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
    width: 200,
    fixed: 'right',
    slotName: 'action',
  },
];

// 数据
const registries = ref<DockerRegistry[]>([]);
const loading = ref(false);
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
});

// 表单
const modalVisible = ref(false);
const currentRegistry = ref<DockerRegistry | null>(null);
const form = reactive({
  name: '',
  type: 'public' as const,
  url: '',
  username: '',
  password: '',
});

// 获取镜像仓库列表
const fetchRegistries = async () => {
  loading.value = true;
  try {
    const res = await getDockerRegistries({
      current: pagination.current,
      pageSize: pagination.pageSize,
    });
    registries.value = res.data.list;
    pagination.total = res.data.total;
  } catch (error) {
    Message.error('获取镜像仓库列表失败');
  } finally {
    loading.value = false;
  }
};

// 表格分页变化
const handleTableChange = (page: number, pageSize: number) => {
  pagination.current = page;
  pagination.pageSize = pageSize;
  fetchRegistries();
};

// 显示对话框
const showModal = (record?: any) => {
  currentRegistry.value = record;
  if (record) {
    Object.assign(form, record);
  } else {
    Object.assign(form, {
      name: '',
      type: 'public',
      url: '',
      username: '',
      password: '',
    });
  }
  modalVisible.value = true;
};

// 提交表单
const handleSubmit = async () => {
  try {
    if (currentRegistry.value) {
      await updateDockerRegistry(currentRegistry.value.id, form);
      Message.success('更新成功');
    } else {
      await createDockerRegistry(form);
      Message.success('创建成功');
    }
    modalVisible.value = false;
    fetchRegistries();
  } catch (error) {
    Message.error('操作失败');
  }
};

// 删除镜像仓库
const handleDelete = async (id: number) => {
  try {
    await deleteDockerRegistry(id);
    Message.success('删除成功');
    fetchRegistries();
  } catch (error) {
    Message.error('删除失败');
  }
};

// 测试连接
const handleTestConnection = async () => {
  try {
    const res = await testDockerRegistryConnection(form);
    if (res.data.status === 'success') {
      Message.success('连接成功');
    } else {
      Message.error('连接失败');
    }
  } catch (error) {
    Message.error('连接失败');
  }
};

// 初始化
fetchRegistries();
</script>

<style scoped>
.registry-container {
  padding: 20px;
}

.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

<script lang="ts" setup>
  import { ref, reactive, onMounted } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import {
    IconPlus,
    IconSync,
    IconEye,
    IconDelete,
    IconEdit,
  } from '@arco-design/web-vue/es/icon';
  import { useI18n } from 'vue-i18n';
  import dayjs from 'dayjs';
  import {
    getRepositories,
    createRepository,
    updateRepository,
    deleteRepository,
    testRepository,
  } from '@/api/warehouse';
  import type { Repository } from '@/api/warehouse';

  const { t } = useI18n();

  const loading = ref(false);
  const repositories = ref<Repository[]>([]);
  const modalVisible = ref(false);
  const editingId = ref<number | null>(null);
  const formRef = ref();

  const form = ref({
    name: '',
    url: '',
    token: '',
  });

  const rules = {
    name: [{ required: true, message: '请输入仓库名称' }],
    url: [{ required: true, message: '请输入仓库地址' }],
    token: [{ required: true, message: '请输入访问令牌' }],
  };

  // 获取仓库列表
  const fetchRepositories = async () => {
    loading.value = true;
    try {
      console.log('获取仓库列表');
      const res = await getRepositories();
      console.log(res,"111")
      // repositories.value
      console.log('获取到的仓库列表:', repositories.value);
      repositories.value = res
    } catch (error) {
      console.error('获取仓库列表失败:', error);
      Message.error('获取仓库列表失败');
    } finally {
      loading.value = false;
    }
  };

  // 打开模态框
  const openModal = (record?: Repository) => {
    if (record) {
      editingId.value = record.id;
      form.value = { ...record };
    } else {
      editingId.value = null;
      form.value = {
        name: '',
        url: '',
        token: '',
      };
    }
    modalVisible.value = true;
  };

  // 关闭模态框
  const closeModal = () => {
    modalVisible.value = false;
    formRef.value?.resetFields();
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const valid = await formRef.value?.validate();
      // if (!valid) {
      //   console.log('表单验证失败');
      //   return;
      // }

      const formData = {
        name: form.value.name,
        platform: 'github', // 默认使用 GitHub
        url: form.value.url,
        token: form.value.token,
        defaultBranch: 'main', // 默认分支
        status: 'active', // 默认状态
      };

      console.log('准备提交数据:', formData);

      if (editingId.value) {
        console.log('更新仓库:', editingId.value);
        await updateRepository(editingId.value, formData);
        Message.success('仓库更新成功');
      } else {
        console.log('创建新仓库');
        await createRepository(formData);
        Message.success('仓库创建成功');
      }
      closeModal();
      fetchRepositories();
    } catch (error) {
      console.error('提交失败:', error);
      Message.error('操作失败');
    }
  };

  // 删除仓库
  const handleDelete = async (id: number) => {
    try {
      await deleteRepository(id);
      Message.success('仓库删除成功');
      fetchRepositories();
    } catch (error) {
      Message.error('删除失败');
    }
  };

  // 获取平台颜色
  const getPlatformColor = (platform: Repository['platform']) => {
    const colors: Record<Repository['platform'], string> = {
      github: 'blue',
      gitlab: 'orange',
      gitee: 'red',
    };
    return colors[platform] || 'gray';
  };

  // 获取状态颜色
  const getStatusColor = (status: Repository['status']) => {
    const colors: Record<Repository['status'], string> = {
      active: 'green',
      inactive: 'gray',
      error: 'red',
    };
    return colors[status] || 'gray';
  };

  // 格式化日期
  const formatDate = (date: string) => {
    return new Date(date).toLocaleString();
  };

  // 测试仓库连接
  const testConnection = async (url: string, token: string) => {
    try {
      loading.value = true;
      const result = await testRepository({ url, token });
      if (result.success) {
        Message.success('仓库连接测试成功');
      } else {
        Message.error(result.message || '仓库连接测试失败');
      }
    } catch (error) {
      console.error('测试连接失败:', error);
      Message.error('测试连接失败');
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    fetchRepositories();
  });
</script>

<template>
  <div class="warehouse-container">
    <a-card class="general-card" :title="$t('仓库管理')">
      <template #extra>
        <a-button type="primary" @click="openModal()">
          <template #icon><icon-plus /></template>
          {{ $t('添加仓库') }}
        </a-button>
      </template>
      <a-table
        :data="repositories"
        :loading="loading"
        :pagination="false"
        :bordered="false"
      >
        <template #columns>
          <a-table-column :title="$t('仓库名称')" data-index="name" />
          <a-table-column :title="$t('仓库地址')" data-index="url" />
          <a-table-column :title="$t('默认分支')" data-index="defaultBranch" />
          <a-table-column :title="$t('最后同步时间')" data-index="lastSyncAt">
            <template #cell="{ record }">
              {{ record.lastSyncAt ? formatDate(record.lastSyncAt) : '-' }}
            </template>
          </a-table-column>
          <a-table-column :title="$t('状态')" data-index="status" />
          <a-table-column :title="$t('操作')" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openModal(record)">
                  <template #icon><icon-edit /></template>
                  {{ $t('编辑') }}
                </a-button>
                <a-popconfirm
                  :content="$t('确定要删除这个仓库吗？')"
                  @ok="handleDelete(record.id)"
                >
                  <a-button type="text" status="danger" size="small">
                    <template #icon><icon-delete /></template>
                    {{ $t('删除') }}
                  </a-button>
                </a-popconfirm>
                <a-button 
                  type="text" 
                  size="small" 
                  status="success"
                  @click="testConnection(record.url, record.token)"
                >
                  {{ $t('测试连接') }}
                </a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="modalVisible"
      :title="editingId ? $t('编辑仓库') : $t('添加仓库')"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
        <a-form-item field="name" :label="$t('仓库名称')">
          <a-input v-model="form.name" :placeholder="$t('请输入仓库名称')" />
        </a-form-item>
        <a-form-item field="url" :label="$t('仓库地址')">
          <a-input v-model="form.url" :placeholder="$t('请输入仓库地址')" />
        </a-form-item>
        <a-form-item field="token" :label="$t('访问令牌')">
          <a-input-password
            v-model="form.token"
            :placeholder="$t('请输入访问令牌')"
          />
        </a-form-item>
        <div class="form-footer">
          <a-button type="primary" @click="handleSubmit">
            {{ editingId ? $t('更新') : $t('创建') }}
          </a-button>
          <a-button 
            type="outline" 
            status="success" 
            @click="testConnection(form.url, form.token)"
            style="margin-left: 8px"
          >
            {{ $t('测试连接') }}
          </a-button>
        </div>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
  .warehouse-container {
    padding: 20px;
  }

  .general-card {
    margin-bottom: 20px;
  }

  .form-footer {
    text-align: right;
  }
</style>

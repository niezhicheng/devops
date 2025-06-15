<template>
  <div class="table-container">
    <a-card>
      <template #title>
        <div class="card-header">
          <span>用户列表</span>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            添加用户
          </a-button>
        </div>
      </template>

      <a-table
        :columns="columns"
        :data="tableData"
        :pagination="pagination"
        @page-change="onPageChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'red'">
            {{ record.status === 'active' ? '启用' : '禁用' }}
          </a-tag>
        </template>

        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">
              编辑
            </a-button>
            <a-button type="text" status="danger" size="small" @click="handleDelete(record)">
              删除
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 添加/编辑用户对话框 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="modalTitle"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form :model="formData" ref="formRef">
        <a-form-item field="name" label="姓名" :rules="[{ required: true, message: '请输入姓名' }]">
          <a-input v-model="formData.name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item field="email" label="邮箱" :rules="[{ required: true, message: '请输入邮箱' }]">
          <a-input v-model="formData.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-select v-model="formData.status">
            <a-option value="active">启用</a-option>
            <a-option value="inactive">禁用</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
export default {
  name: 'TableDemo',
  data() {
    return {
      columns: [
        {
          title: '姓名',
          dataIndex: 'name'
        },
        {
          title: '邮箱',
          dataIndex: 'email'
        },
        {
          title: '状态',
          dataIndex: 'status',
          slotName: 'status'
        },
        {
          title: '操作',
          slotName: 'operations'
        }
      ],
      tableData: [
        {
          id: 1,
          name: '张三',
          email: 'zhangsan@example.com',
          status: 'active'
        },
        {
          id: 2,
          name: '李四',
          email: 'lisi@example.com',
          status: 'inactive'
        },
        {
          id: 3,
          name: '王五',
          email: 'wangwu@example.com',
          status: 'active'
        }
      ],
      pagination: {
        total: 3,
        current: 1,
        pageSize: 10
      },
      modalVisible: false,
      modalTitle: '添加用户',
      formData: {
        name: '',
        email: '',
        status: 'active'
      },
      isEdit: false,
      currentRecord: null
    }
  },
  methods: {
    handleAdd() {
      this.modalTitle = '添加用户'
      this.isEdit = false
      this.formData = {
        name: '',
        email: '',
        status: 'active'
      }
      this.modalVisible = true
    },
    handleEdit(record) {
      this.modalTitle = '编辑用户'
      this.isEdit = true
      this.currentRecord = record
      this.formData = { ...record }
      this.modalVisible = true
    },
    handleDelete(record) {
      this.$modal.confirm({
        title: '确认删除',
        content: `确定要删除用户 ${record.name} 吗？`,
        onOk: () => {
          const index = this.tableData.findIndex(item => item.id === record.id)
          if (index > -1) {
            this.tableData.splice(index, 1)
            this.$message.success('删除成功')
          }
        }
      })
    },
    handleModalOk() {
      this.$refs.formRef.validate((errors) => {
        if (!errors) {
          if (this.isEdit) {
            const index = this.tableData.findIndex(item => item.id === this.currentRecord.id)
            if (index > -1) {
              this.tableData[index] = { ...this.currentRecord, ...this.formData }
            }
          } else {
            const newId = Math.max(...this.tableData.map(item => item.id)) + 1
            this.tableData.push({
              id: newId,
              ...this.formData
            })
          }
          this.modalVisible = false
          this.$message.success(this.isEdit ? '编辑成功' : '添加成功')
        }
      })
    },
    handleModalCancel() {
      this.modalVisible = false
    },
    onPageChange(page) {
      this.pagination.current = page
    }
  }
}
</script>

<style scoped>
.table-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

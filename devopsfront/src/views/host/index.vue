<template>
  <div class="container">
    <Breadcrumb :items="['菜单.列表', '主机管理']" />
    <a-card class="general-card" title="查询表格">
      <a-row>
        <a-col :flex="1">
          <a-form
            :model="formModel"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
            label-align="left"
          >
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item field="status" label="状态">
                  <a-select
                    v-model="formModel.status"
                    :options="statusOptions"
                    placeholder="请选择状态"
                  />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="search">
              <template #icon>
                <icon-search />
              </template>
              查询
            </a-button>
            <a-button @click="reset">
              <template #icon>
                <icon-refresh />
              </template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin-top: 0" />
      <a-row style="margin-bottom: 16px">
        <a-col :span="12">
          <a-space>
            <a-button type="primary" @click="openAddModal">
              <template #icon>
                <icon-plus />
              </template>
              新增
            </a-button>
            <a-upload action="/">
              <template #upload-button>
                <a-button> 导入 </a-button>
              </template>
            </a-upload>
          </a-space>
        </a-col>
        <a-col
          :span="12"
          style="display: flex; align-items: center; justify-content: end"
        >
          <a-button>
            <template #icon>
              <icon-download />
            </template>
            下载
          </a-button>
          <a-tooltip content="刷新">
            <div class="action-icon" @click="search">
              <icon-refresh size="18" />
            </div>
          </a-tooltip>
          <a-dropdown @select="handleSelectDensity">
            <a-tooltip content="表格密度">
              <div class="action-icon"><icon-line-height size="18" /></div>
            </a-tooltip>
            <template #content>
              <a-doption
                v-for="item in densityList"
                :key="item.value"
                :value="item.value"
                :class="{ active: item.value === size }"
              >
                <span>{{ item.name }}</span>
              </a-doption>
            </template>
          </a-dropdown>
          <a-tooltip content="列设置">
            <a-popover
              trigger="click"
              position="bl"
              @popup-visible-change="popupVisibleChange"
            >
              <div class="action-icon"><icon-settings size="18" /></div>
              <template #content>
                <div id="tableSetting">
                  <div
                    v-for="(item, index) in showColumns"
                    :key="item.dataIndex"
                    class="setting"
                  >
                    <div style="margin-right: 4px; cursor: move">
                      <icon-drag-arrow />
                    </div>
                    <div>
                      <a-checkbox
                        v-model="item.checked"
                        @change="
                          handleChange($event, item, index)
                        "
                      >
                      </a-checkbox>
                    </div>
                    <div class="title">
                      {{ item.title === '#' ? '序列号' : item.title }}
                    </div>
                  </div>
                </div>
              </template>
            </a-popover>
          </a-tooltip>
        </a-col>
      </a-row>
      <a-table
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        :columns="cloneColumns"
        :data="renderData"
        :bordered="false"
        :size="size"
        @page-change="onPageChange"
      >
        <template #index="{ rowIndex }">
          {{ rowIndex + 1 + (pagination.current - 1) * pagination.pageSize }}
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleView(record)">
              查看
            </a-button>
            <a-button type="text" size="small" @click="handleSftp(record)">
              <template #icon><icon-folder /></template>
              SFTP
            </a-button>
            <a-button type="text" size="small" @click="handleWebShell(record)">
              <template #icon><icon-code /></template>
              WebShell
            </a-button>
            <a-button
              type="text"
              size="small"
              status="danger"
              @click="handleDelete(record)"
            >
              删除
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增表单对话框 -->
    <a-modal
      v-model:visible="addModalVisible"
      title="新增主机"
      :mask-closable="false"
      :unmount-on-close="true"
      :ok-button-props="{ disabled: false }"
      :cancel-button-props="{ disabled: false }"
      @ok="handleAdd"
      @cancel="closeAddModal"
    >
      <a-form
        ref="addFormRef"
        :model="addForm"
        :rules="addFormRules"
        label-align="right"
        :label-col-props="{ span: 6 }"
        :wrapper-col-props="{ span: 18 }"
        @submit.prevent="handleAdd"
      >
        <a-form-item field="name" label="主机名称" validate-trigger="blur">
          <a-input
            v-model="addForm.name"
            placeholder="请输入主机名称"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="ip" label="IP地址" validate-trigger="blur">
          <a-input
            v-model="addForm.ip"
            placeholder="请输入IP地址"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="port" label="端口" validate-trigger="blur">
          <a-input-number
            v-model="addForm.port"
            placeholder="请输入端口"
            :min="1"
            :max="65535"
          />
        </a-form-item>
        <a-form-item field="username" label="用户名" validate-trigger="blur">
          <a-input
            v-model="addForm.username"
            placeholder="请输入用户名"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="password" label="密码" validate-trigger="blur">
          <a-input-password
            v-model="addForm.password"
            placeholder="请输入密码"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="description" label="描述">
          <a-textarea
            v-model="addForm.description"
            placeholder="请输入描述信息"
            allow-clear
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 删除确认对话框 -->
    <a-modal
      v-model:visible="deleteModalVisible"
      title="删除确认"
      :mask-closable="false"
      @ok="confirmDelete"
      @cancel="cancelDelete"
    >
      <div
      >确定要删除主机 "{{ currentRecord?.name || currentRecord?.ip }}"
        吗？</div
      >
    </a-modal>

    <!-- SFTP对话框 -->
    <a-modal
      v-model:visible="sftpModalVisible"
      title="SFTP文件管理"
      :width="800"
      :mask-closable="false"
      :footer="false"
    >
      <div class="sftp-container">
        <div class="sftp-header">
          <a-space>
            <a-button type="primary" @click="handleSftpRefresh">
              <template #icon><icon-refresh /></template>
              刷新
            </a-button>
            <a-button @click="handleUploadClick">
              <template #icon><icon-upload /></template>
              上传
            </a-button>
          </a-space>
          <a-input-search
            v-model="sftpSearchPath"
            placeholder="输入路径"
            style="width: 300px"
            @search="handleSftpPathSearch"
          />
        </div>
        <div class="sftp-breadcrumb">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <a-button type="text" @click="handleBreadcrumbClick('/')">
                <icon-home />
              </a-button>
            </a-breadcrumb-item>
            <template v-for="(item, index) in breadcrumbItems" :key="index">
              <a-breadcrumb-item>
                <a-button type="text" @click="handleBreadcrumbClick(item.path)">
                  {{ item.name }}
                </a-button>
              </a-breadcrumb-item>
            </template>
          </a-breadcrumb>
        </div>
        <div class="sftp-content">
          <a-table
            :data="sftpFileList"
            :loading="sftpLoading"
            :pagination="false"
            :bordered="false"
            :scroll="{ x: 800, y: 500 }"
            class="sftp-table"
          >
            <template #columns>
              <a-table-column title="名称" data-index="name" :width="300">
                <template #cell="{ record }">
                  <a-space>
                    <icon-folder v-if="record.type === 'directory'" />
                    <icon-file v-else />
                    <a-button type="text" @click="handleSftpItemClick(record)">
                      {{ record.name }}
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
              <a-table-column title="大小" data-index="size" :width="120">
                <template #cell="{ record }">
                  {{
                    record.type === 'directory'
                      ? '-'
                      : formatFileSize(record.size)
                  }}
                </template>
              </a-table-column>
              <a-table-column
                title="修改时间"
                data-index="modifyTime"
                :width="180"
              />
              <a-table-column
                title="权限"
                data-index="permissions"
                :width="120"
              />
              <a-table-column title="操作" :width="180" fixed="right">
                <template #cell="{ record }">
                  <a-space>
                    <a-button
                      v-if="record.type === 'file'"
                      type="text"
                      size="mini"
                      @click="handleSftpDownload(record)"
                    >
                      <template #icon><icon-download /></template>
                    </a-button>
                    <a-button
                      v-if="record.type === 'directory'"
                      type="text"
                      size="mini"
                      @click="handleSftpCompress(record)"
                    >
                      <template #icon><icon-zip /></template>
                    </a-button>
                    <a-button
                      type="text"
                      size="mini"
                      @click="handleSftpRename(record)"
                    >
                      <template #icon><icon-edit /></template>
                    </a-button>
                    <a-button
                      type="text"
                      size="mini"
                      status="danger"
                      @click="handleSftpFileDelete(record)"
                    >
                      <template #icon><icon-delete /></template>
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
            </template>
          </a-table>
        </div>
      </div>
    </a-modal>

    <!-- 重命名对话框 -->
    <a-modal
      v-model:visible="renameModalVisible"
      title="重命名"
      @ok="handleRenameConfirm"
      @cancel="closeRenameModal"
    >
      <a-form :model="renameForm" layout="vertical">
        <a-form-item field="newName" label="新名称">
          <a-input v-model="renameForm.newName" placeholder="请输入新名称" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 文件上传对话框 -->
    <a-modal
      v-model:visible="uploadModalVisible"
      title="上传文件"
      @cancel="handleUploadCancel"
      @before-ok="handleUploadConfirm"
    >
      <a-upload
        :custom-request="handleSftpUpload"
        :show-file-list="false"
        accept="*"
        :limit="1"
        :auto-upload="true"
      >
        <template #upload-button>
          <a-button type="primary">
            <template #icon><icon-upload /></template>
            选择文件
          </a-button>
        </template>
      </a-upload>
    </a-modal>

    <!-- WebShell模态框 -->
    <a-modal
      v-model:visible="webShellVisible"
      title="WebShell"
      :footer="false"
      :mask-closable="false"
      :unmount-on-close="true"
      :width="'90vw'"
      :height="'90vh'"
      @close="handleWebShellClose"
    >
      <div class="web-shell-container">
        <div ref="terminalRef" class="terminal"></div>
      </div>
    </a-modal>
  </div>
</template>

<script>
import { computed, ref, reactive, watch, nextTick } from 'vue';
import useLoading from '@/hooks/loading';
import { Message, Modal } from '@arco-design/web-vue';
import { queryHostList, addHost, deleteHost, uploadSftpFile, fetchSftpFiles, deleteSftpFile, downloadSftpFile, renameSftpFile, compressSftpDir } from '@/api/host';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SearchAddon } from 'xterm-addon-search';
import 'xterm/css/xterm.css';
// 添加中文文本映射
const contentTypeText = {
  img: '图文',
  horizontalVideo: '横版短视频',
  verticalVideo: '竖版短视频',
};

const filterTypeText = {
  artificial: '人工筛选',
  rules: '规则筛选',
};

const statusText = {
  online: '已上线',
  offline: '已下线',
};

const generateFormModel = () => {
  return {
    number: '',
    name: '',
    contentType: '',
    filterType: '',
    createdTime: [],
    status: '',
  };
};

export default {
  setup() {
    const { loading, setLoading } = useLoading(true);
    const renderData = ref([]);
    const formModel = ref(generateFormModel());
    const cloneColumns = ref([]);
    const showColumns = ref([]);

    const size = ref('medium');

    const basePagination = reactive({
      current: 1,
      pageSize: 20,
    });
    const pagination = reactive({
      ...basePagination,
    });
    const densityList = computed(() => [
      {
        name: '迷你',
        value: 'mini',
      },
      {
        name: '偏小',
        value: 'small',
      },
      {
        name: '中等',
        value: 'medium',
      },
      {
        name: '偏大',
        value: 'large',
      },
    ]);
    const columns = computed(() => [
      {
        title: '序号',
        dataIndex: 'index',
        slotName: 'index',
      },
      {
        title: '主机名称',
        dataIndex: 'name',
      },
      {
        title: 'IP地址',
        dataIndex: 'ip',
      },
      {
        title: '端口',
        dataIndex: 'port',
      },
      {
        title: '用户名',
        dataIndex: 'username',
      },
      {
        title: '密码',
        dataIndex: 'password',
      },
      {
        title: '创建时间',
        dataIndex: 'createdTime',
      },
      {
        title: '操作',
        dataIndex: 'operations',
        slotName: 'operations',
      },
    ]);
    const statusOptions = computed(() => [
      {
        label: '已上线',
        value: 'online',
      },
      {
        label: '已下线',
        value: 'offline',
      },
    ]);

    const fetchData = async (params = { current: 1, pageSize: 20 }) => {
      setLoading(true);
      try {
        console.log('发送请求参数:', params);
        const { data } = await queryHostList(params);
        console.log('收到响应数据:', data);
        renderData.value = data.list;
        pagination.current = params.current;
        pagination.total = data.total;
      } catch (err) {
        console.error('获取数据失败:', err);
        Message.error('获取数据失败');
      } finally {
        setLoading(false);
      }
    };

    const search = () => {
      fetchData({
        ...basePagination,
        ...formModel.value,
      });
    };

    const onPageChange = (current) => {
      fetchData({ ...basePagination, current });
    };

    // 初始化加载数据
    fetchData();
    const reset = () => {
      formModel.value = generateFormModel();
    };

    const handleSelectDensity = (val) => {
      size.value = val;
    };

    const handleChange = (checked, column, index) => {
      if (!checked) {
        cloneColumns.value = showColumns.value.filter(
          (item) => item.dataIndex !== column.dataIndex
        );
      } else {
        cloneColumns.value.splice(index, 0, column);
      }
    };

    const exchangeArray = (array, beforeIdx, newIdx, isDeep = false) => {
      const newArray = isDeep ? JSON.parse(JSON.stringify(array)) : array;
      if (beforeIdx > -1 && newIdx > -1) {
        // 先替换后面的，然后拿到替换的结果替换前面的
        newArray.splice(
          beforeIdx,
          1,
          newArray.splice(newIdx, 1, newArray[beforeIdx]).pop()
        );
      }
      return newArray;
    };

    const popupVisibleChange = (val) => {
      if (val) {
        nextTick(() => {
          const el = document.getElementById('tableSetting');
          const sortable = new Sortable(el, {
            onEnd(e) {
              const { oldIndex, newIndex } = e;
              exchangeArray(cloneColumns.value, oldIndex, newIndex);
              exchangeArray(showColumns.value, oldIndex, newIndex);
            },
          });
        });
      }
    };

    watch(
      () => columns.value,
      (val) => {
        cloneColumns.value = JSON.parse(JSON.stringify(val));
        cloneColumns.value.forEach((item, index) => {
          item.checked = true;
        });
        showColumns.value = JSON.parse(JSON.stringify(cloneColumns.value));
      },
      { deep: true, immediate: true }
    );

    // 新增表单相关
    const addModalVisible = ref(false);
    const addFormRef = ref();
    const addForm = ref({
      name: '',
      ip: '',
      port: 22,
      username: '',
      password: '',
      description: '',
    });

    const addFormRules = {
      name: [
        { required: true, message: '请输入主机名称' },
        { minLength: 2, message: '主机名称至少2个字符' },
      ],
      ip: [
        { required: true, message: '请输入IP地址' },
        {
          match: /^(\d{1,3}\.){3}\d{1,3}$/,
          message: 'IP地址格式不正确',
        },
      ],
      port: [
        { required: true, message: '请输入端口' },
        {
          type: 'number',
          min: 1,
          max: 65535,
          message: '端口号必须在1-65535之间',
        },
      ],
      username: [
        { required: true, message: '请输入用户名' },
        { minLength: 2, message: '用户名至少2个字符' },
      ],
      password: [
        { required: true, message: '请输入密码' },
        { minLength: 2, message: '密码至少2个字符' },
      ],
    };

    const openAddModal = () => {
      addModalVisible.value = true;
    };

    const closeAddModal = () => {
      addForm.value = {
        name: '',
        ip: '',
        port: 22,
        username: '',
        password: '',
        description: '',
      };
      addFormRef.value?.resetFields();
      addModalVisible.value = false;
    };

    const handleAdd = async () => {
      try {
        if (
          !addForm.value.name ||
          !addForm.value.ip ||
          !addForm.value.username ||
          !addForm.value.password
        ) {
          Message.error('请填写必填字段');
          return;
        }

        // 验证IP地址格式
        const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/;
        if (!ipRegex.test(addForm.value.ip)) {
          Message.error('IP地址格式不正确');
          return;
        }

        // 验证端口范围
        if (addForm.value.port < 1 || addForm.value.port > 65535) {
          Message.error('端口号必须在1-65535之间');
          return;
        }

        // 验证用户名长度
        if (addForm.value.username.length < 1) {
          Message.error('用户名至少1个字符');
          return;
        }

        // 验证密码长度
        if (addForm.value.password.length < 1) {
          Message.error('密码至少1个字符');
          return;
        }

        await addHost(addForm.value);
        Message.success('添加成功');
        closeAddModal();
        search(); // 刷新列表
      } catch (error) {
        console.error('添加失败:', error);
        Message.error('添加失败');
      }
    };

    // 删除相关
    const deleteModalVisible = ref(false);
    const currentRecord = ref(null);

    const handleDelete = (record) => {
      currentRecord.value = record;
      deleteModalVisible.value = true;
    };

    const confirmDelete = async () => {
      if (!currentRecord.value) return;

      try {
        await deleteHost(currentRecord.value.id);
        Message.success('删除成功');
        deleteModalVisible.value = false;
        search(); // 刷新列表
      } catch (error) {
        console.error('删除失败:', error);
        Message.error('删除失败');
      }
    };

    const cancelDelete = () => {
      deleteModalVisible.value = false;
      currentRecord.value = null;
    };

    const handleView = (record) => {
      console.log('查看详情:', record);
    };

    // SFTP相关
    const sftpModalVisible = ref(false);
    const sftpLoading = ref(false);
    const sftpFileList = ref([]);
    const sftpSearchPath = ref('');
    const currentSftpHost = ref(null);
    const uploadModalVisible = ref(false);
    const breadcrumbItems = ref([]);
    const renameModalVisible = ref(false);
    const renameForm = ref({
      newName: '',
      currentFile: null,
    });

    const handleSftp = (record) => {
      currentSftpHost.value = record;
      sftpModalVisible.value = true;
      loadSftpFiles('/');
    };

    // 更新面包屑
    const updateBreadcrumb = (path) => {
      if (path === '/') {
        breadcrumbItems.value = [];
        return;
      }

      const parts = path.split('/').filter(Boolean);
      const items = [];
      let currentPath = '';

      parts.forEach((part) => {
        currentPath += `/${part}`;
        items.push({
          name: part,
          path: currentPath,
        });
      });

      breadcrumbItems.value = items;
    };

    const handleBreadcrumbClick = (path) => {
      loadSftpFiles(path);
    };

    const loadSftpFiles = async (path) => {
      if (!currentSftpHost.value) return;

      sftpLoading.value = true;
      try {
        const { data } = await fetchSftpFiles(currentSftpHost.value.id, path);
        console.log(data, '这是数据data');
        sftpFileList.value = data.list;
        sftpSearchPath.value = path;
        updateBreadcrumb(path);
      } catch (error) {
        console.error('加载SFTP文件失败:', error);
        Message.error('加载文件列表失败');
      } finally {
        sftpLoading.value = false;
      }
    };

    const handleSftpRefresh = () => {
      loadSftpFiles(sftpSearchPath.value);
    };

    const handleSftpPathSearch = () => {
      loadSftpFiles(sftpSearchPath.value);
    };

    const handleSftpItemClick = (record) => {
      if (record.type === 'directory') {
        loadSftpFiles(record.path);
      }
    };

    const handleSftpUpload = async (option) => {
      console.log('Upload option:', option);
      if (!currentSftpHost.value) return;

      try {
        const formData = new FormData();
        const { file } = option.fileItem;
        console.log('File to upload:', file);

        formData.append('file', file);
        formData.append('path', sftpSearchPath.value || '/');

        Message.loading({ content: '正在上传文件...', duration: 0 });
        await uploadSftpFile(currentSftpHost.value.id, formData);
        Message.clear();
        Message.success('上传成功');
        loadSftpFiles(sftpSearchPath.value);
        uploadModalVisible.value = false;
      } catch (error) {
        Message.clear();
        console.error('上传失败:', error);
        Message.error('上传失败');
      }
    };

    const handleSftpDownload = async (record) => {
      if (!currentSftpHost.value) return;

      try {
        const response = await downloadSftpFile(
          currentSftpHost.value.id,
          record.path
        );
        const blob = new Blob([response.data]);
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = record.name;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
        Message.success('下载成功');
      } catch (error) {
        console.error('下载失败:', error);
        Message.error('下载失败');
      }
    };

    const handleSftpCompress = async (record) => {
      if (!currentSftpHost.value) return;

      try {
        Message.loading({ content: '正在压缩文件夹...', duration: 0 });
        await compressSftpDir(currentSftpHost.value.id, record.path);
        Message.clear();
        Message.success('压缩成功');
        loadSftpFiles(sftpSearchPath.value);
      } catch (error) {
        Message.clear();
        console.error('压缩失败:', error);
        Message.error('压缩失败');
      }
    };

    const handleSftpRename = (record) => {
      renameForm.value.currentFile = record;
      renameForm.value.newName = record.name;
      renameModalVisible.value = true;
    };

    const handleRenameConfirm = async () => {
      if (!currentSftpHost.value || !renameForm.value.currentFile) return;

      const oldPath = renameForm.value.currentFile.path;
      const newPath = oldPath.replace(
        renameForm.value.currentFile.name,
        renameForm.value.newName
      );

      try {
        await renameSftpFile(currentSftpHost.value.id, oldPath, newPath);
        Message.success('重命名成功');
        renameModalVisible.value = false;
        await loadSftpFiles(sftpSearchPath.value);
      } catch (error) {
        console.error('重命名失败:', error);
        Message.error('重命名失败');
      }
    };

    const closeRenameModal = () => {
      renameModalVisible.value = false;
      renameForm.value = {
        newName: '',
        currentFile: null,
      };
    };

    const handleSftpFileDelete = async (record) => {
      if (!currentSftpHost.value) return;

      Modal.warning({
        title: '确认删除',
        content: `确定要删除 ${record.name} 吗？`,
        okText: '确认',
        cancelText: '取消',
        onOk: async () => {
          try {
            await deleteSftpFile(currentSftpHost.value.id, record.path);
            Message.success('删除成功');
            loadSftpFiles(sftpSearchPath.value);
          } catch (error) {
            console.error('删除失败:', error);
            Message.error('删除失败');
          }
        },
      });
    };

    const handleUploadClick = () => {
      uploadModalVisible.value = true;
    };

    const handleUploadCancel = () => {
      uploadModalVisible.value = false;
    };

    const handleUploadConfirm = () => {
      uploadModalVisible.value = false;
    };

    const formatFileSize = (size) => {
      if (size < 1024) return `${size} B`;
      if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
      if (size < 1024 * 1024 * 1024)
        return `${(size / (1024 * 1024)).toFixed(2)} MB`;
      return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
    };

    // WebShell相关状态
    const webShellVisible = ref(false);
    const terminalRef = ref(null);
    const terminal = ref(null);
    const wsConnection = ref(null);

    // 初始化终端
    const initTerminal = () => {
      if (!terminalRef.value) return;

      terminal.value = new Terminal({
        cursorBlink: true,
        fontSize: 14,
        fontFamily: 'Menlo, Monaco, "Courier New", monospace',
        theme: {
          background: '#1e1e1e',
          foreground: '#ffffff',
        },
        cols: 200,
        rows: 40,
      });

      const fitAddon = new FitAddon();
      const webLinksAddon = new WebLinksAddon();
      const searchAddon = new SearchAddon();

      terminal.value.loadAddon(fitAddon);
      terminal.value.loadAddon(webLinksAddon);
      terminal.value.loadAddon(searchAddon);

      terminal.value.open(terminalRef.value);

      nextTick(() => {
        fitAddon.fit();
        if (wsConnection.value?.readyState === WebSocket.OPEN) {
          const { rows, cols } = terminal.value;
          const resizeMessage = `\x1b[8;${rows};${cols}t`;
          console.log(
            'Sending resize message:',
            resizeMessage,
            'rows:',
            rows,
            'cols:',
            cols
          );
          wsConnection.value.send(resizeMessage);
        }
      });

      const resizeObserver = new ResizeObserver(() => {
        fitAddon.fit();
        if (wsConnection.value?.readyState === WebSocket.OPEN) {
          const { rows, cols } = terminal.value;
          const resizeMessage = `\x1b[8;${rows};${cols}t`;
          console.log(
            'Sending resize message:',
            resizeMessage,
            'rows:',
            rows,
            'cols:',
            cols
          );
          wsConnection.value.send(resizeMessage);
        }
      });

      resizeObserver.observe(terminalRef.value);

      terminal.value.onData((data) => {
        if (wsConnection.value?.readyState === WebSocket.OPEN) {
          setTimeout(() => {
            wsConnection.value?.send(data);
          }, 0);
        }
      });
    };

    const connectWebSocket = () => {
      if (!currentSftpHost.value || !terminal.value) return;

      const wsUrl = `ws://localhost:8080/api/host/${currentSftpHost.value.id}/webshell`;
      wsConnection.value = new WebSocket(wsUrl);

      wsConnection.value.onopen = () => {
        terminal.value?.writeln('WebShell连接已建立');
      };

      wsConnection.value.onmessage = (event) => {
        setTimeout(() => {
          terminal.value?.write(event.data);
        }, 0);
      };

      wsConnection.value.onerror = (error) => {
        terminal.value?.writeln(`\r\n\x1b[31m连接错误: ${error}\x1b[0m`);
      };

      wsConnection.value.onclose = () => {
        terminal.value?.writeln('\r\n\x1b[31m连接已关闭\x1b[0m');
      };
    };

    const handleWebShell = (record) => {
      currentSftpHost.value = record;
      webShellVisible.value = true;
      nextTick(() => {
        initTerminal();
        connectWebSocket();
      });
    };

    const handleWebShellClose = () => {
      if (wsConnection.value) {
        wsConnection.value.close();
        wsConnection.value = null;
      }

      if (terminal.value) {
        terminal.value.dispose();
        terminal.value = null;
      }

      webShellVisible.value = false;
    };

    return {
      loading,
      renderData,
      formModel,
      cloneColumns,
      showColumns,
      size,
      pagination,
      densityList,
      columns,
      statusOptions,
      search,
      onPageChange,
      reset,
      handleSelectDensity,
      handleChange,
      popupVisibleChange,
      addModalVisible,
      addFormRef,
      addForm,
      addFormRules,
      openAddModal,
      closeAddModal,
      handleAdd,
      deleteModalVisible,
      currentRecord,
      handleDelete,
      confirmDelete,
      cancelDelete,
      handleView,
      sftpModalVisible,
      sftpLoading,
      sftpFileList,
      sftpSearchPath,
      uploadModalVisible,
      breadcrumbItems,
      renameModalVisible,
      renameForm,
      handleSftp,
      handleBreadcrumbClick,
      handleSftpRefresh,
      handleSftpPathSearch,
      handleSftpItemClick,
      handleSftpUpload,
      handleSftpDownload,
      handleSftpCompress,
      handleSftpRename,
      handleRenameConfirm,
      closeRenameModal,
      handleSftpFileDelete,
      handleUploadClick,
      handleUploadCancel,
      handleUploadConfirm,
      formatFileSize,
      webShellVisible,
      terminalRef,
      handleWebShell,
      handleWebShellClose,
    };
  },
};
</script>

<style scoped lang="less">
.container {
  padding: 0 20px 20px 20px;
}
:deep(.arco-table-th) {
  &:last-child {
    .arco-table-th-item-title {
      margin-left: 16px;
    }
  }
}
.action-icon {
  margin-left: 12px;
  cursor: pointer;
}
.active {
  color: #0960bd;
  background-color: #e3f4fc;
}
.setting {
  display: flex;
  align-items: center;
  width: 200px;
  .title {
    margin-left: 12px;
    cursor: pointer;
  }
}
.sftp-container {
  display: flex;
  flex-direction: column;
  height: 600px;
}

.sftp-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.sftp-content {
  flex: 1;
  overflow: hidden;
  position: relative;

  .sftp-table {
    height: 100%;

    :deep(.arco-table-container) {
      height: 100%;
    }

    :deep(.arco-table-body) {
      overflow-y: auto;
    }

    :deep(.arco-table-th) {
      background-color: var(--color-fill-2);
    }

    :deep(.arco-table-td) {
      .arco-btn {
        padding: 0 4px;
        height: 24px;
        line-height: 24px;
        font-size: 14px;
      }
    }
  }
}

.sftp-breadcrumb {
  margin: 8px 0;
  padding: 8px;
  background-color: var(--color-fill-2);
  border-radius: 4px;

  :deep(.arco-breadcrumb) {
    .arco-breadcrumb-item {
      .arco-btn {
        padding: 0 4px;
        height: 24px;
        line-height: 24px;
        font-size: 14px;

        &:hover {
          background-color: var(--color-fill-3);
        }
      }
    }
  }
}

.web-shell-container {
  width: 100%;
  height: calc(90vh - 120px);
  background-color: #1e1e1e;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.terminal {
  width: 100%;
  height: 100%;
  padding: 8px;
  position: absolute;
  top: 0;
  left: 0;
}

:deep(.arco-modal) {
  padding: 0;
}

:deep(.arco-modal-header) {
  margin-bottom: 0;
  padding: 12px 16px;
}

:deep(.arco-modal-body) {
  padding: 0;
}
</style>

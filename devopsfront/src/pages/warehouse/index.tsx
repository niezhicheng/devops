import React, { useEffect, useState } from 'react';
import { Button, Card, Form, Input, Modal, Popconfirm, Space, Table, message } from 'antd';
import { PlusOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useRequest } from 'ahooks';
import { createRepository, deleteRepository, getRepositories, updateRepository } from '@/services/warehouse';

const Warehouse: React.FC = () => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);

  // 获取仓库列表
  const { data: repositories = [], refresh } = useRequest(getRepositories);

  // 处理表单提交
  const handleSubmit = async (values: any) => {
    try {
      if (editingId) {
        await updateRepository(editingId, values);
        message.success(t('仓库更新成功'));
      } else {
        await createRepository(values);
        message.success(t('仓库创建成功'));
      }
      setModalVisible(false);
      form.resetFields();
      refresh();
    } catch (error) {
      message.error(t('操作失败'));
    }
  };

  // 处理删除
  const handleDelete = async (id: number) => {
    try {
      await deleteRepository(id);
      message.success(t('仓库删除成功'));
      refresh();
    } catch (error) {
      message.error(t('删除失败'));
    }
  };

  // 表格列定义
  const columns = [
    {
      title: t('仓库名称'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('仓库地址'),
      dataIndex: 'url',
      key: 'url',
    },
    {
      title: t('默认分支'),
      dataIndex: 'defaultBranch',
      key: 'defaultBranch',
    },
    {
      title: t('最后同步时间'),
      dataIndex: 'lastSyncAt',
      key: 'lastSyncAt',
      render: (text: string) => text ? new Date(text).toLocaleString() : '-',
    },
    {
      title: t('状态'),
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: t('操作'),
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => {
              setEditingId(record.id);
              form.setFieldsValue(record);
              setModalVisible(true);
            }}
          >
            {t('编辑')}
          </Button>
          <Popconfirm
            title={t('确定要删除这个仓库吗？')}
            onConfirm={() => handleDelete(record.id)}
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              {t('删除')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="p-6">
      <Card
        title={t('仓库管理')}
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              setEditingId(null);
              form.resetFields();
              setModalVisible(true);
            }}
          >
            {t('添加仓库')}
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={repositories}
          rowKey="id"
          pagination={false}
        />
      </Card>

      <Modal
        title={editingId ? t('编辑仓库') : t('添加仓库')}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
        }}
        onOk={() => form.submit()}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="name"
            label={t('仓库名称')}
            rules={[{ required: true, message: t('请输入仓库名称') }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="url"
            label={t('仓库地址')}
            rules={[{ required: true, message: t('请输入仓库地址') }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="token"
            label={t('访问令牌')}
            rules={[{ required: true, message: t('请输入访问令牌') }]}
          >
            <Input.Password />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Warehouse; 
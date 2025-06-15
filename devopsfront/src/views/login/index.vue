<template>
  <div class="login">
    <div class="login-form">
      <div class="login-form-title">Login Arco Admin</div>

      <a-form
        ref="loginForm"
        :model="userInfo"
        layout="vertical"
        @submit="handleSubmit"
      >
        <a-form-item
          field="username"
          :rules="[{ required: true, message: '用户名不能为空' }]"
          :validate-trigger="['change', 'blur']"
          hide-label
        >
          <a-input v-model="userInfo.username" placeholder="用户名：admin" @keyup.enter="handleSubmit">
            <template #prefix>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          field="password"
          :rules="[{ required: true, message: '密码不能为空' }]"
          :validate-trigger="['change', 'blur']"
          hide-label
        >
          <a-input v-model="userInfo.password" placeholder="密码：admin" type="password" @keyup.enter="handleSubmit">
            <template #prefix>
              <icon-lock />
            </template>
          </a-input>
        </a-form-item>

        <a-space :size="16" direction="vertical">
          <div class="login-form-password-actions">
            <a-checkbox v-model="rememberPassword" @change="handleRememberPasswordChange">
              记住密码
            </a-checkbox>
            <a-link>忘记密码？</a-link>
          </div>
          <a-button type="primary" html-type="submit" long :loading="loading">
            登录
          </a-button>
          <a-button type="text" long class="login-form-register-btn">
            注册账号
          </a-button>
        </a-space>
      </a-form>
    </div>
  </div>
</template>

<script setup>
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { reactive, ref } from 'vue'
import useLoading from '@/hooks/loading'

const userInfo = reactive({
  username: 'admin',
  password: 'admin',
})
const { loading, setLoading } = useLoading()
const rememberPassword = ref(false)

const store = useStore()
const router = useRouter()

const handleRememberPasswordChange = (value) => {
  rememberPassword.value = value
  if (value) {
    localStorage.setItem('rememberedUser', JSON.stringify(userInfo))
  } else {
    localStorage.removeItem('rememberedUser')
  }
}

// 检查是否有记住的密码
const checkRememberedUser = () => {
  const rememberedUser = localStorage.getItem('rememberedUser')
  if (rememberedUser) {
    const { username, password } = JSON.parse(rememberedUser)
    userInfo.username = username
    userInfo.password = password
    rememberPassword.value = true
  }
}

const handleSubmit = async () => {
  setLoading(true)
  try {
    const result = await store.dispatch('user/login', userInfo)
    if (result) {
      if (rememberPassword.value) {
        localStorage.setItem('rememberedUser', JSON.stringify(userInfo))
      }
      await router.push('/')
    }
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    setLoading(false)
  }
}

// 初始化时检查记住的密码
checkRememberedUser()
</script>

<style scoped>
.login {
  position: fixed;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-fill-2);
}

.login-form {
  width: 352px;
  padding: 32px;
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.login-form-title {
  color: var(--color-text-1);
  font-weight: 500;
  font-size: 24px;
  line-height: 32px;
  margin-bottom: 24px;
  text-align: center;
}

.login-form-password-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.login-form-register-btn {
  color: var(--color-text-3) !important;
}
</style>

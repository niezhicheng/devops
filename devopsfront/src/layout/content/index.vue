<template>
  <div class="page-wrapper">
    <navbar />

    <div class="content-wrapper">
      <side-bar />

      <div class="main-content">
        <transition name="fade-transform" mode="out-in">
          <router-view class="content" />
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup>
import Navbar from '@/layout/navbar'
import SideBar from '@/layout/sidebar'
</script>

<style scoped>
.page-wrapper {
  box-sizing: border-box;
  padding-top: 50px;
  width: 100%;
  height: 100vh;
  background-color: var(--color-bg-1);
}

.content-wrapper {
  display: flex;
  width: 100%;
  height: calc(100vh - 50px);
}

.main-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.content {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  padding: 20px;
  box-sizing: border-box;
}

/* 路由切换动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>

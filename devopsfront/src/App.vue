<template>
  <router-view />
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
  font-size: 14px;
  line-height: 1.5;
  color: var(--color-text-1);
  background-color: var(--color-bg-1);
}

#app {
  width: 100%;
  height: 100%;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--color-bg-2);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: var(--color-fill-3);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--color-fill-4);
}
</style>

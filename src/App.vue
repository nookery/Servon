<script setup lang="ts">
import MainLayout from './layouts/MainLayout.vue'
import Toast from './components/Toast.vue'
import ErrorAlert from './components/ErrorAlert.vue'
import GlobalConfirm from './components/GlobalConfirm.vue'
import LockScreen from './components/LockScreen.vue'
import { useLockStore } from './stores/lockStore'
import { ref, provide } from 'vue'

const lockStore = useLockStore()
const errorAlertRef = ref<InstanceType<typeof ErrorAlert> | null>(null)

// 提供错误显示方法给其他组件使用
provide('showError', (message: string) => {
  errorAlertRef.value?.showError(message)
})
</script>

<template>
  <template v-if="!lockStore.isLocked">
    <main-layout>
      <router-view></router-view>
    </main-layout>
  </template>
  <LockScreen v-else />
  <Toast />
  <ErrorAlert ref="errorAlertRef" />
  <GlobalConfirm />
</template>

<style>
@import 'remixicon/fonts/remixicon.css';

body {
  margin: 0;
  padding: 0;
}
</style>

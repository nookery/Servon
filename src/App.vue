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
const headerErrors = ref<{ id: number; message: string }[]>([])
let nextHeaderErrorId = 1

// 提供错误显示方法给其他组件使用，默认在头部显示
provide('showError', (message: string, showInHeader = true) => {
  if (showInHeader) {
    const id = nextHeaderErrorId++
    headerErrors.value.push({ id, message })
  } else {
    errorAlertRef.value?.showError(message)
  }
})

const removeHeaderError = (id: number) => {
  headerErrors.value = headerErrors.value.filter(error => error.id !== id)
}
</script>

<template>
  <template v-if="!lockStore.isLocked">
    <main-layout>
      <!-- 头部错误提示 -->
      <div v-for="error in headerErrors" :key="error.id" class="alert alert-error shadow-lg mb-4 animate-slide-in-down">
        <div class="flex-1 flex items-center gap-2">
          <i class="ri-error-warning-line text-xl" />
          <span>{{ error.message }}</span>
        </div>
        <IconButton icon="ri-close-line" variant="ghost" circle size="sm" @click="removeHeaderError(error.id)" />
      </div>
      <router-view></router-view>
    </main-layout>
  </template>
  <LockScreen v-else />
  <Toast />
  <ErrorAlert ref="errorAlertRef" />
  <GlobalConfirm />
</template>

<style>
.animate-slide-in-down {
  animation: slide-in-down 0.3s ease-out;
}

@keyframes slide-in-down {
  from {
    transform: translateY(-100%);
    opacity: 0;
  }

  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
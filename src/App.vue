<script setup lang="ts">
import MainLayout from './layouts/MainLayout.vue'
import Toast from './components/Toast.vue'
import GlobalConfirm from './components/GlobalConfirm.vue'
import LockScreen from './components/LockScreen.vue'
import { useLockStore } from './stores/lockStore'

const lockStore = useLockStore()
</script>

<template>
  <template v-if="!lockStore.isLocked">
    <main-layout>
      <router-view></router-view>
    </main-layout>
  </template>
  <LockScreen v-else />
  <Toast />
  <GlobalConfirm />
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>

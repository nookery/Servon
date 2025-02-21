<script setup lang="ts">
import MainLayout from './layouts/MainLayout.vue'
import Toast from './components/Toast.vue'
import GlobalConfirm from './components/GlobalConfirm.vue'
</script>

<template>
  <main-layout>
    <router-view></router-view>
  </main-layout>
  <Toast />
  <GlobalConfirm />
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>

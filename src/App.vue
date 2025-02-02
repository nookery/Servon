<script setup lang="ts">
import MainLayout from './layouts/MainLayout.vue'
import { NConfigProvider, NMessageProvider } from 'naive-ui'
</script>

<template>
  <n-config-provider>
    <n-message-provider>
      <main-layout>
        <router-view></router-view>
      </main-layout>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>

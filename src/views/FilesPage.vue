<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import FileManager from '../components/FileManager.vue'

const currentPath = ref('/')
</script>

<template>
    <PageContainer title="文件管理">
        <FileManager v-model:path="currentPath" :initial-path="currentPath" :show-breadcrumbs="true"
            :show-toolbar="true" :show-pagination="true" :show-shortcuts="true" :read-only="false" />
    </PageContainer>
</template>

<style>
@import 'remixicon/fonts/remixicon.css';
</style>
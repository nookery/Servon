<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import FileManager from '../components/files/FileManager.vue'
import type { SortBy, SortOrder } from '../models/FileInfo'

const currentPath = ref('/')
const currentSort = ref<SortBy>('name')
const sortOrder = ref<SortOrder>('asc')
</script>

<template>
    <PageContainer title="文件管理">
        <FileManager v-model:path="currentPath" :initial-path="currentPath" :show-breadcrumbs="true"
            :show-toolbar="true" :show-pagination="true" :show-shortcuts="true" :read-only="false"
            v-model:sort-by="currentSort" v-model:sort-order="sortOrder" />
    </PageContainer>
</template>

<style>
@import 'remixicon/fonts/remixicon.css';
</style>
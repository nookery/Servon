<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import FileManager from '../components/files/FileManager.vue'

const currentPath = ref('/data')
</script>

<template>
    <PageContainer title="数据中心">
        <FileManager v-model:path="currentPath" :initial-path="currentPath" :show-breadcrumbs="true"
            :show-toolbar="true" :show-pagination="true" :show-shortcuts="true" :read-only="false" />
    </PageContainer>
</template>

<style>
@import 'remixicon/fonts/remixicon.css';
</style>
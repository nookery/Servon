<script setup lang="ts">
import PageContainer from '../layouts/PageContainer.vue'
import FileManager from '../components/files/FileManager.vue'

const deployLogsPath = '/data/deploy'
</script>

<template>
    <PageContainer title="部署日志">
        <FileManager :initial-path="deployLogsPath" :show-toolbar="true" :show-breadcrumbs="true"
            :show-pagination="true" :read-only="false" :show-shortcuts="false" sort-by="modTime" sort-order="desc" />
    </PageContainer>
</template>
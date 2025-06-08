<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import FileManager from '../components/files/FileManager.vue'
import LogView from '../components/logs/LogView.vue'
import { RiListCheck, RiFolderOpenLine } from '@remixicon/vue'

const currentDir = ref('')
const logDir = '/data/logs'
const activeTab = ref('simple')
const error = ref<string | null>(null)

const tabs = [
    { key: 'simple', title: '简单视图', icon: RiListCheck },
    { key: 'explorer', title: '文件浏览', icon: RiFolderOpenLine }
]
</script>

<template>
    <PageContainer title="日志管理" :error="error" :tabs="tabs" v-model="activeTab" :full-height="true">
        <template #simple>
            <LogView :current-dir="currentDir" />
        </template>

        <template #explorer>
            <FileManager :initial-path="logDir" :show-toolbar="true" :show-breadcrumbs="true" :show-pagination="true"
                :read-only="false" :show-shortcuts="false" />
        </template>
    </PageContainer>
</template>
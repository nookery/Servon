<script setup lang="ts">
import { ref } from 'vue'
import FileManager from '../FileManager.vue'

const props = defineProps<{
    initialPath?: string
}>()

const currentPath = ref(props.initialPath || '/data/github/config')
const fileManagerRef = ref<InstanceType<typeof FileManager> | null>(null)

</script>

<template>
    <FileManager ref="fileManagerRef" v-model:path="currentPath" :initial-path="currentPath" :show-breadcrumbs="true"
        :show-toolbar="false" :show-pagination="true" :show-shortcuts="false" :read-only="true" />
</template>
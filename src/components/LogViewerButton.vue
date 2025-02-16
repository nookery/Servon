<script setup lang="ts">
import { useLogViewerStore } from '../stores/logViewer'

const logViewerStore = useLogViewerStore()
</script>

<template>
    <div>
        <!-- 日志查看器按钮 -->
        <button @click="logViewerStore.toggleVisibility" class="btn btn-ghost btn-circle" title="日志查看器">
            <i class="ri-file-list-line text-xl"></i>
        </button>
    </div>
</template>
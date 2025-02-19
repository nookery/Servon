<script setup lang="ts">
defineProps<{
    title?: string
}>()
</script>

<template>
    <div class="card bg-base-100 shadow-xl h-full flex flex-col">
        <div class="card-body p-0 flex flex-col h-full">
            <!-- 固定头部 -->
            <div class="sticky top-0 bg-base-100 z-10 p-6 pb-0">
                <!-- <h2 class="card-title mb-4">{{ title }}</h2> -->
                <!-- 头部内容插槽 -->
                <slot name="header"></slot>
                <!-- Tab栏插槽 -->
                <slot name="tabs"></slot>
            </div>

            <!-- 可滚动的内容区域 -->
            <div class="flex-1 overflow-auto p-6 pt-0">
                <slot></slot>
            </div>
        </div>
    </div>
</template>
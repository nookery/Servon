<script setup lang="ts">
import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent } from 'naive-ui'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
import { ref } from 'vue'

const collapsed = ref(false)
</script>

<template>
    <n-layout position="absolute">
        <n-layout-header bordered style="height: 64px; padding: 0">
            <Header />
        </n-layout-header>
        <n-layout position="absolute" style="top: 64px; bottom: 0" has-sider>
            <n-layout-sider bordered show-trigger collapse-mode="width" :collapsed-width="64" :width="240"
                :native-scrollbar="false" :collapsed="collapsed" @collapse="collapsed = true"
                @expand="collapsed = false">
                <Sidebar :collapsed="collapsed" />
            </n-layout-sider>
            <n-layout>
                <n-layout-content style="padding: 16px">
                    <slot></slot>
                </n-layout-content>
            </n-layout>
        </n-layout>
    </n-layout>
</template>

<style>
body {
    margin: 0;
    padding: 0;
}

.n-layout-content {
    background: #f5f5f5;
    min-height: 100%;
}
</style>
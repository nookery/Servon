<script setup lang="ts">
import Header from './Header.vue'
import DockBar from './DockBar.vue'
</script>

<template>
    <div class="flex flex-col min-h-screen bg-gradient-to-br from-cyan-300 to-cyan-500 overflow-hidden">
        <!-- Header (顶部菜单栏) -->
        <div class="w-full h-16">
            <Header />
        </div>

        <!-- Main Layout (桌面区域) -->
        <div class="flex-1 overflow-hidden relative p-6 pt-2">
            <!-- Content (窗口内容) -->
            <div class="bg-base-100 rounded-xl shadow-lg p-6 h-[calc(100vh-9rem)] overflow-auto">
                <slot></slot>
            </div>

            <!-- Dock at bottom -->
            <div class="dock-wrapper fixed bottom-4 left-0 right-0 flex justify-center items-center z-10">
                <DockBar />
            </div>
        </div>
    </div>
</template>
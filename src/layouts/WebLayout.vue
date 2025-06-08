<script setup lang="ts">
import Header from './Header.vue'
import Sidebar from './SideBar.vue'
import { useLayoutStore } from '../stores/layout'

const layoutStore = useLayoutStore()
</script>

<template>
    <div class="flex flex-col min-h-screen bg-base-200">
        <!-- Header -->
        <div class="w-full h-16">
            <Header />
        </div>

        <!-- Main Layout -->
        <div class="flex flex-1 overflow-hidden">
            <!-- Sidebar -->
            <div :class="[
                'transition-all duration-300 border-r border-base-300 h-[calc(100vh-4rem)] bg-base-100 fixed left-0 top-16',
                layoutStore.collapsed ? 'w-16' : 'w-40'
            ]">
                <div class="sticky top-16 h-full">
                    <Sidebar :collapsed="layoutStore.collapsed" />
                </div>
            </div>

            <div class="flex flex-row w-full">
                <!-- Content -->
                <div :class="[
                    'flex-1 p-0 transition-all duration-300 overflow-auto h-[calc(100vh-4rem)]',
                    layoutStore.collapsed ? 'ml-16' : 'ml-40',
                ]">
                    <slot></slot>
                </div>
            </div>
        </div>
    </div>
</template>
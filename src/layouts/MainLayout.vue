<script setup lang="ts">
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
import { ref } from 'vue'

const collapsed = ref(false)
</script>

<template>
    <div class="flex flex-col min-h-screen bg-base-200">
        <!-- Header -->
        <div class="w-full h-16">
            <Header />
        </div>

        <!-- Main Layout -->
        <div class="flex flex-1">
            <!-- Sidebar -->
            <div :class="[
                'transition-all duration-300 border-r border-base-300 h-full bg-base-100 fixed left-0 top-16',
                collapsed ? 'w-16' : 'w-60'
            ]">
                <div class="sticky top-16 border-0 border-red-500 h-full">
                    <button @click="collapsed = !collapsed"
                        class="btn btn-ghost btn-sm absolute -right-3 top-3 z-50 rounded-full bg-base-100 border border-base-300">
                        <i :class="[
                            collapsed ? 'ri-arrow-right-s-line' : 'ri-arrow-left-s-line',
                            'text-lg'
                        ]"></i>
                    </button>
                    <Sidebar :collapsed="collapsed" />
                </div>
            </div>

            <!-- Content -->
            <div :class="[
                'flex-1 p-4 transition-all duration-300',
                collapsed ? 'ml-16' : 'ml-60'
            ]">
                <slot></slot>
            </div>
        </div>
    </div>
</template>

<style>
body {
    margin: 0;
    padding: 0;
}
</style>
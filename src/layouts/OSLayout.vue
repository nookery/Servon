<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import Header from './Header.vue'
import DockBar from './DockBar.vue'
import WindowManager from '../components/WindowManager.vue'

// 背景动画控制
const bgPosition = ref({ x: 0, y: 0 })
const mousePosition = ref({ x: 0, y: 0 })

// 跟踪鼠标位置，用于视差效果
const handleMouseMove = (e: MouseEvent) => {
    mousePosition.value = {
        x: e.clientX / window.innerWidth,
        y: e.clientY / window.innerHeight
    }
}

// 动画帧
let animationFrame: number

const animateBackground = () => {
    // 缓慢跟随鼠标移动，创建平滑视差效果
    bgPosition.value.x += (mousePosition.value.x * 10 - bgPosition.value.x) * 0.02
    bgPosition.value.y += (mousePosition.value.y * 10 - bgPosition.value.y) * 0.02

    animationFrame = requestAnimationFrame(animateBackground)
}

onMounted(() => {
    window.addEventListener('mousemove', handleMouseMove)
    animationFrame = requestAnimationFrame(animateBackground)
})

onUnmounted(() => {
    window.removeEventListener('mousemove', handleMouseMove)
    cancelAnimationFrame(animationFrame)
})
</script>

<template>
    <div class="desktop-container min-h-screen overflow-hidden">
        <!-- 使用 Tailwind 实现的简洁渐变背景 -->
        <div class="absolute inset-0 bg-gradient-to-br from-slate-900 to-indigo-950 z-0"></div>

        <!-- Header (顶部菜单栏) -->
        <div class="w-full h-16 z-30 relative">
            <Header />
        </div>

        <!-- Main Layout (桌面区域) -->
        <div class="flex-1 overflow-hidden relative pt-16 z-20">
            <!-- Desktop (桌面) -->
            <div class="w-full h-full">
                <!-- 窗口管理器 -->
                <WindowManager />
            </div>

            <!-- Dock at bottom -->
            <div class="dock-wrapper fixed bottom-4 w-auto left-1/2 -translate-x-1/2 z-10">
                <DockBar />
            </div>
        </div>
    </div>
</template>

<style scoped>
.desktop-container {
    position: relative;
    display: flex;
    flex-direction: column;
}
</style>
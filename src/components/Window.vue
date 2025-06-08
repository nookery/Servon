<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { useWindowStore, type WindowState } from '../stores/windowStore'
import { RiCloseLine, RiSubtractLine, RiCheckboxMultipleLine } from '@remixicon/vue'

// 导入所有可能的视图组件
import DashboardView from '../views/Dashboard.vue'
import SoftwareView from '../views/Software.vue'
import ProcessesView from '../views/Processes.vue'
import FilesView from '../views/FilesPage.vue'
import PortsView from '../views/Ports.vue'
import UsersView from '../views/Users.vue'
import IntegrationsView from '../views/Integrations.vue'
import DataView from '../views/DataPage.vue'
import LogsView from '../views/Logs.vue'
import ProjectsView from '../views/ProjectsPage.vue'
import EmptyView from '../views/EmptyView.vue'

// 创建组件映射
const components = {
    DashboardView,
    SoftwareView,
    ProcessesView,
    FilesView,
    PortsView,
    UsersView,
    IntegrationsView,
    DataView,
    LogsView,
    ProjectsView,
    EmptyView
}

const props = defineProps<{
    window: WindowState
}>()

const windowStore = useWindowStore()

// 引用窗口元素以进行拖动操作
const windowRef = ref<HTMLElement | null>(null)
const windowHeaderRef = ref<HTMLElement | null>(null)
const resizeHandleRef = ref<HTMLElement | null>(null)

// 是否正在拖动或调整大小
const isDragging = ref(false)
const isResizing = ref(false)
const dragOffset = ref({ x: 0, y: 0 })

// 本地状态用于平滑调整大小
const localSize = reactive({
    width: props.window.size.width,
    height: props.window.size.height
})

// 计算当前组件
const currentComponent = computed(() => {
    // 如果组件名为空或组件不存在，返回空白组件
    if (!props.window.component || !components[props.window.component as keyof typeof components]) {
        return EmptyView;
    }
    return components[props.window.component as keyof typeof components];
})

// 计算样式
const windowStyle = computed(() => {
    if (props.window.isMaximized) {
        return {
            width: '100%',
            height: 'calc(100vh - 8rem)',
            left: '0',
            top: '0',
            transform: 'none',
            zIndex: props.window.zIndex
        }
    } else {
        return {
            width: `${isResizing.value ? localSize.width : props.window.size.width}px`,
            height: `${isResizing.value ? localSize.height : props.window.size.height}px`,
            left: `${props.window.position.x}px`,
            top: `${props.window.position.y}px`,
            zIndex: props.window.zIndex
        }
    }
})

// 处理窗口激活
function handleWindowClick() {
    if (!props.window.isActive) {
        windowStore.activateWindow(props.window.id)
    }
}

// 开始拖动
function startDrag(e: MouseEvent) {
    if (props.window.isMaximized) return

    isDragging.value = true
    const rect = windowRef.value?.getBoundingClientRect()
    if (rect) {
        dragOffset.value = {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        }
    }
}

// 开始调整大小
function startResize(e: MouseEvent) {
    if (props.window.isMaximized) return

    isResizing.value = true
    dragOffset.value = {
        x: e.clientX,
        y: e.clientY
    }

    // 初始化本地尺寸
    localSize.width = props.window.size.width
    localSize.height = props.window.size.height
}

// 处理关闭窗口
function closeWindow() {
    windowStore.closeWindow(props.window.id)
}

// 处理最大化/还原窗口
function toggleMaximize() {
    windowStore.toggleMaximize(props.window.id)
}

// 使用requestAnimationFrame优化性能
let rafId: number | null = null
let pendingUpdate = false

onMounted(() => {
    // 设置全局鼠标移动和释放事件处理程序
    const handleMouseMove = (e: MouseEvent) => {
        if (isDragging.value) {
            const newX = e.clientX - dragOffset.value.x
            const newY = e.clientY - dragOffset.value.y
            windowStore.updateWindowPosition(props.window.id, newX, newY)
        } else if (isResizing.value && windowRef.value) {
            // 实时更新本地尺寸以获得即时反馈
            localSize.width = Math.max(300, localSize.width + (e.clientX - dragOffset.value.x))
            localSize.height = Math.max(200, localSize.height + (e.clientY - dragOffset.value.y))

            // 更新拖动偏移量
            dragOffset.value = {
                x: e.clientX,
                y: e.clientY
            }

            // 使用rAF进行节流更新store
            if (!pendingUpdate) {
                pendingUpdate = true
                rafId = requestAnimationFrame(() => {
                    windowStore.updateWindowSize(props.window.id, localSize.width, localSize.height)
                    pendingUpdate = false
                })
            }
        }
    }

    const handleMouseUp = () => {
        if (isResizing.value) {
            // 调整结束时，确保store中的尺寸与本地尺寸一致
            windowStore.updateWindowSize(props.window.id, localSize.width, localSize.height)
            if (rafId !== null) {
                cancelAnimationFrame(rafId)
                rafId = null
            }
        }
        isDragging.value = false
        isResizing.value = false
    }

    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseup', handleMouseUp)

    return () => {
        window.removeEventListener('mousemove', handleMouseMove)
        window.removeEventListener('mouseup', handleMouseUp)
        if (rafId !== null) {
            cancelAnimationFrame(rafId)
        }
    }
})
</script>

<template>
    <div ref="windowRef" class="window bg-base-100 rounded-lg shadow-xl overflow-hidden flex flex-col"
        :class="{ 'active': window.isActive }" :style="windowStyle" @mousedown="handleWindowClick">
        <!-- Window Header -->
        <div ref="windowHeaderRef" class="window-header p-2 bg-base-200 flex items-center justify-between cursor-move"
            @mousedown="startDrag">
            <div class="flex items-center gap-2">
                <i :class="[window.icon, 'text-lg']"></i>
                <span class="font-medium">{{ window.title }}</span>
            </div>

            <div class="flex gap-1">
                <button class="btn btn-circle btn-ghost btn-xs" @click.stop="windowStore.toggleMaximize(window.id)">
                    <RiSubtractLine />
                </button>
                <button class="btn btn-circle btn-ghost btn-xs" @click.stop="toggleMaximize">
                    <RiCheckboxMultipleLine />
                </button>
                <button class="btn btn-circle btn-ghost btn-xs hover:bg-error hover:text-white"
                    @click.stop="closeWindow">
                    <RiCloseLine />
                </button>
            </div>
        </div>

        <!-- Window Content -->
        <div class="window-content flex-1 overflow-auto p-4">
            <component :is="currentComponent" v-bind="window.props" />
        </div>

        <!-- Resize Handle -->
        <div v-if="!window.isMaximized" ref="resizeHandleRef" class="resize-handle" @mousedown.stop="startResize"></div>
    </div>
</template>

<style scoped>
.window {
    position: absolute;
    /* 移除过渡效果，使大小调整更加即时 */
    /* transition: width 0.2s, height 0.2s; */
}

.window.active {
    box-shadow: 0 5px 30px rgba(0, 0, 0, 0.2);
}

.resize-handle {
    position: absolute;
    right: 0;
    bottom: 0;
    width: 15px;
    height: 15px;
    cursor: nwse-resize;
    background: linear-gradient(135deg, transparent 50%, rgba(128, 128, 128, 0.5) 50%);
}
</style>
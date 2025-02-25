import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { v4 as uuidv4 } from 'uuid'

export interface WindowState {
    id: string;
    title: string;
    icon: string;
    component: string;
    props?: any;
    position: {
        x: number;
        y: number;
    };
    size: {
        width: number;
        height: number;
    };
    isActive: boolean;
    isMaximized: boolean;
    zIndex: number;
}

export const useWindowStore = defineStore('window', () => {
    const windows = ref<WindowState[]>([])
    const nextZIndex = ref(100)

    const activeWindow = computed(() =>
        windows.value.find(w => w.isActive)
    )

    // 打开窗口
    function openWindow(windowData: Partial<WindowState>) {
        // 检查是否已经有同类型窗口打开
        const existingWindow = windows.value.find(w => w.component === windowData.component)

        if (existingWindow) {
            // 如果已经有窗口，激活它
            activateWindow(existingWindow.id)
            return existingWindow.id
        }

        // 随机计算新窗口的位置，避免完全重叠
        const randomOffset = windows.value.length * 20

        // 创建新窗口
        const newWindow: WindowState = {
            id: uuidv4(),
            title: windowData.title || '新窗口',
            icon: windowData.icon || 'ri-window-line',
            component: windowData.component || '',
            props: windowData.props || {},
            position: windowData.position || {
                x: 100 + randomOffset,
                y: 100 + randomOffset
            },
            size: windowData.size || {
                width: 800,
                height: 600
            },
            isActive: true,
            isMaximized: false,
            zIndex: nextZIndex.value
        }

        // 将所有其他窗口设为非活动
        windows.value.forEach(w => {
            w.isActive = false
        })

        // 添加新窗口
        windows.value.push(newWindow)
        nextZIndex.value += 1

        return newWindow.id
    }

    // 关闭窗口
    function closeWindow(id: string) {
        const index = windows.value.findIndex(w => w.id === id)
        if (index !== -1) {
            windows.value.splice(index, 1)

            // 如果还有其他窗口，激活最顶层的
            if (windows.value.length > 0) {
                const topWindow = windows.value.reduce((prev, current) =>
                    (prev.zIndex > current.zIndex) ? prev : current
                )
                activateWindow(topWindow.id)
            }
        }
    }

    // 激活窗口（置于顶层）
    function activateWindow(id: string) {
        windows.value.forEach(w => {
            if (w.id === id) {
                w.isActive = true
                w.zIndex = nextZIndex.value
                nextZIndex.value += 1
            } else {
                w.isActive = false
            }
        })
    }

    // 最大化/还原窗口
    function toggleMaximize(id: string) {
        const window = windows.value.find(w => w.id === id)
        if (window) {
            window.isMaximized = !window.isMaximized
        }
    }

    // 更新窗口位置
    function updateWindowPosition(id: string, x: number, y: number) {
        const window = windows.value.find(w => w.id === id)
        if (window && !window.isMaximized) {
            window.position.x = x
            window.position.y = y
        }
    }

    // 更新窗口大小
    function updateWindowSize(id: string, width: number, height: number) {
        const window = windows.value.find(w => w.id === id)
        if (window && !window.isMaximized) {
            window.size.width = width
            window.size.height = height
        }
    }

    return {
        windows,
        activeWindow,
        openWindow,
        closeWindow,
        activateWindow,
        toggleMaximize,
        updateWindowPosition,
        updateWindowSize
    }
}) 
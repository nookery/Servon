<script setup lang="ts">
import { watch } from 'vue'
import { confirmState } from '../composables/useConfirmState'

// 处理确认
function handleConfirm() {
    confirmState.resolve?.(true)
    confirmState.show = false
}

// 处理取消
function handleCancel() {
    confirmState.resolve?.(false)
    confirmState.show = false
}

// 监听显示状态
watch(() => confirmState.show, (show) => {
    if (!show) {
        confirmState.resolve?.(false)
    }
})

// 根据类型获取图标路径
function getIconPath() {
    switch (confirmState.options.type) {
        case 'error':
            return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
        case 'info':
            return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
        default: // warning
            return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
    }
}

// 根据类型获取图标颜色
function getIconColor() {
    switch (confirmState.options.type) {
        case 'error':
            return 'stroke-error'
        case 'info':
            return 'stroke-info'
        default: // warning
            return 'stroke-warning'
    }
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': confirmState.show }">
        <div class="modal-box p-0">
            <div role="alert" class="alert">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="h-6 w-6 shrink-0"
                    :class="getIconColor()">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPath()">
                    </path>
                </svg>
                <div class="flex-1">
                    <h3 class="font-bold">{{ confirmState.options.title }}</h3>
                    <div class="text-sm">{{ confirmState.options.message }}</div>
                </div>
                <div class="flex gap-2">
                    <button class="btn btn-sm" @click="handleCancel">
                        {{ confirmState.options.cancelText || '取消' }}
                    </button>
                    <button class="btn btn-sm" :class="{
                        'btn-error': confirmState.options.type === 'error',
                        'btn-warning': confirmState.options.type === 'warning',
                        'btn-info': confirmState.options.type === 'info'
                    }" @click="handleConfirm">
                        {{ confirmState.options.confirmText || '确认' }}
                    </button>
                </div>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button @click="handleCancel">关闭</button>
        </form>
    </dialog>
</template>
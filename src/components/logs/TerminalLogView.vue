<script setup lang="ts">
import { ref, onMounted, watch, computed, onBeforeUnmount } from 'vue'
import { RiEmotionHappyLine } from '@remixicon/vue'
import type { LogEntry } from '../../types/log'
import * as monaco from 'monaco-editor'

const props = defineProps<{
    logEntries: LogEntry[],
    visibleFields: string[]
}>()

const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

// 格式化终端日志内容，按照时间倒序排列
const terminalContent = computed(() => {
    if (props.logEntries.length === 0) return ''

    // 创建日志条目的副本，以便排序
    const sortedEntries = [...props.logEntries].reverse()

    return sortedEntries.map(entry => {
        const parts = []

        if (props.visibleFields.includes('time')) {
            parts.push(entry.time)
        }

        if (props.visibleFields.includes('level')) {
            const level = entry.level.toUpperCase().padEnd(5, ' ')
            parts.push(`[${level}]`)
        }

        if (props.visibleFields.includes('caller') && entry.caller) {
            parts.push(`[${entry.caller}]`)
        }

        if (props.visibleFields.includes('message')) {
            parts.push(entry.message)
        }

        return parts.join(' ')
    }).join('\n')
})

// 初始化Monaco编辑器
function initEditor() {
    if (!editorContainer.value) return

    // 如果编辑器已存在，先销毁
    if (editor) {
        editor.dispose()
    }

    // 注册自定义日志语言
    monaco.languages.register({ id: 'custom-log' });

    // 定义日志语言的语法高亮规则
    monaco.languages.setMonarchTokensProvider('custom-log', {
        tokenizer: {
            root: [
                // 错误级别
                [/\[ERROR\s*\]/, 'log-error'],
                // 警告级别
                [/\[WARN\s*\]/, 'log-warning'],
                // 信息级别
                [/\[INFO\s*\]/, 'log-info'],
                // 调试级别
                [/\[DEBUG\s*\]/, 'log-debug'],
                // 时间戳
                [/\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}/, 'log-date'],
                // 调用位置
                [/\[[^\]]+\]/, 'log-caller'],
            ]
        }
    });

    // 创建编辑器实例
    editor = monaco.editor.create(editorContainer.value, {
        value: terminalContent.value,
        language: 'custom-log', // 使用自定义日志语言
        theme: 'custom-log-theme', // 使用自定义主题
        readOnly: true, // 只读模式
        minimap: { enabled: true }, // 启用小地图
        scrollBeyondLastLine: false, // 不允许滚动超过最后一行
        automaticLayout: true, // 自动调整布局
        wordWrap: 'on', // 启用自动换行
        fontSize: 12, // 字体大小
        fontFamily: 'Consolas, "Courier New", monospace', // 等宽字体
        lineNumbers: 'on', // 显示行号
        fixedOverflowWidgets: true, // 修复溢出部件
        scrollbar: {
            // 自定义滚动条
            verticalScrollbarSize: 10,
            horizontalScrollbarSize: 10,
            alwaysConsumeMouseWheel: false
        }
    })

    // 为不同日志级别添加语法高亮规则
    monaco.editor.defineTheme('custom-log-theme', {
        base: 'vs-dark',
        inherit: true,
        rules: [
            { token: 'log-error', foreground: 'ff5252', fontStyle: 'bold' },
            { token: 'log-warning', foreground: 'ffab40', fontStyle: 'bold' },
            { token: 'log-info', foreground: '40c4ff', fontStyle: 'bold' },
            { token: 'log-debug', foreground: 'aaaaaa', fontStyle: 'bold' },
            { token: 'log-date', foreground: 'bbbbbb' },
            { token: 'log-caller', foreground: '7986cb' },
        ],
        colors: {
            'editor.background': '#1e1e1e',
        }
    })

    monaco.editor.setTheme('custom-log-theme')

    // 添加窗口大小变化监听，确保编辑器正确调整大小
    window.addEventListener('resize', handleResize)
}

// 处理窗口大小变化
function handleResize() {
    if (editor) {
        editor.layout()
    }
}

// 监听日志条目变化，更新编辑器内容
watch(() => props.logEntries, () => {
    if (editor) {
        editor.setValue(terminalContent.value)
        // 滚动到底部
        if (terminalContent.value) {
            const lineCount = terminalContent.value.split('\n').length
            editor.revealLine(lineCount)
        }
    }
}, { deep: true })

// 监听可见字段变化，更新编辑器内容
watch(() => props.visibleFields, () => {
    if (editor) {
        editor.setValue(terminalContent.value)
    }
}, { deep: true })

// 组件挂载时初始化编辑器
onMounted(() => {
    initEditor()

    // 确保编辑器在挂载后正确布局
    setTimeout(() => {
        if (editor) {
            editor.layout()
            // 初始化时滚动到底部
            if (terminalContent.value) {
                const lineCount = terminalContent.value.split('\n').length
                editor.revealLine(lineCount)
            }
        }
    }, 100)
})

// 组件卸载时销毁编辑器
onBeforeUnmount(() => {
    // 移除窗口大小变化监听
    window.removeEventListener('resize', handleResize)

    if (editor) {
        editor.dispose()
        editor = null
    }
})
</script>

<template>
    <div class="h-full">
        <div v-if="logEntries.length === 0" class="flex flex-col items-center justify-center h-full gap-2">
            <RiEmotionHappyLine class="w-8 h-8 text-success" />
            <span>暂无符合条件的日志 </span>
            <span class="text-xs opacity-50">这说明系统运行得很顺利呢！</span>
        </div>
        <div v-else ref="editorContainer" class="h-full"></div>
    </div>
</template>

<style scoped>
/* 确保编辑器容器正确填充高度 */
:deep(.monaco-editor) {
    height: 100% !important;
}

:deep(.monaco-editor .overflow-guard) {
    height: 100% !important;
}
</style>
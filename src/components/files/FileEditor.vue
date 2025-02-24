<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { fileAPI } from '../../api/file_api'
import type { FileInfo } from '../../types/FileInfo'
import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import IconButton from '../IconButton.vue'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'
import {
    RiCodeLine,
    RiFileCopyLine,
    RiDownloadLine,
    RiRestartLine,
    RiSaveLine,
    RiCloseLine,
    RiAlertLine,
    RiMapLine,
    RiDeleteBinLine,
} from '@remixicon/vue'
import { getLanguageFromFileName, getSupportedLanguages } from '../../utils/languages'

const props = defineProps<{
    show: boolean
    file: FileInfo | null
    initialContent?: string  // 添加可选的初始内容属性
}>()

const emit = defineEmits<{
    'update:show': [value: boolean]
    'saved': []
}>()

const content = ref('')
const error = ref<string | null>(null)
const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

const toast = useToast()
const confirm = useConfirm()

const showMinimap = ref(true)

// 添加自动刷新相关的状态
const autoRefresh = ref(false)
const refreshInterval = ref(5) // 默认5秒
let refreshTimer: ReturnType<typeof setInterval> | null = null

const currentLanguage = ref('plaintext')
const supportedLanguages = getSupportedLanguages()

// 添加行数统计
const lineCount = ref(0)
const characterCount = ref(0)

// 配置 Monaco Editor workers
window.MonacoEnvironment = {
    getWorker(_, label) {
        switch (label) {
            case 'json':
                return new jsonWorker()
            case 'css':
            case 'scss':
            case 'less':
                return new cssWorker()
            case 'html':
                return new htmlWorker()
            case 'typescript':
            case 'javascript':
                return new tsWorker()
            default:
                return new editorWorker()
        }
    }
}

// 监听文件变化
watch(() => props.file, (newFile) => {
    if (newFile) {
        const detectedLanguage = getLanguageFromFileName(newFile.name)
        currentLanguage.value = detectedLanguage

        // 如果编辑器已存在，更新其语言
        if (editor) {
            const model = editor.getModel()
            if (model) {
                monaco.editor.setModelLanguage(model, detectedLanguage)
            }
        }
    }
}, { immediate: true })

onMounted(() => {
    if (editorContainer.value) {
        const initialLanguage = props.file
            ? getLanguageFromFileName(props.file.name)
            : 'plaintext'

        editor = monaco.editor.create(editorContainer.value, {
            value: content.value,
            language: initialLanguage,
            theme: 'vs-dark',
            automaticLayout: true,
            minimap: { enabled: showMinimap.value },
            scrollBeyondLastLine: false,
            formatOnPaste: true,
        })

        editor.addCommand(monaco.KeyMod.Shift | monaco.KeyMod.Alt | monaco.KeyCode.KeyF, () => {
            formatCode()
        })

        editor.onDidChangeModelContent(() => {
            content.value = editor?.getValue() || ''
            updateEditorStats()
        })

        // 初始化统计
        updateEditorStats()

        // 初始化当前语言
        currentLanguage.value = initialLanguage
    }
})

onBeforeUnmount(() => {
    editor?.dispose()
    stopAutoRefresh()
})

watch(() => props.show, async (newVal) => {
    if (newVal && props.file) {
        await loadFileContent()
        editor?.setValue(content.value)
    }
})

watch(() => props.file?.name, (newName) => {
    if (editor && newName) {
        const model = editor.getModel()
        if (model) {
            monaco.editor.setModelLanguage(model, getLanguageFromFileName(newName))
        }
    }
})

async function loadFileContent() {
    if (!props.file) return
    try {
        // 如果提供了初始内容，直接使用
        if (props.initialContent !== undefined) {
            content.value = props.initialContent
            return
        }
        // 否则从文件系统读取
        const res = await fileAPI.getFileContent(props.file.path)
        content.value = res.data.content
        error.value = null
    } catch (err: any) {
        error.value = `打开文件失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}

async function saveFile() {
    if (!props.file) return
    try {
        await fileAPI.saveFileContent(props.file.path, content.value)
        error.value = null
        emit('saved')
        emit('update:show', false)
    } catch (err: any) {
        error.value = `保存文件失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}

async function formatCode() {
    if (!editor) return
    try {
        const content = editor.getValue()
        let formatted: string

        const ext = props.file?.name.split('.').pop()?.toLowerCase()
        if (ext === 'json') {
            formatted = JSON.stringify(JSON.parse(content), null, 2)
            editor.setValue(formatted)
        } else {
            const action = editor.getAction('editor.action.formatDocument')
            if (action) {
                await action.run()
            }
        }
        toast.success('格式化成功')
    } catch (e) {
        error.value = '格式化失败，请检查文件格式'
    }
}

function copyAll() {
    if (!editor) return
    const content = editor.getValue()
    navigator.clipboard.writeText(content)
        .then(() => toast.success('已复制到剪贴板'))
        .catch(() => error.value = '复制失败')
}

function downloadFile() {
    if (!props.file) return
    const content = editor?.getValue() || ''
    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = props.file.name
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    toast.success('开始下载')
}

async function resetChanges() {
    if (!editor || !props.file) return
    if (await confirm.warning('重置更改', '确定要放弃所有更改吗？此操作不可撤销。')) {
        await loadFileContent()
        editor.setValue(content.value)
        toast.info('已重置更改')
    }
}

const autoSave = ref(false)
let autoSaveTimer: ReturnType<typeof setTimeout> | null = null

watch(() => content.value, () => {
    if (autoSave.value && props.file) {
        if (autoSaveTimer) clearTimeout(autoSaveTimer)
        autoSaveTimer = setTimeout(() => {
            saveFile()
        }, 1000)
    }
    updateEditorStats()
})

onBeforeUnmount(() => {
    if (autoSaveTimer) clearTimeout(autoSaveTimer)
})

// 切换小地图
function toggleMinimap() {
    if (editor) {
        showMinimap.value = !showMinimap.value
        editor.updateOptions({ minimap: { enabled: showMinimap.value } })
        toast.info(showMinimap.value ? '已显示小地图' : '已隐藏小地图')
    }
}

// 处理自动刷新
function toggleAutoRefresh() {
    autoRefresh.value = !autoRefresh.value
    if (autoRefresh.value) {
        startAutoRefresh()
    } else {
        stopAutoRefresh()
    }
}

function startAutoRefresh() {
    if (refreshTimer) clearInterval(refreshTimer)
    refreshTimer = setInterval(async () => {
        if (!props.file) return
        try {
            const res = await fileAPI.getFileContent(props.file.path)
            if (res.data.content !== content.value) {
                if (editor) {
                    const position = editor.getPosition()
                    editor.setValue(res.data.content)
                    editor.setPosition(position || { lineNumber: 1, column: 1 })
                }
                toast.info('文件内容已更新')
            }
        } catch (err) {
            console.error('自动刷新失败:', err)
        }
    }, refreshInterval.value * 1000)
}

function stopAutoRefresh() {
    if (refreshTimer) {
        clearInterval(refreshTimer)
        refreshTimer = null
    }
}

// 添加清空内容的函数
async function clearContent() {
    if (!editor) return
    if (await confirm.warning('清空内容', '确定要清空所有内容吗？此操作不可撤销。')) {
        editor.setValue('')
        toast.info('内容已清空')
    }
}

// 更新语言的函数
function changeLanguage(newLanguage: string) {
    if (editor) {
        const model = editor.getModel()
        if (model) {
            monaco.editor.setModelLanguage(model, newLanguage)
            currentLanguage.value = newLanguage
            toast.success(`已切换到 ${newLanguage} 语法高亮`)
        }
    }
}

// 在编辑器内容变化时更新统计
function updateEditorStats() {
    if (editor) {
        const model = editor.getModel()
        if (model) {
            lineCount.value = model.getLineCount()
            characterCount.value = model.getValueLength()
        }
    }
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box w-11/12 max-w-5xl h-4/5 flex flex-col pb-0 px-0">
            <!-- 标题栏和主要操作按钮 -->
            <div class="flex justify-between items-center mb-4 px-4">
                <h3 class="font-bold text-lg">编辑文件: {{ file?.name }}</h3>
                <div class="flex items-center gap-2">
                    <div class="form-control">
                        <label class="label cursor-pointer">
                            <span class="label-text mr-2">自动保存</span>
                            <input type="checkbox" v-model="autoSave" class="toggle toggle-primary toggle-sm" />
                        </label>
                    </div>
                    <IconButton @click="$emit('update:show', false)">
                        <RiCloseLine />
                        取消
                    </IconButton>
                    <IconButton variant="primary" @click="saveFile">
                        <RiSaveLine />
                        保存
                    </IconButton>
                </div>
            </div>

            <!-- 次要工具栏 -->
            <div
                class="flex justify-between items-center mx-4 py-2 mb-4 bg-base-200 rounded-lg px-4 shadow-sm transition-colors">
                <div class="flex gap-2">
                    <IconButton @click="formatCode" title="格式化 (Shift+Alt+F)">
                        <RiCodeLine />
                        格式化
                    </IconButton>
                    <IconButton @click="copyAll" title="复制全部">
                        <RiFileCopyLine />
                        复制
                    </IconButton>
                    <IconButton @click="downloadFile" title="下载文件">
                        <RiDownloadLine />
                        下载
                    </IconButton>
                    <IconButton variant="warning" @click="resetChanges" title="重置更改">
                        <RiRestartLine />
                        重置
                    </IconButton>
                    <IconButton variant="error" @click="clearContent" title="清空内容">
                        <RiDeleteBinLine />
                        清空
                    </IconButton>
                    <IconButton @click="toggleMinimap" :title="showMinimap ? '隐藏小地图' : '显示小地图'">
                        <RiMapLine />
                        {{ showMinimap ? '隐藏地图' : '显示地图' }}
                    </IconButton>
                    <IconButton :icon="autoRefresh ? 'ri-time-fill' : 'ri-time-line'"
                        :variant="autoRefresh ? 'primary' : 'default'" @click="toggleAutoRefresh"
                        :title="`自动刷新 (${refreshInterval}秒)`">
                        {{ autoRefresh ? '停止刷新' : '自动刷新' }}
                    </IconButton>
                    <div v-if="autoRefresh" class="flex items-center gap-2">
                        <select v-model="refreshInterval" class="select select-bordered select-sm"
                            @change="startAutoRefresh">
                            <option value="3">3秒</option>
                            <option value="5">5秒</option>
                            <option value="10">10秒</option>
                            <option value="30">30秒</option>
                        </select>
                    </div>
                </div>
            </div>

            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <RiAlertLine />
                    <span>{{ error }}</span>
                </div>
            </div>

            <div ref="editorContainer" class="flex-1"></div>

            <!-- 状态栏 -->
            <div class="flex justify-between items-center py-0 px-4 bg-base-300 text-sm mt-0 rounded-b-lg">
                <!-- 左侧状态信息 -->
                <div class="flex items-center gap-4 text-base-content/70">
                    <span>{{ lineCount }} 行</span>
                    <span>{{ characterCount }} 个字符</span>
                </div>

                <!-- 语言选择器 -->
                <div class="dropdown dropdown-top dropdown-end">
                    <label tabindex="0" class="btn btn-ghost btn-sm normal-case">
                        {{supportedLanguages.find(lang => lang.id === currentLanguage)?.name || currentLanguage}}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20"
                            fill="currentColor">
                            <path fill-rule="evenodd"
                                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                                clip-rule="evenodd" />
                        </svg>
                    </label>
                    <ul tabindex="0"
                        class="dropdown-content menu p-2 shadow bg-base-200 rounded-box w-52 max-h-60 overflow-y-auto">
                        <li v-for="lang in supportedLanguages" :key="lang.id">
                            <a @click="changeLanguage(lang.id)" :class="{ 'active': currentLanguage === lang.id }">
                                {{ lang.name }}
                            </a>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </dialog>
</template>

<style scoped>
.dropdown-content {
    bottom: 100%;
    margin-bottom: 4px;
}
</style>
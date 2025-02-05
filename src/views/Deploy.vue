<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useToast } from '../composables/useToast'
import Modal from '../components/Modal.vue'
import LogViewer from '../components/LogViewer.vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'

interface Project {
    id: number
    name: string
    type: string
    git_repo: string
    branch: string
    build_cmd: string
    output_dir: string
    domain: string
    path: string
    port: number
    status: string
    last_deploy: string
    environment: { key: string, value: string }[]
}

// 创建默认项目模板
const createDefaultProject = (): Project => ({
    id: 0,
    name: '',
    type: 'static',
    git_repo: '',
    branch: 'main',
    build_cmd: '',
    output_dir: '',
    domain: '',
    path: '/',
    port: 3000,
    status: 'stopped',
    last_deploy: '',
    environment: []
})

const projects = ref<Project[]>([])
const showModal = ref(false)
const editingProject = ref<Project>(createDefaultProject())
const showBuildLogs = ref(false)
const buildLogs = ref<string[]>([])
const currentProject = ref<Project | null>(null)
const showDeleteConfirm = ref(false)
const projectToDelete = ref<number | null>(null)
const loading = ref(false)
const toast = useToast()

// 获取所有项目
const fetchProjects = async () => {
    try {
        const res = await axios.get('/web_api/deploy/projects')
        projects.value = res.data
    } catch (error: any) {
        toast.error('获取项目列表失败: ' + error.message)
    }
}

// 保存项目
const saveProject = async () => {
    try {
        loading.value = true
        const project = editingProject.value

        if (project.id) {
            await axios.put(`/web_api/deploy/projects/${project.id}`, project)
            toast.success('项目更新成功')
        } else {
            await axios.post('/web_api/deploy/projects', project)
            toast.success('项目创建成功')
        }

        showModal.value = false
        await fetchProjects()
    } catch (error: any) {
        const message = error.response?.data?.error || error.message
        toast.error('保存失败: ' + message)
    } finally {
        loading.value = false
    }
}

// 构建项目
const buildProject = async (project: Project) => {
    currentProject.value = project
    buildLogs.value = []
    showBuildLogs.value = true

    const eventSource = new EventSource(`/web_api/deploy/projects/${project.id}/build`)

    eventSource.onmessage = (event) => {
        buildLogs.value.push(event.data)
        buildLogs.value = [...buildLogs.value]
    }

    eventSource.onerror = () => {
        eventSource.close()
        fetchProjects()
    }
}

// 删除项目
const handleDelete = async () => {
    if (!projectToDelete.value) return

    try {
        await axios.delete(`/web_api/deploy/projects/${projectToDelete.value}`)
        toast.success('项目删除成功')
        await fetchProjects()
    } catch (error: any) {
        toast.error('删除失败: ' + error.message)
    } finally {
        showDeleteConfirm.value = false
        projectToDelete.value = null
    }
}

// 编辑项目
const editTask = (project: Project) => {
    editingProject.value = { ...project }
    showModal.value = true
}

// 新建项目
const createProject = () => {
    editingProject.value = createDefaultProject()
    showModal.value = true
}

onMounted(fetchProjects)
</script>

<template>
    <div class="p-6">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-bold">项目部署</h1>
            <button class="btn btn-primary" @click="createProject">
                <i class="ri-add-line mr-1"></i>新建项目
            </button>
        </div>

        <!-- 项目列表 -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="project in projects" :key="project.id" class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title">
                        {{ project.name }}
                        <div class="badge" :class="{
                            'badge-success': project.status === 'running',
                            'badge-error': project.status === 'failed',
                            'badge-warning': project.status === 'stopped'
                        }">
                            {{ project.status }}
                        </div>
                    </h2>
                    <p class="text-sm opacity-70">{{ project.domain }}{{ project.path }}</p>
                    <div class="divider my-2"></div>
                    <div class="space-y-2 text-sm">
                        <p><span class="font-semibold">类型：</span>{{ project.type }}</p>
                        <p><span class="font-semibold">仓库：</span>{{ project.git_repo }}</p>
                        <p><span class="font-semibold">分支：</span>{{ project.branch }}</p>
                        <p v-if="project.last_deploy">
                            <span class="font-semibold">最后部署：</span>
                            {{ new Date(project.last_deploy).toLocaleString() }}
                        </p>
                    </div>
                    <div class="card-actions justify-end mt-4">
                        <button class="btn btn-sm btn-ghost" @click="editTask(project)">
                            <i class="ri-edit-line"></i>
                        </button>
                        <button class="btn btn-sm btn-primary" @click="buildProject(project)">
                            <i class="ri-rocket-line mr-1"></i>部署
                        </button>
                        <button class="btn btn-sm btn-ghost text-error"
                            @click="projectToDelete = project.id; showDeleteConfirm = true">
                            <i class="ri-delete-bin-line"></i>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- 项目表单 -->
        <Modal v-model:show="showModal" :title="editingProject.id ? '编辑项目' : '新建项目'">
            <form @submit.prevent="saveProject" class="space-y-4 p-4">
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">项目名称</span>
                    </label>
                    <input type="text" v-model="editingProject.name" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">项目类型</span>
                    </label>
                    <select v-model="editingProject.type" class="select select-bordered" required>
                        <option value="static">静态网站</option>
                        <option value="node">Node.js</option>
                        <option value="go">Go</option>
                        <option value="python">Python</option>
                    </select>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Git仓库</span>
                    </label>
                    <input type="text" v-model="editingProject.git_repo" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">分支</span>
                    </label>
                    <input type="text" v-model="editingProject.branch" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">构建命令</span>
                    </label>
                    <input type="text" v-model="editingProject.build_cmd" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">输出目录</span>
                    </label>
                    <input type="text" v-model="editingProject.output_dir" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">域名</span>
                    </label>
                    <input type="text" v-model="editingProject.domain" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">URL路径</span>
                    </label>
                    <input type="text" v-model="editingProject.path" class="input input-bordered" required />
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">服务端口</span>
                    </label>
                    <input type="number" v-model="editingProject.port" class="input input-bordered"
                        :required="editingProject.type !== 'static'" />
                </div>

                <!-- 环境变量 -->
                <div class="space-y-2">
                    <label class="label">
                        <span class="label-text">环境变量</span>
                        <button type="button" class="btn btn-sm btn-ghost"
                            @click="editingProject.environment.push({ key: '', value: '' })">
                            <i class="ri-add-line"></i>
                        </button>
                    </label>
                    <div v-for="(env, index) in editingProject.environment" :key="index" class="flex gap-2">
                        <input type="text" v-model="env.key" placeholder="KEY" class="input input-bordered flex-1" />
                        <input type="text" v-model="env.value" placeholder="VALUE"
                            class="input input-bordered flex-1" />
                        <button type="button" class="btn btn-ghost btn-sm"
                            @click="editingProject.environment.splice(index, 1)">
                            <i class="ri-delete-bin-line"></i>
                        </button>
                    </div>
                </div>

                <div class="modal-action">
                    <button type="button" class="btn" @click="showModal = false">取消</button>
                    <button type="submit" class="btn btn-primary" :disabled="loading">
                        {{ loading ? '保存中...' : '保存' }}
                    </button>
                </div>
            </form>
        </Modal>

        <!-- 构建日志 -->
        <Modal v-model:show="showBuildLogs" :title="`${currentProject?.name} 部署日志`">
            <LogViewer :logs="buildLogs" />
        </Modal>

        <!-- 删除确认 -->
        <ConfirmDialog v-model:show="showDeleteConfirm" title="确认删除" message="删除项目将同时删除所有相关文件和配置，此操作不可撤销。是否继续？"
            type="error" confirm-text="删除" @confirm="handleDelete" />
    </div>
</template>
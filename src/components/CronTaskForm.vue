<script setup lang="ts">
import { ref } from 'vue'
import Alert from './Alert.vue'

interface CronTask {
    id: number
    name: string
    command: string
    schedule: string
    description: string
    enabled: boolean
    last_run?: string
    next_run?: string
}

const props = defineProps<{
    task: CronTask
    isEditing: boolean
}>()

const emit = defineEmits<{
    (e: 'submit', task: CronTask): void
    (e: 'cancel'): void
}>()

const formError = ref('')
const fieldErrors = ref<Record<string, string>>({})
const showCronHelp = ref(false)

// 预定义的任务模板
const taskTemplates = {
    systemClean: {
        name: '系统清理',
        command: 'apt clean && apt autoremove -y',
        description: '定期清理系统包缓存'
    },
    databaseBackup: {
        name: '数据库备份',
        command: 'mysqldump -u root -p database > backup.sql',
        description: '数据库定时备份'
    },
    logClean: {
        name: '日志清理',
        command: 'find /var/log -type f -name "*.log" -mtime +30 -delete',
        description: '清理30天前的日志文件'
    },
    systemUpdate: {
        name: '系统更新',
        command: 'apt update && apt upgrade -y',
        description: '自动更新系统包'
    }
}

// 应用任务模板
const applyTemplate = (templateName: keyof typeof taskTemplates) => {
    const template = taskTemplates[templateName]
    props.task.name = template.name
    props.task.command = template.command
    props.task.description = template.description
}

const handleSubmit = () => {
    emit('submit', props.task)
}
</script>

<template>
    <form @submit.prevent="handleSubmit">
        <div class="flex flex-wrap gap-2 mb-4">
            <button type="button" class="btn btn-sm" @click="() => applyTemplate('systemClean')">
                系统清理
            </button>
            <button type="button" class="btn btn-sm" @click="() => applyTemplate('databaseBackup')">
                数据库备份
            </button>
            <button type="button" class="btn btn-sm" @click="() => applyTemplate('logClean')">
                日志清理
            </button>
            <button type="button" class="btn btn-sm" @click="() => applyTemplate('systemUpdate')">
                系统更新
            </button>
        </div>

        <div class="form-control mb-4">
            <label class="label">
                <span class="label-text">任务名称</span>
            </label>
            <input type="text" v-model="task.name" class="input input-bordered"
                :class="{ 'input-error': fieldErrors.name }" required />
            <label v-if="fieldErrors.name" class="label">
                <span class="label-text-alt text-error">{{ fieldErrors.name }}</span>
            </label>
        </div>

        <div class="form-control mb-4">
            <label class="label">
                <span class="label-text">执行命令</span>
            </label>
            <input type="text" v-model="task.command" class="input input-bordered"
                :class="{ 'input-error': fieldErrors.command }" required />
            <label v-if="fieldErrors.command" class="label">
                <span class="label-text-alt text-error">{{ fieldErrors.command }}</span>
            </label>
        </div>

        <div class="form-control mb-4">
            <label class="label">
                <span class="label-text">定时表达式</span>
                <span class="label-text-alt">
                    <a href="#" @click.prevent="showCronHelp = true" class="link">
                        帮助
                    </a>
                </span>
            </label>
            <input type="text" v-model="task.schedule" class="input input-bordered"
                :class="{ 'input-error': fieldErrors.schedule }" required placeholder="0 */5 * * * *" />
            <div class="flex flex-col gap-1 mt-1">
                <label class="label py-0">
                    <span class="label-text-alt text-base-content/70">格式: 秒 分 时 日 月 星期</span>
                </label>
                <div class="flex flex-wrap gap-2 mt-1">
                    <button type="button" class="btn btn-xs" @click="task.schedule = '0 0 0 * * *'">
                        每天0点
                    </button>
                    <button type="button" class="btn btn-xs" @click="task.schedule = '0 0 3 * * *'">
                        每天3点
                    </button>
                    <button type="button" class="btn btn-xs" @click="task.schedule = '0 */30 * * * *'">
                        每30分钟
                    </button>
                    <button type="button" class="btn btn-xs" @click="task.schedule = '0 0 * * * *'">
                        每小时
                    </button>
                    <button type="button" class="btn btn-xs" @click="task.schedule = '0 0 12 * * 1-5'">
                        工作日12点
                    </button>
                </div>
                <label v-if="fieldErrors.schedule" class="label py-0">
                    <span class="label-text-alt text-error">{{ fieldErrors.schedule }}</span>
                </label>
            </div>
        </div>

        <div class="form-control mb-4">
            <label class="label">
                <span class="label-text">描述</span>
            </label>
            <textarea v-model="task.description" class="textarea textarea-bordered" rows="3">
            </textarea>
        </div>

        <!-- 通用错误信息显示区域 -->
        <Alert v-if="formError" type="error" :message="formError" class="mb-4" />

        <div class="modal-action">
            <button type="button" class="btn" @click="emit('cancel')">取消</button>
            <button type="submit" class="btn btn-primary">保存</button>
        </div>
    </form>
</template>
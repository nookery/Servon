<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { githubAPI } from '../../api/github_api'
import ActionButton from '../../components/ActionButton.vue'

const formData = reactive({
    name: 'Servon',
    description: 'Servon GitHub integration for automation',
    base_url: 'http://43.142.208.212:9754'
})

const isLocalhost = computed(() => {
    const url = formData.base_url.toLowerCase()
    return url.includes('localhost') || url.includes('127.0.0.1')
})

const error = ref<string | null>(null)
const showUrlWarning = ref(false)

const emit = defineEmits<{
    success: []
}>()

const handleSubmit = async () => {
    if (isLocalhost.value) {
        showUrlWarning.value = true
        error.value = '请使用公网可访问的URL，localhost或127.0.0.1将无法接收GitHub的webhook请求'
        return
    }
    error.value = null

    // 将 reactive 对象转换为普通对象
    const data = {
        name: formData.name,
        description: formData.description || undefined,
        base_url: formData.base_url
    }

    console.log('Sending setup request with data:', data)

    try {
        const response = await githubAPI.setup(data)

        const div = document.createElement('div')
        div.innerHTML = response.data
        document.body.appendChild(div)

        const form = div.querySelector('#github-form') as HTMLFormElement
        if (form) {
            form.submit()
        } else {
            throw new Error('表单创建失败')
        }

        setTimeout(() => {
            document.body.removeChild(div)
        }, 2000)

        emit('success')
    } catch (err: any) {
        error.value = err.response?.data?.message ||
            err.response?.data?.error ||
            (typeof err.response?.data === 'string' ? err.response.data : null) ||
            err.message ||
            '创建失败'
    }
}
</script>

<template>
    <form @submit.prevent="handleSubmit">
        <div class="form-control w-full flex flex-row gap-3">
            <label class="label">
                <span class="label-text">名称</span>
            </label>
            <input v-model="formData.name" type="text" class="input input-bordered w-full" required />
        </div>

        <div class="form-control w-full mt-4 hidden">
            <label class="label">
                <span class="label-text">描述</span>
            </label>
            <textarea v-model="formData.description" class="textarea textarea-bordered"></textarea>
        </div>

        <div class="form-control w-full mt-4">
            <label class="label">
                <span class="label-text">公网访问地址</span>
                <span class="label-text-alt text-warning">必须是公网可访问的URL</span>
            </label>
            <input v-model="formData.base_url" type="url" placeholder="例如: https://your-domain.com"
                class="input input-bordered w-full" required />
            <label v-if="isLocalhost" class="label">
                <span class="label-text-alt text-error">
                    localhost或127.0.0.1无法接收GitHub的webhook请求！
                </span>
            </label>
        </div>

        <!-- 错误提示横幅 -->
        <div v-if="error" class="alert alert-error mt-4">
            <i class="ri-error-warning-line"></i>
            <span>{{ error }}</span>
        </div>

        <div class="modal-action">
            <slot name="actions" :submitting="false">
                <ActionButton type="submit" variant="primary">创建</ActionButton>
            </slot>
        </div>
    </form>
</template>
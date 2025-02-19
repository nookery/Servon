<script setup lang="ts">
import { reactive, ref } from 'vue'
import { githubAPI } from '../../api/github'
import ActionButton from '../../components/ActionButton.vue'

const formData = reactive({
    name: 'Servon',
    description: 'Servon GitHub integration for automation'
})
const error = ref<string | null>(null)

const emit = defineEmits<{
    success: []
}>()

const handleSubmit = async () => {
    error.value = null
    try {
        const response = await githubAPI.setup({
            name: formData.name,
            description: formData.description || undefined,
        })

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
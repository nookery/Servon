<template>
    <div>
        <button @click="showModal" class="btn btn-ghost btn-circle" title="GitHub集成">
            <i class="ri-github-line text-xl"></i>
        </button>

        <dialog ref="modalRef" class="modal">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">创建 GitHub App</h3>
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

                    <div class="modal-action">
                        <button type="button" class="btn" @click="closeModal">取消</button>
                        <button type="submit" class="btn btn-primary">创建</button>
                    </div>
                </form>
            </div>
        </dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { githubAPI } from '../api/github'

const modalRef = ref<HTMLDialogElement | null>(null)
const formData = reactive({
    name: 'Servon',
    description: 'Servon GitHub integration for automation'
})

const showModal = () => {
    modalRef.value?.showModal()
}

const closeModal = () => {
    modalRef.value?.close()
}

const handleGitHubSuccess = () => {
    // 可以添加成功提示
}

const handleGitHubError = (error: string) => {
    alert('启动GitHub集成失败: ' + error)
}

const handleSubmit = async () => {
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

        handleGitHubSuccess()
        closeModal()
    } catch (error: any) {
        handleGitHubError(error.response?.data?.error || error.message)
    }
}
</script>
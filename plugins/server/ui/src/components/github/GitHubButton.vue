<script setup lang="ts">
import { ref } from 'vue'
import IconButton from '../IconButton.vue'
import ActionButton from '../ActionButton.vue'
import GitHubSetupForm from './GitHubSetupForm.vue'

const modalRef = ref<HTMLDialogElement | null>(null)

const showModal = () => {
    modalRef.value?.showModal()
}

const closeModal = () => {
    modalRef.value?.close()
}

const handleGitHubSuccess = () => {
    closeModal()
}
</script>

<template>
    <div>
        <IconButton @click="showModal" icon="ri-github-line" size="lg" variant="ghost" circle title="GitHub集成" />

        <dialog ref="modalRef" class="modal">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">创建 GitHub App</h3>

                <GitHubSetupForm @success="handleGitHubSuccess">
                    <template #actions>
                        <ActionButton type="button" @click="closeModal">取消</ActionButton>
                        <ActionButton type="submit" variant="primary">创建</ActionButton>
                    </template>
                </GitHubSetupForm>
            </div>
        </dialog>
    </div>
</template>
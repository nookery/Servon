<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import GitHubIntegration from '../components/github/GitHubIntegration.vue'

const currentTab = ref('github')
</script>

<template>
    <PageContainer>
        <template #tabs>
            <div class="tabs tabs-box">
                <a class="tab tab-bordered" :class="{ 'tab-active': currentTab === 'github' }"
                    @click="currentTab = 'github'">
                    <i class="ri-github-fill mr-2"></i>
                    GitHub
                </a>
                <a class="tab tab-bordered" :class="{ 'tab-active': currentTab === 'gitlab' }"
                    @click="currentTab = 'gitlab'">
                    <i class="ri-gitlab-fill mr-2"></i>
                    GitLab
                </a>
            </div>
        </template>

        <!-- GitHub 集成内容 -->
        <GitHubIntegration v-if="currentTab === 'github'" />

        <!-- GitLab 集成内容 -->
        <div v-else-if="currentTab === 'gitlab'" class="py-4">
            <div class="flex flex-col items-center justify-center py-12 text-base-content/60">
                <i class="ri-gitlab-fill text-6xl mb-4"></i>
                <p class="text-lg mb-2">GitLab 集成配置</p>
                <p class="text-sm">这里将显示 GitLab 相关的配置和状态</p>
            </div>
        </div>
    </PageContainer>
</template>
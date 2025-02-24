<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import GitHubIntegration from '../components/github/GitHubIntegration.vue'
import GitLabIntegration from '../components/gitlab/GitLabIntegration.vue'
import { RiGithubFill, RiGitlabFill } from '@remixicon/vue'

const activeTab = ref('github')

const tabs = [
    { key: 'github', title: 'GitHub', icon: RiGithubFill },
    { key: 'gitlab', title: 'GitLab', icon: RiGitlabFill }
]
</script>

<template>
    <PageContainer title="集成管理" :tabs="tabs" v-model="activeTab">
        <template #github>
            <GitHubIntegration />
        </template>

        <template #gitlab>
            <GitLabIntegration />
        </template>
    </PageContainer>
</template>
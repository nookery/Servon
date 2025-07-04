import { createRouter, createWebHistory } from 'vue-router'
import CronTasks from './views/CronTasks.vue'
import Dashboard from './views/Dashboard.vue'
import Software from './views/Software.vue'
import Processes from './views/Processes.vue'
import FilesPage from './views/FilesPage.vue'
import Ports from './views/Ports.vue'
import Users from './views/Users.vue'
import DeployLogs from './views/DeployLogs.vue'
import Integrations from './views/Integrations.vue'
import DataPage from './views/DataPage.vue'
import Logs from './views/Logs.vue'
import ProjectsPage from './views/ProjectsPage.vue'
const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/dashboard'
        },
        {
            path: '/dashboard',
            component: Dashboard
        },
        {
            path: '/software',
            component: Software
        },
        {
            path: '/processes',
            component: Processes
        },
        {
            path: '/files',
            component: FilesPage
        },
        {
            path: '/ports',
            component: Ports
        },
        {
            path: '/users',
            component: Users
        },
        {
            path: '/cron',
            name: 'cron',
            component: CronTasks,
            meta: {
                title: '定时任务'
            }
        },
        {
            path: '/deploy-logs',
            name: 'deploy-logs',
            component: DeployLogs
        },
        {
            path: '/integrations',
            name: 'integrations',
            component: Integrations
        },
        {
            path: '/data',
            name: 'data',
            component: DataPage
        },
        {
            path: '/logs',
            name: 'Logs',
            component: Logs,
            meta: {
                title: '日志管理',
                icon: 'ri-file-list-line'
            }
        },
        {
            path: '/projects',
            name: 'Projects',
            meta: {
                title: '项目管理',
                icon: 'ri-folder-line'
            },
            component: ProjectsPage
        }
    ]
})

export default router 
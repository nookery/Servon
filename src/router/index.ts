import { createRouter, createWebHistory } from 'vue-router'
import CronTasks from '../views/CronTasks.vue'
import Dashboard from '../views/Dashboard.vue'
import Software from '../views/Software.vue'
import Processes from '../views/Processes.vue'
import Files from '../views/Files.vue'
import Ports from '../views/Ports.vue'
import Deploy from '../views/Deploy.vue'

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
            component: Files
        },
        {
            path: '/ports',
            component: Ports
        },
        {
            path: '/deploy',
            component: Deploy
        },
        {
            path: '/cron',
            name: 'cron',
            component: CronTasks,
            meta: {
                title: '定时任务'
            }
        }
    ]
})

export default router 
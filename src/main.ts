import { createApp } from 'vue'
import './app.css'
import App from './App.vue'
import router from './router'
import 'remixicon/fonts/remixicon.css'

const app = createApp(App)
app.use(router)
app.mount('#app')

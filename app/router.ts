import { createMemoryHistory, createRouter } from 'vue-router'

import login from './login.vue'
import connected from './connected.vue'

const routes = [
    { name: "login_page", path: '/', component: login },
    { name: "connected_page", path: '/connected', component: connected },
]

const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

export default router
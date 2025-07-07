import { createMemoryHistory, createRouter } from 'vue-router'

import login from './views/login.vue'
import connected from './views/connected.vue'
import arena from './views/arena.vue'

const routes = [
    { name: "login_page", path: '/', component: login },
    { name: "connected_page", path: '/connected', component: connected },
    { name: "arena_page", path: '/arena', component: arena },
]

const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

export default router
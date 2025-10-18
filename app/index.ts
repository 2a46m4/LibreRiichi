import { createApp, ref, Ref } from 'vue'
import App from './views/app.vue'
import { createPinia, defineStore } from 'pinia'
import './index.css'
import { Connection, websocket_address } from './messaging/connection'
import { ServerMessageBus } from './messaging/event_handler'
import { createMemoryHistory, createRouter } from 'vue-router'

import login from './views/login.vue'
import connected from './views/connected.vue'
import arena from './views/arena.vue'
import { IncomingMessage, validate_message } from './messaging/message'

const routes = [
  { name: 'login_page', path: '/', component: login },
  { name: 'connected_page', path: '/connected', component: connected },
  { name: 'arena_page', path: '/arena', component: arena },
]

export const router = createRouter({
  history: createMemoryHistory(),
  routes,
})

const pinia = createPinia()
const app = createApp(App)
app.use(router).use(pinia).mount('#app')

export const use_websocket_state = defineStore(
  'websocket_state',
  (): {
    conn: Connection
    ready: Ref<boolean>
  } => {
    let conn = new Connection(
      new WebSocket(websocket_address),
      (data: MessageEvent) => {
        let msg = JSON.parse(data.data) as IncomingMessage
        if (validate_message(msg).isValid) {
          ServerMessageBus.handle(msg)
        }
      },
    )
    let ready = ref(false)
    conn.wait_until_ready().then(() => {
      ready.value = true
    })

    return { conn, ready }
  },
)

export const use_player_state = defineStore(
  'player_state',
  (): {
    username: Ref<string>
  } => {
    let username = ref('')
    return { username }
  },
)

export const use_room_state = defineStore('room_state', () => {
  let room_name = ref('')
  let room_set = ref(false)
  return { room_name, room_set }
})
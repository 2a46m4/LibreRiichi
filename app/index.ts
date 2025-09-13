import {createApp, ref, Ref} from "vue";
import App from "./views/app.vue";
import router from './router'
import {createPinia, defineStore} from "pinia";
import './index.css'
import {Connection, websocket_address} from "./messaging/connection";
import {ServerMessageBus} from "./messaging/event_handler";
import {MessageRouter} from "./messaging/message_router";

const pinia = createPinia();
const app = createApp(App);
app.use(router).use(pinia).mount("#app");

export const use_websocket_state = defineStore('websocket_state', (): {
    conn: Connection;
    msg_router: MessageRouter;
    ready: Ref<boolean>;
} => {
    let conn = new Connection(new WebSocket(websocket_address), ServerMessageBus.handle.bind(ServerMessageBus))
    let ready = ref(false)
    conn.wait_until_ready().then(() => {ready.value = true})

    let msg_router = new MessageRouter()
    ServerMessageBus.register(msg_router.match_message.bind(msg_router))
    return {conn, msg_router, ready}
})

export const use_player_state = defineStore('player_state', (): {
    username: Ref<string>;
} => {
    let username = ref("")
    return {username}
})
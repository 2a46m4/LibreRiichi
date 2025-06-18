import {createApp, ref} from "vue";
import App from "./app.vue";
import router from './router'
import {createPinia} from "pinia";
import './index.css'

const pinia = createPinia();
const app = createApp(App);
app.use(router).use(pinia).mount("#app");
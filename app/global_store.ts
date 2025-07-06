import {defineStore} from "pinia";
import {Application} from "./application";
import {LoginState} from "./app_state/login_state";

export const useGlobalStore = defineStore('globalStore', ()=>{
    let application = new Application()
    return {application}
})
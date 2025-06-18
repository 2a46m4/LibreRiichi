import {defineStore} from "pinia";
import {Application} from "./app_state";

export const useGlobalStore = defineStore('globalStore', ()=>{
    let application = new Application()
    return {application}
})
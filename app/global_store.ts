import {defineStore} from "pinia";
import {Ref, ref} from "vue";
import {Connection} from "./connection";

export const useGlobalStore = defineStore('globalStore', ()=>{
    let connection: Ref<Connection | null> = ref(null)
    function setConnection(conn: Connection): void {
        connection.value = conn
    }
    function getConnection(): Connection {
        return connection.value
    }
    return {connection, setConnection, getConnection}
})
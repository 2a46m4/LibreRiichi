<script setup lang="ts">
import {computed, reactive, ref} from 'vue'
import {JoinArenaActionData} from "./message";
import {Connection} from "./connection";
import {useRoute, useRouter} from "vue-router";
import {useGlobalStore} from "./global_store";
import {BoxStyling, ButtonStyling, H1Styling, InputStyling} from "./styling";

const globalStore = useGlobalStore();

const user_name = ref('')

const router = useRouter()

async function connect() {
  console.log("connecting...")
  let connection = new Connection(user_name.value, new WebSocket("ws://localhost:3000/game"))
  await connection.WaitUntilReady()
  globalStore.setConnection(connection);
  await router.push({name: 'connected_page'})
}
</script>

<template>
  <div :class="BoxStyling">
    <h1 :class="H1Styling">LibreRiichi</h1>
    <p>Username</p>
    <input
        :class="InputStyling"
        v-model="user_name">
    <button
        :class="ButtonStyling"
        @click="connect">Connect</button>
  </div>
</template>

<style>

button {
  margin-left: 10px;
}
</style>

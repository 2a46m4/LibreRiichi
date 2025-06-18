<script setup lang="ts">
import {computed, reactive, ref} from 'vue'
import {JoinArenaActionData} from "./message";
import {Connection} from "./connection";
import {useRoute, useRouter} from "vue-router";
import {useGlobalStore} from "./global_store";
import {BoxStyling, ButtonStyling, H1Styling, InputStyling} from "./styling";

const globalStore = useGlobalStore();

const user_name = ref('')

const app = globalStore.application

const set_username = (username: string) => app.set_username(username)

async function connect() {
  set_username(app.username)
  app.connect()
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

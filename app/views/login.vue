<script setup lang="ts">
import {ref} from 'vue'
import {useGlobalStore} from "../global_store";
import {BoxStyling, ButtonStyling, H1Styling, InputStyling} from "../styling";
import Tile from "../components/tile.vue";
import {test} from "../assets/tiles";

const globalStore = useGlobalStore();
const app = globalStore.application

const user_name = ref('')

async function connect() {
  app.set_username(user_name.value)
  await app.action.connect()
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
      @click="connect"
      @keyup.enter="connect">Connect</button>
  </div>
  <Tile :tile_path="test"></Tile>
</template>

<style>

button {
  margin-left: 10px;
}
</style>

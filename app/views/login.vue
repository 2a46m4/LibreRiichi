<script setup lang="ts">
import {ref} from 'vue'
import {useGlobalStore} from "../global_store";
import {BoxStyling, ButtonStyling, H1Styling, InputStyling} from "../styling";
import TileComponent from "../components/tile_component.vue";
import {test, test2} from "../assets/tiles";
import ErrorDisplay from "../components/error_display.vue";

const globalStore = useGlobalStore();
const app = globalStore.application

const user_name = ref('')
const status = ref("")

async function connect() {
  app.set_username(user_name.value)
  try {
    await app.action.connect()
  } catch (error) {
    if (error instanceof Error) {
      status.value = error.message
    }
  }
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
    <ErrorDisplay :error="status" v-if="status.length !== 0"></ErrorDisplay>
  </div>

</template>

<style>

button {
  margin-left: 10px;
}
</style>

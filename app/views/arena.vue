<script setup lang="ts">

import {BoxStyling, FlexBox, H1Styling, ULStyling} from "../styling";
import {useGlobalStore} from "../global_store";
import {ref, Ref} from "vue";
import ListItem from "../components/list_item.vue";

const store = useGlobalStore()
const app = store.application
const action = app.action

const players: Ref<string[]> = ref([])

async function get_arena_info() {
  let arena = await action.get_arena_info()
  players.value = arena.agents.map(x => x.name)
}

get_arena_info()

</script>

<template>
  <div :class="BoxStyling">
    <div :class="FlexBox">
      <h1 :class="H1Styling">Players</h1>
    </div>
    <ul :class="ULStyling" v-if="players.length !== 0">
      <ListItem v-for="player in players">{{ player }}</ListItem>
    </ul>
  </div>
</template>

<style scoped>

</style>
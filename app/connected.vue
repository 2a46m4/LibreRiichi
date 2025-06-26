<script setup lang="ts">
import {useGlobalStore} from "./global_store";
import {BoxStyling, ButtonStyling, H1Styling, InputStyling} from "./styling";
import {ref} from "vue";
import ListItem from "./components/list_item.vue";

const globalStore = useGlobalStore();

const room_name = ref('')
const create_room_name = ref('')
const show_error = ref(false)
const avail_rooms = ref([])

function make_room() {

}

async function check_avail_rooms() {
  avail_rooms.value = await globalStore.application.list_rooms()
  console.log(avail_rooms.value)
}


async function find_room() {
  try {
    await globalStore.application.connect_room(
        room_name.value,
    )
  } catch (error) {
    console.log(error)
    show_error.value = true
  }
}

async function create_room() {
    await globalStore.application.create_room(create_room_name.value)
}

</script>

<template>
  <div :class="BoxStyling">
  <h1 :class="H1Styling">Join Room</h1>
  <p>Room Name</p>
  <input
      :class="InputStyling"
      v-model="room_name">
  <button
      :class="ButtonStyling"
      @click="find_room">Find</button>
  </div>
  <div :class="BoxStyling">
    <h1 :class="H1Styling">Create room</h1>
    <input
        :class="InputStyling"
        v-model="create_room_name">
    <button
        :class="ButtonStyling"
        @click="create_room">Create</button>
  </div>

  <div :class="BoxStyling">
    <h1 :class="H1Styling">Available rooms</h1>
    <button :class="ButtonStyling" @click="check_avail_rooms">Find</button>
    <ul v-if="avail_rooms.length !== 0">
      <ListItem v-for="avail_room in avail_rooms">{{ avail_room }}</ListItem>
    </ul>
  </div>


</template>

<style scoped>

</style>
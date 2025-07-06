<script setup lang="ts">
import {useGlobalStore} from "../global_store";
import {BoxStyling, ButtonStyling, FlexBox, H1Styling, InputStyling, ULStyling} from "../styling";
import {Ref, ref} from "vue";
import ListItem from "../components/list_item.vue";

const globalStore = useGlobalStore();
const app = globalStore.application
const action = app.action

const room_name = ref('')
const create_room_name = ref('')
const show_error = ref(false)
const avail_rooms: Ref<string[]> = ref([])

async function check_avail_rooms() {
  avail_rooms.value = await app.action.list_rooms()
}

async function find_room() {
  try {
    await action.connect_room(
        room_name.value,
    )
  } catch (error) {
    console.log(error)
    show_error.value = true
  }
}

async function create_room() {
    await action.create_room(create_room_name.value)
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
    <div :class="FlexBox">
      <h1 :class="H1Styling">Available rooms</h1>
      <button :class="ButtonStyling" @click="check_avail_rooms">Find</button>
      <br>
    </div>
    <ul :class="ULStyling" v-if="avail_rooms.length !== 0">
      <ListItem v-for="avail_room in avail_rooms">{{ avail_room }}</ListItem>
    </ul>
  </div>



</template>

<style scoped>

</style>
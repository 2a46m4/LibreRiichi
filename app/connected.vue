<script setup lang="ts">
import {JoinArenaActionData} from "./message";
import {useGlobalStore} from "./global_store";
import {BoxStyling, ButtonStyling, H1Styling, InputStyling} from "./styling";
import {ref} from "vue";
const globalStore = useGlobalStore();

const room_name = ref('')
const create_room_name = ref('')

function make_room() {

}

function check_avail_rooms() {}


function find_room() {
  try {
    globalStore.connection.Send(new JoinArenaActionData(room_name.value))
    globalStore.connection.Receive()
  } catch (error) {
    console.log(error)
  }
}

function create_room() {}

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
  </div>


</template>

<style scoped>

</style>
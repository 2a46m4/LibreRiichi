<script setup lang="ts">
import {JoinArenaActionData} from "./message";
import {useGlobalStore} from "./global_store";
const globalStore = useGlobalStore();
console.log("NMIGGSFD")
console.log(globalStore.getConnection().Send("SDKLFJSLDFKJ"));

function make_room() {

}

function find_room() {
  connection.Send(new JoinArenaActionData())
}
</script>

<template>
  <h2>Join Room</h2>
  <p>Room Name</p>
  <input v-model="room_name">
  <button @click="find_room">Find</button>

  <h2>Available rooms</h2>
  <h2>Create room</h2>
  <p>Room name</p>
  <input v-model="create_room_name">
  <button @click="make_room">Create</button>
</template>

<style scoped>

</style>
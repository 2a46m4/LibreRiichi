<script setup>
import { ref } from 'vue'
import {MessageType} from "./Message";

// Try to connect to websockets
let socket = new WebSocket("ws://localhost:3000/game");
socket.onopen = function(e) {
  console.log("Socket opened:", e)
}

socket.onmessage = function (e) {
  console.log(e.data)
  console.log("Message recv:", JSON.parse(e.data))
}

const message = ref('Hello World!')

const room_name = ref('')

const user_name = ref('')

const create_room_name = ref('')

function make_room() {

}

function reverseMessage() {
  // Access/mutate the value of a ref via
  // its .value property.
  message.value = message.value.split('').reverse().join('')
}

function notify() {
  alert('navigation was prevented.')
}

function find_room() {
  console.log(room_name.value)
  let message = JSON.stringify({
    message_type: MessageType.InitialMessageReturn,
    data: {
      name: "dabanya",
      room: "bigjoeroom"
    }
  })
  console.log("Sending message:", message)
  socket.send(message)
  console.log("Message sent:", message)
}
</script>

<template>
  <h1>LibreRiichi</h1>

  <h2>Join Room</h2>
  <p>
    Username
  </p>
  <input v-model="room_name">

  <p>Room Name</p>
    <input v-model="user_name">
    <button @click="find_room">Find</button>

  <h2>Available rooms</h2>
  <h2>Create room</h2>
  <p>Room name</p>
  <input v-model="create_room_name">
  <button @click="make_room">Create</button>

  <!--
    Note we don't need .value inside templates because
    refs are automatically "unwrapped" in templates.
  -->
<!--  <h1>{{ message }}</h1>-->

  <!--
    Bind to a method/function.
    The @click syntax is short for v-on:click.
  -->
<!--  <button @click="reverseMessage">Reverse Message</button>-->

  <!-- Can also be an inline expression statement -->
<!--  <button @click="message += '!'">Append "!"</button>-->

  <!--
    Vue also provides modifiers for common tasks
    such as e.preventDefault() and e.stopPropagation()
  -->
<!--  <a href="https://vuejs.org" @click.prevent="notify">-->
<!--    A link with e.preventDefault()-->
<!--  </a>-->
</template>

<style>
button {
  margin-left: 10px;
}
</style>

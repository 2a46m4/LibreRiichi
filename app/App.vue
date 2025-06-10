<script setup>
import { ref } from 'vue'



const message = ref('Hello World!')

const room_name = ref('')

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
}
</script>

<template>
  <h1>LibreRiichi</h1>

  <p>Join Room</p>
  <p>Room Name
    <input v-model="room_name">
    <button @click="find_room">Find</button>
  </p>

  <!--
    Note we don't need .value inside templates because
    refs are automatically "unwrapped" in templates.
  -->
  <h1>{{ message }}</h1>

  <!--
    Bind to a method/function.
    The @click syntax is short for v-on:click.
  -->
  <button @click="reverseMessage">Reverse Message</button>

  <!-- Can also be an inline expression statement -->
  <button @click="message += '!'">Append "!"</button>

  <!--
    Vue also provides modifiers for common tasks
    such as e.preventDefault() and e.stopPropagation()
  -->
  <a href="https://vuejs.org" @click.prevent="notify">
    A link with e.preventDefault()
  </a>
</template>

<style>
button {
  margin: 10px;
}
</style>
<script setup lang="ts">

let props = defineProps<{
  text: string
}>()
</script>

<template>
  <div class="container outline mt-2 p-2 bg-white rounded shadow-md text-red-400">
    <p class="text-xl">{{props.text}}</p>
  </div>
</template>

<style scoped></style>

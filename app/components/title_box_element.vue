<script setup lang="ts">
let props = defineProps<{
  text: string
}>()
</script>

<template>
  <div class="container outline mt-5 mb-5 p-10 bg-white rounded shadow-md">
    <h1 class="font-bold text-xl">{{ props.text }}</h1>
  </div>
</template>

<style scoped></style>

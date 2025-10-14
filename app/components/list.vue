<script setup lang="ts">
import ListItem from './list_item.vue'

defineProps<{
  items: string[]
}>()
</script>

<template>
  <ul class="container rounded mt-4">
    <ListItem v-for="item in items">{{ item }}</ListItem>
  </ul>
</template>

<style scoped></style>

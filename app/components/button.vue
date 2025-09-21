<script setup lang="ts">
defineProps<{
  condition: boolean
  on_click: () => void
  text: string
}>()
</script>

<template>
  <button
      class="outline p-1 rounded shadow-md bg-white"
      v-if="condition"
      @click="on_click"
  >
    {{text}}
  </button>
</template>

<style scoped>

</style>
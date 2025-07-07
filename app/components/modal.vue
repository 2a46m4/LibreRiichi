<script setup lang="ts">
const props = defineProps({
  show: Boolean
})
</script>

<template>
  <div class="flex items-center justify-center h-screen">
    <div class="bg-gray-200 p-4" v-if="props.show">CenteredHello</div>
  </div>
</template>

<style scoped>

</style>
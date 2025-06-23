<script setup lang="ts">
import {ListStyling} from "../styling";

const props = defineProps<{room_name: string}>()
</script>

<template>
<li :class="ListStyling">
  {{props.room_name}}
</li>
</template>

<style scoped>

</style>
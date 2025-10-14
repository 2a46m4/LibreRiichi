<script setup lang="ts">
import { ErrorStyling } from '../styling'

const props = defineProps<{ error: string }>()
</script>

<template>
  <p :class="ErrorStyling">{{ props.error }}</p>
</template>

<style scoped></style>

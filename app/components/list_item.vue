<script setup lang="ts">
import {ListStyling} from "../styling";
</script>

<template>
<li :class="ListStyling">
  <slot/>
</li>
</template>

<style scoped>

</style>
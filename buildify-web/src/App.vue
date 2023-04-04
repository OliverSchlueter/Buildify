<script lang="ts">
import History from "./components/History.vue";

export default {
    components: {
      History
    }
}
</script>

<template>
  <h1>Build history</h1>
  <History max-amount=6 link="localhost:1337"/>
</template>

<style scoped>
  h1{
    text-align: center;
  }
</style>
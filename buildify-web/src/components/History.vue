<script lang="ts">
export default {
  props: [
    "link"
  ],
  computed:{
    builds(){
      fetch(this.link + "builds").then(res => {
        console.log(res);
      })
    }
  }
}
</script>

<template>
  <h1> link: {{ link }}</h1>
</template>

<style>
  
</style>
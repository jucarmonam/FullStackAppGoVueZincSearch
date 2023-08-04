<script setup lang="ts">
import NavBar from './components/NavBar.vue'
import HomeView from './views/HomeView.vue'
</script>

<template>
  <NavBar />
  <HomeView />
</template>

<style>
* {
  box-sizing: border-box;
}

html {
  font-family: sans-serif;
}

body {
  margin: 0;
}
</style>

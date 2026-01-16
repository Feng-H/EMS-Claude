<template>
  <router-view />
</template>

<script setup lang="ts">
</script>

<style>
/* Global Reset - Fix Scrollbar and Borders */
html {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  margin: 0 !important;
  padding: 0 !important;
  width: 100vw !important;
  overflow-x: hidden !important;
  overflow-y: auto !important;
}

body {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  margin: 0 !important;
  padding: 0 !important;
  width: 100vw !important;
  overflow-x: hidden !important;
  overflow-y: auto !important;
}

#app {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  margin: 0 !important;
  padding: 0 !important;
  width: 100% !important;
  display: block !important;
  overflow-x: hidden !important;
}

/* Hide scrollbar completely - remove the black bar */
html::-webkit-scrollbar,
body::-webkit-scrollbar,
#app::-webkit-scrollbar {
  display: none !important;
  width: 0 !important;
  height: 0 !important;
}

/* Firefox scrollbar */
html, body, #app {
  scrollbar-width: none !important;
  -ms-overflow-style: none !important;
}

/* Remove ALL borders and outlines */
* {
  border: none !important;
  outline: none !important;
}

*:focus,
*:focus-visible,
*:focus-within {
  outline: none !important;
  box-shadow: none !important;
}

/* Reset box model */
*,
*::before,
*::after {
  box-sizing: border-box;
}
</style>

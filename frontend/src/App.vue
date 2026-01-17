<template>
  <router-view v-slot="{ Component, route }">
    <component
      :is="route.meta?.layout === 'full-screen' ? 'div' : MainLayout"
      v-if="Component"
    >
      <component :is="Component" />
    </component>
  </router-view>
</template>

<script setup lang="ts">
import MainLayout from '@/views/layout/MainLayout.vue'
</script>

<style>
/* Global Reset */
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

#app {
  width: 100%;
  min-height: 100vh;
}

/* Hide scrollbar */
html::-webkit-scrollbar,
body::-webkit-scrollbar {
  display: none;
}

html, body {
  scrollbar-width: none;
  -ms-overflow-style: none;
}

/* Reset box model */
*,
*::before,
*::after {
  box-sizing: border-box;
}
</style>

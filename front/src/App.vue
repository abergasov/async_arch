<template>
  <h1 v-if="auth">1231</h1>
  <HelloWorld msg="Hello Vue 3 + Vite" />
</template>

<script setup>
import HelloWorld from './components/HelloWorld.vue'
import { computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()
const auth = computed(() => +store.state.auth === 1)


</script>

<style>

</style>

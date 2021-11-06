<template>
  <van-skeleton v-if="!initDone" :row="20" />
  <div v-else>
    <Auth v-if="!auth"/>
    <HelloWorld msg="Hello Vue 3 + Vite" />
  </div>
</template>

<script setup>
import HelloWorld from './components/HelloWorld.vue'
import Auth from './components/Auth.vue'
import { getCurrentInstance } from 'vue'
const app = getCurrentInstance()
import { computed } from 'vue'
import { useStore } from 'vuex'

const askBackend = app.appContext.config.globalProperties.askBackend;



const store = useStore()
const auth = computed(() => +store.state.auth === 1)
const initDone = computed(() => store.state.init_done)

const getJWT = () => {
  const urlParams = new URLSearchParams(window.location.search);
  let code = urlParams.get('code');
  if (code.length > 0) {
    askBackend("exchange", {code: code}).then(
        data => {
          if (!data.ok) {
            return;
          }
          console.log(data);
          store.commit("initDone");
          store.commit("setJWT", data.code);
        },
        err => console.error(err),
    )
  }
}
getJWT();

</script>

<style>

</style>

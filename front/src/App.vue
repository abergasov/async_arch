<template>
  <van-skeleton v-if="!initDone" :row="20" />
  <div v-else>
    <Auth v-if="!auth"/>
    <Dashboard v-else :user="user" v-on:change_role="changeRole"/>
  </div>
</template>

<script setup>
import Dashboard from './components/Dashboard.vue'
import Auth from './components/Auth.vue'
import { getCurrentInstance } from 'vue'
const app = getCurrentInstance()
import { computed, ref } from 'vue'
import { useStore } from 'vuex'

const askBackend = app.appContext.config.globalProperties.askBackend;

const store = useStore()
const auth = computed(() => +store.state.auth === 1)
const initDone = computed(() => store.state.init_done)
const user = ref({})
// const props = defineProps({
//   user = {
//     email: "",
//     name: "",
//     role: "",
//     public_id: "",
//   }
// })
// let ;

const changeRole = (payload) => {
  askBackend("change_role", payload).then(
      data => {
        if (!data.ok) {
          return
        }
        user.value = data.data;
        store.commit("setJWT", data.token);
      },
      err => console.error(err),
  )
}

const getJWT = () => {
  const urlParams = new URLSearchParams(window.location.search);
  let code = urlParams.get('code');
  if (code) {
    // code exist in GET, try exchange to jwt
    askBackend("exchange", {code: code}).then(
        data => {
          if (!data.ok) {
            return;
          }
          store.commit("initDone");
          store.commit("setJWT", data.code);
          window.location.href = "/";
        },
        err => {
          console.error(err);
          window.location.href = "/";
        },
    )
  } else {
    let jwt = store.state.key;
    if (!jwt) {
      //code not exist
      store.commit("initDone");
      return
    }
    // check jwt valid
    askBackend("dashboard", {}).then(
      data => {
        if (!data.ok) {
          return
        }
        user.value = data.data;
        store.commit("initDone");
        store.commit("setAuth", 1);
        store.commit("setJWT", jwt);

      },
      err => {
        store.commit("initDone");
      },
    );
  }
}
getJWT();

</script>

<style>

</style>

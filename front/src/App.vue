<template>
  <el-skeleton v-if="!initDone" :row="20" />
  <div v-else>
    <Auth v-if="!auth"/>
    <Dashboard v-else :user="user" :tasks="tasks"
               v-on:assign_tasks="assignTasks"
               v-on:change_role="changeRole"
               v-on:done_tasks="finishTask"
               v-on:create_task="createTask"/>
  </div>
</template>

<script setup>
import Dashboard from './components/Dashboard.vue'
import Auth from './components/Auth.vue'
import { getCurrentInstance, h } from 'vue'
import { computed, ref } from 'vue'
import { useStore } from 'vuex'
import { ElNotification } from 'element-plus'

const app = getCurrentInstance()
const askBackend = app.appContext.config.globalProperties.askBackend;
const store = useStore()
const auth = computed(() => +store.state.auth === 1)
const initDone = computed(() => store.state.init_done)
const user = ref({})
const tasks = ref([])

const changeRole = (payload) => {
  askBackend("auth/change_role", payload).then(
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

const assignTasks = () => {
  askBackend("task/assign", {}).then(
      data => {
        if (!data.ok) {
          return;
        }
        tasks.value = data.data;
      },
      err => console.error(err),
  )
}

const createTask = (payload) => {
  askBackend("task/create", payload).then(
      data => {
        if (!data.ok) {
          return;
        }
        ElNotification({
          title: 'Task created',
          message: h('p', { style: 'color: green' }, 'Task created:' + data.data.publicID),
        });
        loadTasks();
      },
      err => console.error(err),
  );
}

const loadTasks = () => {
  askBackend("task/list", {}).then(
      data => {
        if (!data.ok) {
          return;
        }
        tasks.value = data.data;
      },
      err => console.error(err),
  )
}

const finishTask = (payload) => {
  askBackend("task/finish", payload).then(
      data => {
        ElNotification({
          title: 'Task completed',
          message: h('p', { style: 'color: green' }, 'Task completed:' + payload.task_id),
        });
        loadTasks();
        },
      err => console.error(err),
  )
}


const getJWT = () => {
  const urlParams = new URLSearchParams(window.location.search);
  let code = urlParams.get('code');
  if (code) {
    // code exist in GET, try exchange to jwt
    askBackend("auth/exchange", {code: code}).then(
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
    askBackend("auth/dashboard", {}).then(
      data => {
        if (!data.ok) {
          return
        }
        user.value = data.data;
        store.commit("initDone");
        store.commit("setAuth", 1);
        store.commit("setJWT", jwt);
        store.commit("setUser", data.user);
        loadTasks();
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

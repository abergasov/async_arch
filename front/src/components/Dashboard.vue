<template>
  <h1>Das is dashboard</h1>
  <h4>ID: {{ user.public_id }}</h4>
  <h4>Email: {{ user.email }}</h4>
  <h4>User name: {{ user.name }}</h4>
  <h4>Role: {{ user.role }}</h4>
  <p>
    <button @click="clickIT">awdawdw</button>
  </p>
</template>

<script>
export default {
  name: "Dashboard",
  props: {
    user: Object,
  },
  data() {
    return {
      msg: "awdwad",
      content: "",
    };
  },
  methods: {
    clickIT() {
      this.$store.commit('setAuth', 1);
      this.askBackend("auth/awdawd", {})
    },
  }
}

</script>

<style scoped>
a {
  color: #42b983;
}
</style>

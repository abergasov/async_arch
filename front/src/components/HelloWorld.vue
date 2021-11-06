<template>
  <h1>{{ msg }}</h1>
  <p>
    <button @click="clickIT">awdawdw</button>
  </p>
</template>

<script>
export default {
  name: "HelloWorld",
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

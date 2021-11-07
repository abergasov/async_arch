<template>
  <h1>Das is dashboard</h1>
  <h4>ID: {{ user.public_id }}</h4>
  <h4>Email: {{ user.user_mail }}</h4>
  <h4>User name: {{ user.user_name }}</h4>
  <h4>Role: {{ user.user_role }}</h4>
  <hr>
  <p>
    Change role to
    <select v-model="role">
      <option value="admin">Admin</option>
      <option value="worker">Worker</option>
    </select>
    <button @click="changeRole">change role</button>
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
      role: this.user.user_role,
    };
  },
  methods: {
    changeRole() {
      this.$emit("change_role", {role: this.role});
    },
  }
}

</script>

<style scoped>
a {
  color: #42b983;
}
</style>

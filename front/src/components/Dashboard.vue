<template>
  <el-container>
    <el-main>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="box-card">
            <template #header>
              <div class="card-header">
                <h1>User profile</h1>
              </div>
            </template>
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
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card">
            <template #header>
              <div class="card-header">
                <h1>User tasks</h1>
              </div>
            </template>
            <hr>
            <el-form>
              <el-form-item label="task title">
                <el-input v-model="newTask.title"/>
              </el-form-item>
              <el-form-item label="task description">
                <el-input v-model="newTask.desc"/>
              </el-form-item>
            </el-form>
            <el-button :disabled="newTask.desc.length === 0 || newTask.title.length === 0" type="primary" @click="createTask">create task</el-button>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </el-container>

</template>

<script>
export default {
  name: "Dashboard",
  props: {
    user: Object,
  },
  data() {
    return {
      newTask: {
        title: "",
        desc: ""
      },
      role: this.user.user_role,
    };
  },
  methods: {
    changeRole() {
      this.$emit("change_role", {role: this.role});
    },
    createTask() {
      this.$emit("create_task", {task: this.newTask});
      this.newTask = {
        title: "",
        desc: ""
      };
    }
  }
}

</script>

<style scoped>
a {
  color: #42b983;
}
</style>

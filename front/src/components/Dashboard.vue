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
            <h4>ID: {{ user.publicID }}</h4>
            <h4>Email: {{ user.userMail }}</h4>
            <h4>User name: {{ user.userName }}</h4>
            <h4>Role: {{ user.userRole }}</h4>
            <hr>
            <p>
              Change role to
              <select v-model="role">
                <option value="admin">Admin</option>
                <option value="worker">Worker</option>
                <option value="manager">Manager</option>
              </select>
              <button @click="changeRole">change role</button>
              <br/>
              <br/>
              <button v-if="user.userRole !== 'worker'" @click="assignTasks">assign tasks</button>
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
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card class="box-card">
            <el-table :data="tasks" style="width: 100%"
                      :row-class-name="tableRowClassName">
              <el-table-column prop="publicID" label="public_id" width="280" />
              <el-table-column prop="title" label="title" width="180" />
              <el-table-column prop="jiraID" label="jiraID" width="180" />
              <el-table-column prop="createdAt" label="created_at" />
              <el-table-column prop="publicStatus" label="status" />
              <el-table-column label="Operations">
                <template #default="scope">
                  <el-button v-if="showBtn(scope.row)" size="mini" @click="handleEdit(scope.row)"
                  >Finish</el-button>
                </template>
              </el-table-column>
            </el-table>
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
    tasks: Array,
  },
  data() {
    return {
      newTask: {
        title: "",
        desc: ""
      },
      role: this.user.userRole,
    };
  },
  methods: {
    handleEdit(row) {
      this.$emit("done_tasks", {task_id: row.publicID});
    },
    tableRowClassName({ row, rowIndex }) {
      if (row.assigned_to === this.$store.state.user) {
        return 'warning-row'
      }
      return ''
    },
    showBtn(row) {
      if (row.assignedTo === this.$store.state.user && row.status !== "done") {
        return true
      }
      return false
    },

    changeRole() {
      this.$emit("change_role", {role: this.role});
    },
    assignTasks() {
      this.$emit("assign_tasks");
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
.warning-row {
  background-color: greenyellow;
}
</style>

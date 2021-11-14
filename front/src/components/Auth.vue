<template>
  <el-button type="primary" @click="googleLogin">Login with google</el-button>
</template>

<script setup>
  function googleLogin() {
    window.location.href = "http://localhost:8090/auth/google/login";
  }
</script>
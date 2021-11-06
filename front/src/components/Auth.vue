<template>
  <van-button @click="googleLogin" type="primary">Login with google</van-button>
</template>

<script setup>
  function googleLogin() {
    window.location.href = "http://localhost:8090/auth/google/login";
  }
</script>
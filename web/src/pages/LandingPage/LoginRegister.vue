<template>
  <div class="loginContainer">
    <h4 class="authTitle text-center text-weight-bold q-mt-none q-mb-sm">Login</h4>
    <p class="authTitle text-center q-mt-none q-mb-lg text-grey">
      Enter your credentials to continue
    </p>
    <BaseInput v-model="username" placeholder="Username" class="q-mb-md" />
    <PasswordInput v-model="password" placeholder="Password" class="q-mb-xl" />
    <q-btn label="Login" color="secondary" style="width: 100%" @click="notImplemented" />
    <p class="q-mt-md q-mb-none">
      Don't have an account?
      <span
        class="text-secondary cursor-pointer"
        style="text-decoration: underline; font-weight: bold"
        @click="notImplemented"
        >Create one</span
      >
    </p>
  </div>
</template>

<script setup lang="ts">
import PasswordInput from 'src/components/input/PasswordInput.vue';
import BaseInput from 'src/components/input/BaseInput.vue';
import { useNotify } from 'src/composables/useNotify';
import { ref } from 'vue';

const { notImplemented } = useNotify();
const username = ref('');
const password = ref('');
</script>

<style lang="scss" scoped>
.loginContainer {
  width: clamp(15rem, 25vw, 25rem);
  border-radius: 5px;
  box-shadow: 0px 0px 8px 1px grey;
  padding: 2rem;
  text-align: center;
}
</style>

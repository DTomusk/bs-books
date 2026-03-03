<template>
  <div class="authContainer">
    <h4 class="authTitle text-center text-weight-bold q-mt-none q-mb-sm">
      {{ isCreatingAcc ? 'Create Account' : 'Login' }}
    </h4>
    <p class="authTitle text-center q-mt-none q-mb-lg text-grey">
      Enter your credentials to continue
    </p>
    <BaseInput v-model="email" placeholder="Email" class="q-mb-md" />
    <BaseInput v-if="isCreatingAcc" v-model="username" placeholder="Username" class="q-mb-md" />
    <PasswordInput v-model="password" placeholder="Password" class="q-mb-xl" />
    <q-btn
      :label="isCreatingAcc ? 'Register' : 'Login'"
      color="secondary"
      style="width: 100%"
      @click="notImplemented"
    />
    <p class="q-mt-md q-mb-none">
      {{ isCreatingAcc ? 'Already have an account?' : "Don't have an account?" }}
      <span
        class="text-secondary cursor-pointer"
        style="text-decoration: underline; font-weight: bold"
        @click="isCreatingAcc = !isCreatingAcc"
        >{{ isCreatingAcc ? 'Sign in' : 'Create one' }}</span
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
const isCreatingAcc = ref(false); // When false, user is logging in
const username = ref('');
const email = ref('');
const password = ref('');
</script>

<style lang="scss" scoped>
.authContainer {
  width: 20rem;
  border-radius: 5px;
  box-shadow: 0px 0px 8px 1px grey;
  padding: 2rem;
  text-align: center;
}
</style>

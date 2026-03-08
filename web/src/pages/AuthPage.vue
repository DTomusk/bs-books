<template>
  <div class="authContainer">
    <h4 class="authTitle text-center text-weight-bold q-mt-none q-mb-sm">
      {{ mode === 'login' ? 'Login' : 'Create Account' }}
    </h4>
    <p class="authTitle text-center q-mt-none q-mb-lg text-grey">
      Enter your credentials to continue
    </p>
    <BaseInput v-model="email" placeholder="Email" />
    <BaseInput v-if="mode === 'register'" v-model="username" placeholder="Username" />
    <PasswordInput v-model="password" placeholder="Password" class="q-mb-xl">
      <template #bottom>
        <PasswordStrength v-if="mode === 'register'" :password="password" />
      </template>
    </PasswordInput>

    <q-btn
      :label="mode === 'login' ? 'Login' : 'Create Account'"
      color="secondary"
      style="width: 100%"
      @click="notImplemented"
    />
    <p class="q-mt-md q-mb-none">
      {{ mode === 'login' ? "Don't have an account?" : 'Already have an account?' }}
      <span
        class="text-secondary cursor-pointer"
        style="text-decoration: underline; font-weight: bold"
        @click="handleToggleMode"
        >{{ mode === 'login' ? 'Create one' : 'Sign in' }}</span
      >
    </p>
  </div>
</template>

<script setup lang="ts">
import PasswordStrength from 'src/components/rating/PasswordStrength.vue';
import PasswordInput from 'src/components/input/PasswordInput.vue';
import BaseInput from 'src/components/input/BaseInput.vue';
import { useNotify } from 'src/composables/useNotify';
import { useRouter } from 'vue-router';
import { ref } from 'vue';

const { mode = 'login', isOwnPage = true } = defineProps<{
  mode?: 'login' | 'register';
  isOwnPage?: boolean; // When true, toggling between login/reg changes page
}>();
const emit = defineEmits<{ (e: 'toggleMode'): void }>();

const { error, notImplemented } = useNotify();
const router = useRouter();

const email = ref('');
const username = ref('');
const password = ref('');

function handleToggleMode() {
  email.value = '';
  username.value = '';
  password.value = '';
  if (isOwnPage) {
    const route = mode === 'login' ? '/auth/create-account' : '/auth/login';
    router.push(route).catch(() => {
      error(`Route not found: ${route}`);
    });
  } else {
    emit('toggleMode');
  }
}
</script>

<style lang="scss" scoped>
.authContainer {
  margin: auto;
  width: 20rem;
  border-radius: 5px;
  background-color: white;
  box-shadow: 0px 0px 8px 1px grey;
  padding: 2rem;
  margin-block: 2rem 4rem;
  text-align: center;
}
</style>

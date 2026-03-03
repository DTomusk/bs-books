<template>
  <q-layout view="hhh lpR fff" class="bg-primary">
    <Header :variant="variant" @toggleDrawer="toggleDrawer" />

    <q-drawer v-model="rightDrawerOpen" behavior="mobile" side="right" elevated>
      <q-list class="full-height" style="display: flex; flex-direction: column">
        <EssentialLink
          v-for="link in [...baseLinks, ...(isLoggedIn ? loggedInLinks : loggedOutLinks)]"
          :key="link.title"
          v-bind="link"
          @closeDrawer="rightDrawerOpen = false"
        />
        <q-space />
        <div class="flex justify-between q-px-sm q-py-xs bg-grey-10 text-white items-center">
          <p :style="{ margin: '0px', fontSize: 'medium' }">ⓒ BS Books - 2026</p>
          <p :style="{ margin: '0px' }">v{{ version }}</p>
        </div>
      </q-list>
    </q-drawer>

    <Footer />

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import EssentialLink, { type EssentialLinkProps } from 'components/EssentialLink.vue';
import { version } from '../../package.json';
import Footer from 'src/layouts/Footer.vue';
import Header from 'src/layouts/Header.vue';
import { useRoute } from 'vue-router';
import { ref, computed } from 'vue';

const route = useRoute();
const isLoggedIn = ref(false);
const rightDrawerOpen = ref(false);
const variant = computed(() => route.meta.variant);

const baseLinks: EssentialLinkProps[] = [
  {
    title: 'Home',
    icon: 'home',
    route: '/',
  },
];

const loggedInLinks: EssentialLinkProps[] = [
  {
    title: 'Some Link',
    icon: 'fa-regular fa-circle-question',
  },
  {
    title: 'Another Link',
    icon: 'fa-solid fa-person-drowning',
  },
];

const loggedOutLinks: EssentialLinkProps[] = [
  {
    title: 'Login',
    icon: 'mdi-account-circle',
    route: '/auth/login',
  },
  {
    title: 'Create Account',
    icon: 'mdi-account-circle',
    route: '/auth/create-account',
  },
];

function toggleDrawer() {
  rightDrawerOpen.value = !rightDrawerOpen.value;
}
</script>

<style lang="scss" scoped>
.background {
  --multiplier: 0.5;

  background:
    url('src/assets/backgroundPatternHorizontal.png'),
    linear-gradient(180deg, $secondary 0%, $primary 100%);

  background-size: calc(var(--multiplier) * 3656px) calc(var(--multiplier) * 2496px);
}
</style>

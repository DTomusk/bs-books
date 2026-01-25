<template>
  <q-layout view="lHh Lpr lff" class="background">
    <q-header v-if="loggedIn" elevated>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
        <q-toolbar-title> BS Books</q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header> Essential Links </q-item-label>

        <EssentialLink v-for="link in linksList" :key="link.title" v-bind="link" />
      </q-list>
    </q-drawer>

    <q-footer>
      <div class="footerLinks">
        <p v-for="link in ['About', 'FAQs', 'Contact Us', 'Careers']" :key="link" class="hoverable">
          {{ link }}
        </p>
      </div>
      <div class="flex-centre-row" style="margin-block: 1rem 2rem">
        <q-icon
          v-for="social in ['fa-instagram', 'fa-twitter', 'fa-facebook', 'fa-youtube']"
          :key="social"
          :name="'fa-brands ' + social"
          size="x-large"
          class="hoverable"
        />
      </div>
      <p style="color: #999; text-align: center">ⓒ 2026 BS Books - All rights reserved</p>
    </q-footer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import EssentialLink, { type EssentialLinkProps } from 'components/EssentialLink.vue';

const linksList: EssentialLinkProps[] = [
  {
    title: 'Docs',
    caption: 'quasar.dev',
    icon: 'school',
    link: 'https://quasar.dev',
  },
];

const leftDrawerOpen = ref(false);
const loggedIn = false;

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>

<style lang="scss" scoped>
.background {
  --tile: 1200px;

  background:
    url('src/assets/backgroundPatternHorizontal.svg') repeat,
    linear-gradient(180deg, $secondary 0%, $primary 100%);

  background-size:
    var(--tile) var(--tile),
    contain;
}

.footerLinks {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 1rem;
}

.footerLinks > p {
  flex: 1 1 150px;
  max-width: 150px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

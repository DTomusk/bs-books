<template>
  <q-layout view="hhh lpR fff" class="bg-primary">
    <q-header>
      <q-toolbar style="position: absolute; z-index: 3">
        <div class="topHeroContainer">
          <div style="width: 10vw; max-width: 5rem">
            <q-img src="src/assets/Logo.png" />
          </div>
          <q-btn
            v-if="isMobile"
            flat
            dense
            round
            icon="menu"
            aria-label="Menu"
            @click="toggleLeftDrawer"
          />
          <div v-else style="display: flex; align-items: center">
            <div class="linksContainer text-grey">
              <a @click="notImplemented">Home</a>
              <a @click="notImplemented">Explore</a>
              <a @click="notImplemented">About Us</a>
            </div>
            <q-input v-model="search" standout rounded dense bg-color="grey-8">
              <template v-slot:append>
                <q-icon name="search" @click="notImplemented" class="hoverable highlight" />
              </template>
            </q-input>
          </div>
        </div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="rightDrawerOpen" behavior="mobile" side="right" elevated>
      <q-list>
        <q-item-label header> Essential Links </q-item-label>

        <EssentialLink v-for="link in linksList" :key="link.title" v-bind="link" />
      </q-list>
    </q-drawer>

    <q-footer class="bg-dark q-pt-lg">
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
import EssentialLink, { type EssentialLinkProps } from 'components/EssentialLink.vue';
import { useNotify } from 'src/composables/useNotify';
import { computed, ref } from 'vue';
import { useQuasar } from 'quasar';

const linksList: EssentialLinkProps[] = [
  {
    title: 'Docs',
    caption: 'quasar.dev',
    icon: 'school',
    link: 'https://quasar.dev',
  },
];

const $q = useQuasar();
const search = ref('');
const rightDrawerOpen = ref(false);
const { notImplemented } = useNotify();
const isMobile = computed(() => $q.screen.lt.md);

function toggleLeftDrawer() {
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

.topHeroContainer {
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  top: 0;
  right: 0;
  z-index: 4;
  padding: 1rem 0rem;

  @media (min-width: $breakpoint-sm-max) {
    width: 65%; // Same as max-width on heroContainer _LandingPage.vue
    padding: 1rem 2rem 1rem 1rem;
  }
}

.linksContainer {
  display: flex;
  margin-right: 2rem;
  gap: 2rem;
  font-weight: bold;
  font-size: medium;

  a {
    padding: 1rem 0.5rem;
  }
}

.highlight:hover,
a:hover {
  color: white;
}

.footerLinks {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.footerLinks > p {
  flex: 1 1 150px;
  max-width: 150px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

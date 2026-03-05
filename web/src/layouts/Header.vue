<template>
  <q-header>
    <q-toolbar :class="[{ ['absoluteToolbar']: variant === 'main' }]">
      <div class="topHeroContainer">
        <div style="width: 12vw; max-width: 80px">
          <q-img src="src/assets/Logo.png" @click="router.push('/')" />
        </div>
        <div v-if="isMobile" class="row" style="gap: 0.5rem">
          <q-btn
            v-if="variant === 'main'"
            flat
            dense
            round
            icon="mdi-account"
            aria-label="account"
            @click="router.push('/auth/login')"
          />
          <q-btn
            v-else
            flat
            dense
            round
            icon="mdi-home"
            aria-label="home"
            @click="router.push('/')"
          />
          <q-btn flat dense round icon="menu" aria-label="Menu" @click="emit('toggleDrawer')" />
        </div>
        <div v-else style="display: flex; align-items: center">
          <div class="linksContainer text-grey">
            <a
              v-for="link in ['Home', 'Explore', 'About Us']"
              :key="link"
              @click="notImplemented"
              >{{ link }}</a
            >
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
</template>

<script setup lang="ts">
import { useScreen } from 'src/composables/useScreen';
import { useNotify } from 'src/composables/useNotify';
import { useRouter } from 'vue-router';
import { ref } from 'vue';

defineProps<{ variant: 'main' | 'auth' }>();
const emit = defineEmits<{ (e: 'toggleDrawer'): void }>();

const router = useRouter();
const search = ref('');
const { notImplemented } = useNotify();
const { isMobile } = useScreen();
</script>

<style lang="scss" scoped>
.absoluteToolbar {
  position: absolute;
  z-index: 3;
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
    width: 65%; // width & max-width need to match heroContainer _LandingPage.vue
    max-width: 1200px;
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
</style>

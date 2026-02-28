<template>
  <div
    :class="['container bg-offwhite full-height full-width', { glow: glowing }]"
    @mouseenter="!isMobile ? triggerGlow() : null"
    @click="triggerGlow"
  >
    <div class="flex-centre-col">
      <div :key="hoverKey">
        <slot name="icon"></slot>
      </div>
    </div>
    <div style="display: block">
      <div class="row">
        <p class="card-title">{{ title }}</p>
      </div>
      <p v-for="x in text" :key="x" class="text">
        {{ x }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuasar } from 'quasar';

defineProps<{ title: string; text: string[] }>();

const $q = useQuasar();
const hoverKey = ref(0);
const glowing = ref(false);
let isPlaying = false;
const animationDurationMs = 1000; // Approx (may need tweaking)
const isMobile = computed(() => $q.screen.lt.md);

/**
 * Triggers icon animation (& glow [if applicable])
 * Sets timeout to avoid buggy state of animated icons
 */
function triggerGlow() {
  if (isPlaying) {
    return;
  }

  isPlaying = true;
  hoverKey.value += 1;
  setTimeout(() => {
    isPlaying = false;
  }, animationDurationMs);

  if (!isMobile.value) return;
  glowing.value = false;
  requestAnimationFrame(() => {
    glowing.value = true;
  });
}
</script>

<style lang="scss" scoped>
.container {
  display: flex;
  padding: 1rem;
  border-radius: 10px;
  box-shadow: 0px 0px 7px white;
  transition: all 0.2s ease-in-out;
  max-width: 650px;
}

@media only screen and (min-width: $breakpoint-sm-max) {
  .container:hover {
    cursor: pointer;
    box-shadow: 0px 0px 15px white;
  }
}

.glow {
  animation: glow-pulse 1000ms ease-out;
}

@keyframes glow-pulse {
  0% {
    box-shadow: 0px 0px 7px white;
  }
  50% {
    box-shadow: 0px 0px 15px white;
  }
  100% {
    box-shadow: 0px 0px 7px white;
  }
}

.card-title {
  font-size: 110%;
  font-weight: bold;
  line-height: 100%;
  margin-block: 0.5rem;
}

.text {
  font-size: 0.9rem;
  margin-bottom: 0;
}
</style>

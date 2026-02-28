<template>
  <div
    :class="['row justify-between items-center', { 'q-pb-sm': isMobile }, { 'q-pt-xl': !isMobile }]"
  >
    <h4 class="text-offwhite" style="font-weight: bold; text-decoration: underline">{{ title }}</h4>
    <div>
      <q-btn
        style="margin-right: 1rem"
        icon="fa-solid fa-arrow-left"
        round
        :size="isMobile ? 'sm' : 'md'"
        color="secondary"
        @click="scrollByViewport(-1)"
      />
      <q-btn
        icon="fa-solid fa-arrow-right"
        round
        :size="isMobile ? 'sm' : 'md'"
        color="secondary"
        @click="scrollByViewport(1)"
      />
    </div>
  </div>
  <div class="carousel-wrapper">
    <div ref="track" class="carousel">
      <div v-for="book in books" :key="book.title" class="carousel-item">
        <div class="featured-item" :key="book.title">
          <q-img class="featured-img" :src="book.image" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Book } from 'src/components/models';
import { ref, computed } from 'vue';
import { useQuasar } from 'quasar';

defineProps<{
  title: string;
  books: Book[];
}>();

const $q = useQuasar();
const isMobile = computed(() => $q.screen.lt.md);
const track = ref<HTMLElement | null>(null);

function scrollByViewport(direction: 1 | -1) {
  const el = track.value;
  if (!el) return;

  const amount = el.clientWidth * 0.8;
  el.scrollBy({
    left: amount * direction,
    behavior: 'smooth',
  });
}
</script>

<style lang="scss" scoped>
.featured-item {
  width: 125px;
  padding-right: 12px;
  padding-bottom: 12px;

  @media (min-width: $breakpoint-sm-max) {
    width: 150px;
  }
}

.featured-img {
  border-radius: 6px;
  height: 100%;
}

.carousel {
  display: flex;
  gap: 1rem;
  overflow-x: auto;
  scroll-behavior: smooth;
  scroll-snap-type: x proximity;
}

.carousel-item {
  flex: 0 0 auto;
  scroll-snap-align: start;
}

@media (min-width: $breakpoint-sm-max) {
  // Firefox
  .carousel {
    scrollbar-width: thin;
    scrollbar-color: #bbb transparent;
  }

  // Chrome / Edge / Safari
  .carousel::-webkit-scrollbar {
    height: 8px;
  }

  .carousel::-webkit-scrollbar-track {
    background: transparent;
  }
}
</style>

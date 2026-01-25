<template>
  <section class="q-pa-md">
    <p class="title text-white q-mb-sm">{{ title }}</p>
    <q-virtual-scroll
      class="featured-scroll"
      :items="books"
      virtual-scroll-horizontal
      v-slot="{ item, index }"
    >
      <div class="featured-item" :key="index">
        <q-img class="featured-img" :src="item.image" />
      </div>
    </q-virtual-scroll>
  </section>
</template>

<script setup lang="ts">
import type { Book } from 'src/components/models';

defineProps<{
  title: string;
  books: Book[];
}>();
</script>

<style lang="scss" scoped>
.featured-scroll {
  min-height: 150px;
}

.featured-item {
  width: 100px;
  padding-right: 12px;
  padding-bottom: 12px;
}

.featured-img {
  border-radius: 6px;
  height: 100%;
}
</style>

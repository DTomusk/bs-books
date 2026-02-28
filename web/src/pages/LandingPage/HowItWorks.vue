<template>
  <div class="row q-col-gutter-lg q-mt-sm q-mb-xl">
    <div class="col-12 flex justify-center">
      <div class="full-width" :style="{ maxWidth: isMobile ? '650px' : '1400px' }">
        <h4
          ref="howItWorksEl"
          class="text-offwhite q-mb-sm"
          style="font-weight: bold; text-decoration: underline"
        >
          How It Works
        </h4>
      </div>
    </div>
    <div class="col-12 col-md-4 flex justify-center">
      <ActionCard
        title="Rate"
        :text="['Rate books using hearts and poos - Hearts are good and poos are bad']"
        ><template #icon>
          <lord-icon
            trigger="in"
            stroke="bold"
            target=".container"
            src="src/assets/heart-coloured.json"
          ></lord-icon> </template
      ></ActionCard>
    </div>
    <div class="col-12 col-md-4 flex justify-center">
      <ActionCard
        title="Collect"
        :text="[
          'Create your own digital library of books which are brilliant, bad, or strangely both',
        ]"
        ><template #icon>
          <lord-icon
            trigger="in"
            target=".container"
            src="src/assets/bookshelf-coloured.json"
          ></lord-icon> </template
      ></ActionCard>
    </div>
    <div class="col-12 col-md-4 flex justify-center">
      <ActionCard
        title="Contribute"
        :text="[
          'Contribute to our database of books by adding covers, descriptions, authors, etc...',
        ]"
        ><template #icon>
          <lord-icon
            trigger="in"
            target=".container"
            src="src/assets/database-coloured.json"
          ></lord-icon> </template
      ></ActionCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import ActionCard from 'src/pages/LandingPage/ActionCard.vue';
import { defineElement } from '@lordicon/element';
import { computed, ref } from 'vue';
import { useQuasar } from 'quasar';

defineElement();
const $q = useQuasar();
const isMobile = computed(() => $q.screen.lt.md);
const howItWorksEl = ref<HTMLElement | null>(null);

defineExpose({
  howItWorksEl,
});
</script>

<style lang="scss" scoped>
lord-icon {
  --size: clamp(2rem, 10vw, 4.5rem);
  width: var(--size);
  height: var(--size);
  margin-right: 1rem;
}

.howItWorks {
  display: flex;
  justify-content: center;
  margin-block: 0.5rem;

  @media (min-width: $breakpoint-sm-max) {
    justify-content: start;
    margin-block: 0.5rem -1.5rem;
  }
}

.cardContainer {
  display: grid;
  gap: 1rem;
  grid-template-columns: 1fr;

  @media (min-width: $breakpoint-sm-max) {
    grid-template-columns: repeat(3, minmax(0, 1fr));
    justify-items: center;
  }
}
</style>

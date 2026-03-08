<template>
  <div class="container" :style="{ fontSize: iconSize }">
    <span v-for="(fill, i) in fills" :key="i" class="icon-wrapper">
      <!-- Using poop svg because the font-awesome outline poop icon wasn't free (boo!) -->
      <PoopIcon
        v-if="icon === 'poop'"
        :filled="false"
        :colour="getPaletteColor(colour) ?? 'black'"
      />
      <q-icon v-else :name="iconName.outline" class="icon" :color="colour" />
      <span class="icon-fill" :style="{ width: fill }">
        <PoopIcon
          v-if="icon === 'poop'"
          :filled="true"
          :colour="getPaletteColor(colour) ?? 'black'"
        />
        <q-icon v-else :name="iconName.fill" class="icon" :color="colour" />
      </span>
    </span>
  </div>
</template>

<script setup lang="ts">
import PoopIcon from 'src/components/rating/PoopIcon.vue';
import { colors } from 'quasar';
import { computed } from 'vue';

const { getPaletteColor } = colors;

const sizes = {
  sm: '14px',
  md: '20px',
  lg: '28px',
} as const;

const icons = {
  star: { outline: 'fa-regular fa-star', fill: 'fa-solid fa-star' },
  poop: { outline: 'mdi-emoticon-poop-outline', fill: 'mdi-emoticon-poop' }, // Fallback
  heart: { outline: 'fa-regular fa-heart', fill: 'fa-solid fa-heart' },
} as const;

const { size = 'md', ...props } = defineProps<{
  rating: number;
  colour: 'primary' | 'secondary' | 'accent' | 'heart' | 'poop';
  icon: keyof typeof icons;
  size?: keyof typeof sizes;
}>();

const computedRating = computed(() => Math.min(Math.max(props.rating, 0), 5));
const iconSize = computed(() => sizes[size]);
const iconName = computed(() => icons[props.icon]);
const fills = computed(() => {
  return Array.from({ length: 5 }, (_, i) => {
    const fill = Math.min(Math.max(computedRating.value - i, 0), 1);
    return `${fill * 100}%`;
  });
});
</script>

<style lang="scss" scoped>
.container {
  display: flex;
  transform-origin: left center;
  align-items: center;
}

.icon-wrapper {
  position: relative;
  width: 1em;
  height: 1em;
}
.icon-fill {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.icon {
  position: absolute;
  inset: 0;
  width: 1em;
  height: 1em;
  line-height: 1;
}
</style>

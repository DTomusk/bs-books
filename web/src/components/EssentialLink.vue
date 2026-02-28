<template>
  <q-item clickable @click="linkClicked" class="essentialLink">
    <q-item-section v-if="icon" avatar style="min-width: 0">
      <q-icon :name="icon" />
    </q-item-section>

    <q-item-section>
      <q-item-label>{{ title }}</q-item-label>
    </q-item-section>
  </q-item>
</template>

<script setup lang="ts">
import { useNotify } from 'src/composables/useNotify';
import { useRouter } from 'vue-router';

export interface EssentialLinkProps {
  title: string;
  icon: string;
  route?: string;
}

const props = defineProps<EssentialLinkProps>();
const emit = defineEmits<{ (e: 'closeDrawer'): void }>();
const router = useRouter();
const { error, notImplemented } = useNotify();

function linkClicked() {
  if (props.route) {
    emit('closeDrawer');
    router.push(props.route).catch(() => {
      error(`Route not found: ${props.route}`);
    });
  } else {
    notImplemented();
  }
}
</script>

<style lang="scss" scoped>
.essentialLink {
  border-bottom: 1px solid $grey-5;
}
</style>

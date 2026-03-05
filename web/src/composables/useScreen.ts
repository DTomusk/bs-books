import { computed } from 'vue';
import { useQuasar } from 'quasar';

export function useScreen() {
  const $q = useQuasar();

  const isMobile = computed(() => $q.screen.lt.md);

  return {
    isMobile,
  };
}

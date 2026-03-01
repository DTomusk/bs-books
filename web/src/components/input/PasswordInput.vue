<template>
  <BaseInput
    v-model="modelValue as string"
    :inputType="isPwd ? 'password' : 'text'"
    :placeholder="placeholder"
  >
    <template #append>
      <q-icon
        :name="isPwd ? 'visibility_off' : 'visibility'"
        class="cursor-pointer"
        @click="isPwd = !isPwd"
      />
    </template>
  </BaseInput>
</template>

<script setup lang="ts">
import BaseInput from 'src/components/input/BaseInput.vue';
import { ref } from 'vue';

defineProps<{
  placeholder: string;
}>();

const modelValue = defineModel<string>();
const isPwd = ref(true);
</script>

<style lang="scss" scoped></style>

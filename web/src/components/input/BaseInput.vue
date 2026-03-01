<template>
  <q-input
    v-model="modelValue"
    outlined
    color="secondary"
    :type="inputType"
    :placeholder="placeholder"
  >
    <template v-slot:append>
      <slot name="append" />
    </template>
  </q-input>
</template>

<script setup lang="ts">
const modelValue = defineModel<string>();

const { inputType = 'text' } = defineProps<{
  placeholder: string;
  inputType?: 'text' | 'password';
}>();
</script>

<style lang="scss" scoped></style>

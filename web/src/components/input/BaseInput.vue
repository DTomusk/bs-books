<template>
  <q-input
    v-model="modelValue"
    outlined
    bottom-slots
    color="secondary"
    :type="inputType"
    :placeholder="placeholder"
  >
    <template v-slot:append>
      <slot name="append" />
    </template>
    <template v-slot:hint>
      <slot name="bottom" />
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

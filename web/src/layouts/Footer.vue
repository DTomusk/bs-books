<template>
  <q-footer class="bg-dark q-pt-lg">
    <div class="footerLinks">
      <p
        v-for="link in ['About', 'FAQs', 'Contact Us', 'Careers']"
        :key="link"
        @click="notImplemented"
        class="hoverable"
      >
        {{ link }}
      </p>
    </div>
    <div class="flex-centre-row" style="margin-block: 1rem 2rem">
      <q-icon
        v-for="social in ['fa-instagram', 'fa-twitter', 'fa-facebook', 'fa-youtube']"
        :key="social"
        @click="notImplemented"
        :name="'fa-brands ' + social"
        size="x-large"
        class="hoverable"
      />
    </div>
    <p style="color: #999; text-align: center">ⓒ 2026 BS Books - All rights reserved</p>
  </q-footer>
</template>

<script setup lang="ts">
import { useNotify } from 'src/composables/useNotify';

const { notImplemented } = useNotify();
</script>

<style lang="scss" scoped>
.footerLinks {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.footerLinks > p {
  flex: 1 1 150px;
  max-width: 150px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

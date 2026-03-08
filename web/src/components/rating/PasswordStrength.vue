<template>
  <div class="pwdStrengthContainer">
    <BaseRating :rating="poopScore" colour="poop" icon="poop" size="sm" />
    <BaseRating :rating="heartScore" colour="heart" icon="heart" size="sm" />
  </div>
</template>

<script setup lang="ts">
import BaseRating from 'src/components/rating/BaseRating.vue';
import { computed } from 'vue';

const { password } = defineProps<{ password: string }>();

const sequences = ['abcd', '123', 'qwerty', 'asdf'];
const common = ['bsbook', 'password', 'welcome', 'admin', 'abc123'];
const pwdAttrs = computed(() => {
  const lowerPwd = password.toLowerCase();
  return {
    length: password.length,
    hasLower: /[a-z]/.test(password),
    hasUpper: /[A-Z]/.test(password),
    hasDigit: /\d/.test(password),
    hasSpecial: /[^A-Za-z0-9]/.test(password),
    narrowSelection: new Set(password).size < password.length / 2,
    hasCommon: common.some((p) => lowerPwd.includes(p)),
    hasSequences: sequences.some((p) => lowerPwd.includes(p)),
  };
});

/**
 * Calculates an arbitrary score (0 ≤ x ≤ 5, to 1dp) for the strength (hearts) of the given password. A higher score
 * indicates a password has attributes of a strong password, but is not necessarily a strong password inherently.
 */
const heartScore = computed(() => {
  let score = 0;

  // Length (Max 2 bonus hearts)
  score += Math.min(2, length * 0.2);

  // Character variation
  if (pwdAttrs.value.hasLower) score += 0.5;
  if (pwdAttrs.value.hasUpper) score += 0.5;
  if (pwdAttrs.value.hasDigit) score += 0.5;
  if (pwdAttrs.value.hasSpecial) score += 1;

  // Diversity bonus (Here, at BS Books, we believe in diversity)
  const types =
    (pwdAttrs.value.hasLower ? 1 : 0) +
    (pwdAttrs.value.hasUpper ? 1 : 0) +
    (pwdAttrs.value.hasDigit ? 1 : 0) +
    (pwdAttrs.value.hasSpecial ? 1 : 0);
  if (types >= 2) score += 0.5;
  if (types === 4) score += 1;

  // Long password bonus
  if (pwdAttrs.value.length >= 8) score += 0.5;
  if (pwdAttrs.value.length >= 12) score += 0.5;

  score = Math.min(5, score);

  // Passwords with any of the following flaws cannot reach 5 hearts
  if (
    pwdAttrs.value.hasCommon ||
    pwdAttrs.value.hasSequences ||
    pwdAttrs.value.narrowSelection ||
    /(.)\1{3,}/.test(password)
  )
    score -= 1;

  return Math.round(score * 10) / 10;
});

/**
 * Calculates an arbitrary score (0 ≤ x ≤ 5, to 1dp) for the weakness (poops) of the given password. A lower score
 * indicates a password has attributes of a weak password, but is not necessarily a weak password inherently.
 */
const poopScore = computed(() => {
  let score = 0;

  // Character variation
  if (!pwdAttrs.value.hasLower) score += 1;
  if (!pwdAttrs.value.hasUpper) score += 1;
  if (!pwdAttrs.value.hasDigit) score += 1;
  if (!pwdAttrs.value.hasSpecial) score += 1;

  // Consecutive repeated characters
  if (/(.)\1{2,}/.test(password)) score += 2;
  else if (/(.)\1{3,}/.test(password)) score += 3;
  else if (/(.)\1{4,}/.test(password)) score += 4;

  if (pwdAttrs.value.hasCommon) score += 2.5;
  if (pwdAttrs.value.hasSequences) score += 2.5;
  if (pwdAttrs.value.narrowSelection) score += 3;

  // Length (Up to 2 poops for short passwords)
  score += Math.max(2 - pwdAttrs.value.length * 0.25, 0);

  score = Math.min(5, score);

  const poopyBusiness = ['poop', 'poo', 'shit', 'shite', 'diarrhea', 'diarrhoea', 'diarrea'];
  if (poopyBusiness.some((p) => password.toLowerCase().includes(p))) score += 5;

  return Math.round(score * 10) / 10;
});
</script>

<style lang="scss">
.pwdStrengthContainer {
  display: flex;
  justify-self: start;
}

.q-field__bottom {
  padding: 0.25rem 0rem !important;
}
</style>

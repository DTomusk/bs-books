<template>
  <div class="heroImg">
    <q-img src="src/assets/bookshelf.jpg" fit="cover" width="100%" height="100%">
      <div class="top-bottom-gradient">
        <div class="imgContent">
          <h1 class="tagline">The book rating system you never knew you needed</h1>
          <h4 class="subTagline">
            Rate books using <span><q-icon name="fa-solid fa-heart" /></span> and
            <span><q-icon name="fa-solid fa-poop" /></span> because some books deserve more than
            just love
          </h4>
          <q-btn
            label="Learn More"
            color="positive"
            text-color="black"
            no-caps
            :size="isMobile ? 'md' : 'lg'"
            @click="emit('learnMore')"
          />
        </div>
      </div>
    </q-img>
  </div>
</template>

<script setup lang="ts">
import { useScreen } from 'src/composables/useScreen';

const emit = defineEmits(['learnMore']);

const { isMobile } = useScreen();
</script>

<style lang="scss">
.top-bottom-gradient {
  background: linear-gradient(
    180deg,
    rgba(0, 0, 0, 0.9) 0%,
    rgba(0, 0, 0, 0.7) 50%,
    rgba(0, 0, 0, 0.9) 100%
  ) !important;
  height: 100%;
  width: 100%;
  position: absolute;
  bottom: 0;
}

.heroImg {
  width: 100%;
  height: 100%;
  max-width: 1200px;
  overflow: hidden;
  padding-inline: 0px;

  .q-img__content > div {
    padding: 0;
  }
}

.imgContent {
  position: relative;
  width: 100%;
  left: 0%;
  top: 55%;
  transform: translate(0%, -50%);
  text-align: center;

  @media (min-width: $breakpoint-sm-max) {
    width: 84%;
    left: 7%;
    text-align: left;
  }
}

.tagline {
  text-wrap: balance;
  font-weight: bold;
  margin: auto;
  margin-bottom: 0.5rem;
  width: 95%;

  @media only screen and (min-width: $breakpoint-sm-max) {
    font-weight: bold;
    margin: 0rem 0rem 1rem 0rem;
  }
}

.subTagline {
  width: 80%;
  margin: auto;
  color: #ccc;
  margin-bottom: 2rem;

  @media only screen and (min-width: $breakpoint-sm-max) {
    width: 90%;
    margin: 0rem 0rem 3rem;
  }
}
</style>

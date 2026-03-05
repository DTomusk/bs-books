<template>
  <div>
    <div class="row justify-between no-select">
      <div class="heroContainer">
        <HeroSection @learnMore="scrollToLearnMore" />
      </div>
      <div
        v-if="!isMobile"
        class="flex-centre-col col rightContainer bg-offwhite"
        style="width: 50%"
      >
        <AuthPage
          :mode="isLoggingIn ? 'login' : 'register'"
          :isOwnPage="false"
          @toggleMode="isLoggingIn = !isLoggingIn"
        />
      </div>
    </div>
    <div class="flex-centre-row">
      <div class="bg-primary contentContainer">
        <HowItWorks ref="howItWorksComp" />
        <BookCarousel title="Our top picks" :books="featuredBooks" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import HeroSection from 'src/pages/LandingPage/HeroSection.vue';
import BookCarousel from 'src/components/book/BookCarousel.vue';
import HowItWorks from 'src/pages/LandingPage/HowItWorks.vue';
import { useScreen } from 'src/composables/useScreen';
import type { Book } from 'src/components/models';
import AuthPage from 'src/pages/AuthPage.vue';
import { scroll } from 'quasar';
import { ref } from 'vue';

const { getScrollTarget, setVerticalScrollPosition } = scroll;
const howItWorksComp = ref<InstanceType<typeof HowItWorks> | null>(null);
const { isMobile } = useScreen();
const isLoggingIn = ref(true);
const featuredBooks: Book[] = [
  {
    image: 'https://cdn.ecommercedns.uk/files/1/238331/3/11640523/9780440867579.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
  {
    image: 'https://fordhamobserver.com/wp-content/uploads/2012/05/Arts_50-Shades-of-grey1.jpg',
    title: '50 Shades of Grey',
  },
  {
    image: 'https://m.media-amazon.com/images/I/712cDO7d73L.jpg',
    title: 'The Hobbit',
  },
  {
    image: 'https://cdn.waterstones.com/bookjackets/large/9780/2419/9780241952443.jpg',
    title: 'Anne Frank',
  },
  {
    image: 'https://cdn.waterstones.com/bookjackets/large/9781/4071/9781407132082.jpg',
    title: 'The Hunger Games',
  },
  {
    image: 'https://images.booksense.com/images/324/906/9781728906324.jpg',
    title: 'Communist Manifesto',
  },
  {
    image: 'https://m.media-amazon.com/images/I/81suTb0slVL._AC_UF894,1000_QL80_.jpg',
    title: 'Boris Johnson Unleashed',
  },
  {
    image:
      'https://i0.wp.com/vilmairis.com/wp-content/uploads/2019/01/After-new-cover.jpg?resize=782%2C1215',
    title: 'After',
  },
];

function scrollToLearnMore() {
  const el = howItWorksComp.value?.howItWorksEl;
  if (!el) return;

  const scrollTarget = getScrollTarget(el);
  // -60 gives padding to header which is scrolled to so it isn't up against top of the screen
  setVerticalScrollPosition(scrollTarget, el.offsetTop - 60, 500);
}
</script>

<style lang="scss" scoped>
.heroContainer {
  position: relative;
  width: 100%;

  @media (min-width: $breakpoint-sm-max) {
    width: 65%;
    height: 100dvh;
    max-width: 1200px;
    overflow: visible;
    filter: drop-shadow(6px 0 6px rgb(0, 0, 0, 0.75));
  }
}

.contentContainer {
  max-width: 90vw;
  margin-inline: 1rem;
  padding-block: 0rem 3rem;

  @media (min-width: $breakpoint-sm-max) {
    max-width: 1200px;
    padding-block: 1rem 5rem;
    padding-inline: 3rem;
  }
}
</style>

<template>
  <div>
    <q-img
      src="https://img.pixers.pics/pho_wat(s3:700/FO/65/70/77/40/700_FO65707740_a1f7177e6af5a55410a721d7e08b3400.jpg,700,689,cms:2018/10/5bd1b6b8d04b8_220x50-watermark.png,over,480,639,jpg)/wall-murals-bookshelf-full-with-books.jpg.jpg"
    >
      <div class="top-bottom-gradient">
        <div class="flex-left logoContainer" style="position: relative">
          <q-img src="src/assets/LogoRectangleLong.png" width="55%" />
          <div :style="{ display: 'flex', gap: $q.screen.lt.sm ? '0.5rem' : '1rem' }">
            <q-btn
              :size="$q.screen.lt.sm ? 'sm' : 'md'"
              icon="fa-regular fa-user"
              color="accent"
              text-color="black"
              round
            />
            <q-btn
              :size="$q.screen.lt.sm ? 'sm' : 'md'"
              icon="search"
              color="accent"
              text-color="black"
              round
            />
          </div>
        </div>
        <div
          class="flex-centre-col"
          style="height: 100%; position: absolute; top: 0; margin-top: 1rem"
        >
          <p class="tagline">The book rating system you never knew you needed</p>
          <p class="subTagline">
            Rate books using <span><q-icon name="fa-solid fa-heart" /></span> and
            <span><q-icon name="fa-solid fa-poop" /></span> because some books deserve more than
            just love
          </p>
          <q-btn
            label="Create an account"
            color="positive"
            text-color="black"
            class="getStartedButton"
          />
        </div>
      </div>
    </q-img>
    <div class="searchContainer flex-centre-col bold">
      <p class="title" style="margin-block: 1rem 0rem">How it works</p>
      <ul style="margin-top: 0.5rem; margin-right: 1rem">
        <li>
          Rate books using hearts and poos; Hearts are good, poos are bad. Some books deserve some
          of each!
        </li>
        <li>
          Create your own digital library of books which are brilliant, bad, or strangely both
        </li>
        <li>Contribute to our database of books by adding covers, descriptions, authors, etc...</li>
      </ul>
    </div>
    <div>
      <BookCarousel title="Our top picks" :books="featuredBooks" />
    </div>
  </div>
</template>

<script setup lang="ts">
import BookCarousel from 'src/components/book/BookCarousel.vue';
import type { Book } from 'src/components/models';
import { useQuasar } from 'quasar';

const $q = useQuasar();
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
</script>

<style lang="scss" scoped>
.top-bottom-gradient {
  background: #000000;
  background: linear-gradient(
    180deg,
    rgba(0, 0, 0, 0.9) 0%,
    rgba(0, 0, 0, 0.7) 50%,
    rgba(0, 0, 0, 0.9) 100%
  );
  height: 100%;
  width: 100vw;
  position: absolute;
  bottom: 0;
}

.searchBar {
  width: 80%;
  margin-block: 0rem 1rem;
}

.searchContainer {
  width: 100%;
  height: 100%;
  background-color: white;
}

.logoContainer {
  width: 100%;
  margin-block: 0.75rem 3rem;
  justify-content: space-between;
  align-items: center;
}

.getStartedButton {
  margin-top: 1rem;
}

.tagline {
  font-size: clamp(3vw, 5vw, 40px);
  font-weight: bold;
  text-align: center;
  text-wrap: balance;
  margin-bottom: 0.5rem;
}

.subTagline {
  font-size: clamp(2vw, 3vw, 30px);
  text-align: center;
  text-wrap: balance;
}
</style>

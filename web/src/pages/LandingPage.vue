<template>
  <div>
    <q-img
      src="https://img.pixers.pics/pho_wat(s3:700/FO/65/70/77/40/700_FO65707740_a1f7177e6af5a55410a721d7e08b3400.jpg,700,689,cms:2018/10/5bd1b6b8d04b8_220x50-watermark.png,over,480,639,jpg)/wall-murals-bookshelf-full-with-books.jpg.jpg"
    >
      <div class="top-bottom-gradient">
        <div class="flex-left logoContainer">
          <q-img src="src/assets/LogoRectangleLong.png" width="50%" />
        </div>
        <div class="flex-centre-col" style="padding-top: 1rem">
          <p class="tagline">The book rating system you never knew you needed</p>
          <q-btn label="Get started" color="tertiary" class="getStartedButton" />
        </div>
      </div>
    </q-img>
    <div class="searchContainer flex-centre-col bold">
      <p class="title" style="margin-block: 1rem 0rem">How it works</p>
      <ul style="margin-top: 0.5rem">
        <li>Create an account or log in</li>
        <li>Rate books using hearts and poos (because opinions are complicated)</li>
        <li>Contribute to our database of books by adding covers, descriptions, authors, etc...</li>
      </ul>
      <q-input
        v-model="searchValue"
        placeholder="Search from 1000s of books"
        bg-color="white"
        outlined
        class="searchBar"
      >
        <template v-slot:append> <q-btn icon="search" round dense flat /> </template
      ></q-input>
    </div>
    <div>
      <BookCarousel title="Featured" :books="featuredBooks" />
    </div>
  </div>
</template>

<script setup lang="ts">
import BookCarousel from 'src/components/book/BookCarousel.vue';
import type { Book } from 'src/components/models';
import { ref } from 'vue';

const featuredBooks: Book[] = [
  {
    image:
      'https://www.stylist.co.uk/images/app/uploads/2017/06/12153937/91mpkwml6il-1518x1012.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
  {
    image:
      'https://www.stylist.co.uk/images/app/uploads/2017/06/12153937/91mpkwml6il-1518x1012.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
  {
    image:
      'https://www.stylist.co.uk/images/app/uploads/2017/06/12153937/91mpkwml6il-1518x1012.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
  {
    image:
      'https://www.stylist.co.uk/images/app/uploads/2017/06/12153937/91mpkwml6il-1518x1012.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
  {
    image:
      'https://www.stylist.co.uk/images/app/uploads/2017/06/12153937/91mpkwml6il-1518x1012.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
  {
    image:
      'https://www.stylist.co.uk/images/app/uploads/2017/06/12153937/91mpkwml6il-1518x1012.jpg',
    title: 'Tracy Beaker',
  },
  {
    image: 'https://images.wolfgangsvault.com/m/large/ZZZ025983-BK/twilight-book-2005.jpg',
    title: 'Twilight',
  },
];
// const slide = ref(featuredBooks[0]?.title);
const searchValue = ref('');
</script>

<style lang="scss" scoped>
.top-bottom-gradient {
  background: #000000;
  background: linear-gradient(
    180deg,
    rgba(0, 0, 0, 0.8) 0%,
    rgba(0, 0, 0, 0.7) 50%,
    rgba(0, 0, 0, 0.8) 100%
  );
  height: 100%;
  width: 100vw;
  position: absolute;
  bottom: 0;
}

.top-gradient {
  background: #000000;
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0.4) 100%);
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
}

.getStartedButton {
  margin-top: 1rem;
}

.tagline {
  font-size: x-large;
  font-weight: bold;
  text-align: center;
  text-wrap: balance;
}
</style>

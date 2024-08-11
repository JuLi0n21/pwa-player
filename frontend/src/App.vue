<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import NowPlaying from './components/NowPlaying.vue'
import Footer from './components/Footer.vue'
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { defineStore } from 'pinia'
import { useHeaderStore } from '@/stores/headerStore';

const headerStore = useHeaderStore();

const showFooter = ref(true);
const showNowPlaying = ref(true);
const route = useRoute();
function hide() {
  showHeader.value = false;
}

watch(route, async (to) => {

  if (route.path.startsWith("/nowplaying")) {
    showNowPlaying.value = false;
  } else {
    showNowPlaying.value = true;
  }

  if (route.path.startsWith("/menu")) {
    headerStore.hide();
  } else {
    headerStore.show();
  }

})

</script>

<template>

  <contentw class="flex flex-col h-screen max-h-screen wrapper text-pink-500 text-xl">

    <RouterView />

    <Transition>
      <NowPlaying v-show="showNowPlaying" />
    </Transition>

    <Footer />

  </contentw>
</template>

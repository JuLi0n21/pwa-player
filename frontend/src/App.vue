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

function loadColors() {

  document.documentElement.style.setProperty('--background-color', localStorage.getItem('bgColor') || '#1c1719');
  document.documentElement.style.setProperty('--action-color', localStorage.getItem('actionColor') || '#eab308');

  document.documentElement.style.setProperty('--information-color', localStorage.getItem('infoColor') || '#ec4899');
  document.documentElement.style.setProperty('--border-color', localStorage.getItem('borderColor') || '#ec4899');

  console.log(localStorage.getItem('bgColor'));
}

loadColors();

</script>

<template>

  <contentw class="flex flex-col h-screen max-h-screen wrapper info text-xl">

    <RouterView />

    <Transition>
      <NowPlaying v-show="showNowPlaying" />
    </Transition>

    <Footer />

  </contentw>
</template>

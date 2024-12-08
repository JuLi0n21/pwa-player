<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import NowPlaying from '@/components/NowPlaying.vue'
import NowPlayingView from '@/views/NowPlayingView.vue'
import MenuView from '@/views/MenuView.vue'
import Footer from '@/components/Footer.vue'
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { defineStore } from 'pinia'
import { useHeaderStore } from '@/stores/headerStore';
import { isMobile, isPc } from './script/utils.ts'

const headerStore = useHeaderStore();

const showFooter = ref(true);
const showNowPlaying = ref(true);
const route = useRoute();

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

const screenInfo = ref({
  isSmall: false,
  isMedium: false,
});

const checkScreenSize = () => {
  screenInfo.value.isSmall = isMobile();
  screenInfo.value.isMedium = isPc();

};

onMounted(() => {
  checkScreenSize(); // Initial check


  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  // Clean up the listener when component is destroyed
  window.removeEventListener('resize', checkScreenSize);
});

</script>

<template>


  <div v-if="screenInfo.isSmall" class="flex flex-col h-screen max-h-screen wrapper info text-xl ">

    <RouterView />

    <NowPlaying v-show="showNowPlaying" />

    <Footer />
  </div>
  <div v-if="!screenInfo.isSmall" class="flex flex-col h-screen max-h-screen wrapper info text-xl ">
    <main class="flex flex-1 w-full h-full overflow-y-scroll">
      <div class="w-1/12">
        <MenuView />
      </div>
      <div class="flex-1 overflow-y-scroll">
        <RouterView />
      </div>
      <div class="w-1/5">
        <NowPlayingView />
      </div>
    </main>

    <Footer />
  </div>

</template>

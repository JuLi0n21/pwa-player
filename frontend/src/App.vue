<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import NowPlaying from '@/components/NowPlaying.vue'
import NowPlayingView from '@/views/NowPlayingView.vue'
import MenuView from '@/views/MenuView.vue'
import Footer from '@/components/Footer.vue'
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { isMobile, isPc } from './script/utils.ts'
import { useAudio } from './composables/useAudio.ts'

const showNowPlaying = ref(true);
const route = useRoute();
const audioRef = ref<HTMLAudioElement | null>(null)
const audio = useAudio(audioRef);

watch(route, async (to) => {

  if (route.path.startsWith("/nowplaying")) {
    showNowPlaying.value = false;
  } else {
    showNowPlaying.value = true;
  }
/*
  if (route.path.startsWith("/menu")) {
    headerStore.hide();
  } else {
    headerStore.show();
  }
*/
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
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>
<template>
  <div v-if="screenInfo.isSmall" class="flex flex-col h-screen max-h-screen wrapper info text-xl">
    <RouterView />
    <NowPlaying v-show="showNowPlaying" />
    <Footer />
  </div>

  <div v-else class="flex flex-col h-screen max-h-screen wrapper info text-xl">
    <main class="flex flex-1 w-full h-full overflow-y-scroll">
      
      <aside class="w-1/12 bg-primary p-4">
        <p>Sidebar content</p>
      </aside>

      <section class="flex-1 overflow-y-scroll p-4">
        <RouterView />
      </section>

      <section class="w-1/5 p-4">
        <NowPlayingView />
      </section>
    </main>
    <Footer />
  </div>
</template>
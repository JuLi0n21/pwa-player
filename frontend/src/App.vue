<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
import NowPlaying from "@/components/NowPlaying.vue";
import NowPlayingView from "@/views/NowPlayingView.vue";
import MenuView from "@/views/MenuView.vue";
import HistoryView from "@/views/HistoryView.vue";
import Footer from "@/components/Footer.vue";
import { ref, onMounted, watch, onUnmounted } from "vue";
import { useRoute } from "vue-router";
import { isMobile, isPc } from "./script/utils.ts";

const showNowPlaying = ref(true);
const route = useRoute();

watch(route, async (to) => {
  if (route.path.startsWith("/nowplaying")) {
    showNowPlaying.value = false;
  } else {
    showNowPlaying.value = true;
  }
});

function loadColors() {
  document.documentElement.style.setProperty("--background-color", localStorage.getItem("bgColor") || "#1c1719");
  document.documentElement.style.setProperty("--action-color", localStorage.getItem("actionColor") || "#eab308");

  document.documentElement.style.setProperty("--information-color", localStorage.getItem("infoColor") || "#ec4899");
  document.documentElement.style.setProperty("--border-color", localStorage.getItem("borderColor") || "#ec4899");
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
  window.addEventListener("resize", checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener("resize", checkScreenSize);
});
</script>
<template>
  <div v-if="screenInfo.isSmall" class="flex flex-col h-screen max-h-screen overflow-y-auto text-xl wrapper info">
    <RouterView />
    <NowPlaying v-show="showNowPlaying" />
    <Footer />
  </div>

  <div v-else class="flex flex-col h-screen max-h-screen text-xl wrapper info">
    <main class="flex flex-1 w-full h-full overflow-y-hidden">
      <aside class="bg-primary p-4 w-1/12 overflow-y-auto">
        <HistoryView />
      </aside>

      <section class="flex-1 overflow-y-hidden">
        <RouterView />
      </section>

      <section class="flex flex-col w-1/5 overflow-y-auto">
        <NowPlayingView />
      </section>
    </main>
    <Footer />
  </div>
</template>

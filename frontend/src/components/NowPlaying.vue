<script setup lang="ts">
import { useAudio } from "@/composables/useAudio";
import { onMounted, computed } from "vue";
import { useUser } from "@/composables/useUser";

const audioStore = useAudio();
const userStore = useUser();

const title = computed(() => audioStore.currentSong.value?.name || "Unknown Title");
const artist = computed(() => audioStore.currentSong.value?.artist || "Unknown Artist");
const bgimg = computed(() => audioStore.currentSong.value?.previewimage || "/default-bg.jpg");

onMounted(() => {
  audioStore.init();
});
</script>

<template>
  <div>
    <hr />
    <div class="relative p-1 wrapper action">
      <img
        :src="encodeURI(`${userStore.cloudflareUrl.value}${bgimg}?h=150&w=400`)"
        class="top-0 right-0 left-0 absolute w-full h-full"
        :style="{ filter: 'blur(2px)', opacity: '0.5' }"
        alt="Background Image"
      />

      <nav class="z-10 relative flex-col">
        <div class="flex justify-between">
          <RouterLink to="/nowplaying" class="overflow-hidden grow">
            <p class="relative overflow-hidden font-bold text-sm text-left text-ellipsis text-nowrap info">
              {{ title }}
            </p>

            <p class="relative font-bold text-sm text-left text-nowrap info">
              {{ artist }}
            </p>
          </RouterLink>
          <div class="flex flex-col justify-center px-2 text-center" @click="audioStore.togglePlay">
            <i
              :class="[audioStore.isPlaying.value ? ' fa-circle-pause' : 'fa-circle-play']"
              class="text-3xl fa-regular"
            ></i>
          </div>
        </div>

        <div class="bg-gray-200 dark:bg-gray-700 rounded-full w-full h-0.5">
          <div
            class="bg-blue-600 dark:bg-yellow-500 rounded-full h-0.5"
            :style="{ width: audioStore.percentDone.value + '%' }"
          ></div>
        </div>
      </nav>
    </div>
  </div>
</template>

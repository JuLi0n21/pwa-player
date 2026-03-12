<script setup lang="ts">
import { computed } from "vue";
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";

const audioStore = useAudio();
const userStore = useUser();
const title = computed(() => audioStore.currentSong.value?.name || "Unknown Title");
const artist = computed(() => audioStore.currentSong.value?.artist || "Unknown Artist");
const bgimg = computed(() => {
  const preview = audioStore.currentSong.value?.previewimage;
  return preview ? encodeURI(`${userStore.cloudflareUrl.value}${preview}`) : "/default-bg.jpg";
});
</script>

<template>
  <header>
    <div class="wrapper">
      <div class="relative">
        <nav class="flex flex-1 justify-start space-x-1 mx-1 my-2">
          <RouterLink class="shadow-xl backdrop--light p-1 rounded-full" to="/"
            ><i class="fa-arrow-left fa-solid"></i>
          </RouterLink>
          <h1 class="right-0 left-0 absolute text-center">Now Playing</h1>
        </nav>
      </div>
      <hr />
    </div>
  </header>

  <main class="flex flex-col flex-1 justify-center items-center px-4 text-center">
    <div class="flex flex-col items-center space-y-6 w-full max-w-md">
      <div class="relative w-full aspect-square">
        <img
          class="absolute inset-0 shadow-lg rounded-lg w-full h-full object-cover"
          :src="encodeURI(`${bgimg}?h=320&w=320`)"
          :key="bgimg"
          alt="Album Art"
        />
        <i class="absolute inset-0 flex justify-center items-center text-white text-5xl">
          <i class="bg-black bg-opacity-50 p-4 rounded-full fa-solid fa-play"></i>
        </i>
      </div>

      <div class="flex justify-between items-center space-x-6 w-full text-3xl">
        <i class="fa-solid fa-backward-step" @click="audioStore.togglePrevious"></i>
        <i
          :class="[audioStore.isPlaying.value ? 'fa-circle-pause' : 'fa-circle-play']"
          class="text-5xl fa-regular"
          @click="audioStore.togglePlay"
        ></i>
        <i class="fa-solid fa-forward-step" @click="audioStore.toggleNext"></i>
      </div>

      <div class="px-2 w-full text-center">
        <p class="font-semibold text-lg truncate">{{ title }}</p>
        <RouterLink :to="'search?a=' + artist" class="block text-blue-500 text-sm truncate">
          {{ artist }}
        </RouterLink>
      </div>

      <div class="flex justify-between items-center px-4 w-full">
        <i
          @click="audioStore.toggleShuffle"
          :class="[audioStore.shuffle.value ? 'text-yellow-500' : '']"
          class="fa-solid fa-shuffle"
        ></i>
        <i
          @click="audioStore.toggleRepeat"
          :class="[audioStore.repeat.value ? 'text-yellow-500' : '']"
          class="fa-solid fa-repeat"
        ></i>
        <i @click="$router.go(-1)" class="fa-solid fa-arrow-down"></i>
      </div>

      <div class="px-4 w-full">
        <input
          class="bg-yellow-200 bg-opacity-20 rounded-full outline-none w-full h-2 accent-yellow-600 appearance-none"
          type="range"
          @input="audioStore.updateTime(Number(($event.target as HTMLInputElement).value) || 0)"
          :max="100"
          step="0.001"
          :value="audioStore.percentDone.value"
        />
      </div>

      <div class="flex justify-between px-4 w-full text-sm">
        <span>{{ audioStore.currentTime.value }}</span>
        <span>{{ audioStore.duration.value }}</span>
      </div>
    </div>
  </main>
</template>

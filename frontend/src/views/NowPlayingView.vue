<script setup lang="ts">
import { computed } from 'vue';
import { useAudio } from '@/composables/useAudio';

const audioStore = useAudio();

const title = computed(() => audioStore.currentSong.value?.name || 'Unknown Title')
const artist = computed(() => audioStore.currentSong.value?.artist || 'Unknown Artist')
const bgimg = computed(() => audioStore.currentSong.value?.previewimage || '/default-bg.jpg')

</script>

<template>
  <header>
    <div class="wrapper">
      <div class="relative">
        <nav class="flex flex-1 justify-start my-2 mx-1 space-x-1">
          <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/"><i class="fa-solid fa-arrow-left"></i>
          </RouterLink>
          <h1 class="absolute left-0 right-0 text-center"> Now Playing </h1>
        </nav>
      </div>
      <hr>
    </div>
  </header>

<main class="flex-1 flex flex-col items-center justify-center text-center px-4">
  <div class="flex flex-col items-center w-full max-w-md space-y-6">

    <div class="relative w-full aspect-square">
      <img
        class="absolute inset-0 w-full h-full object-cover rounded-lg shadow-lg"
        :src="encodeURI(bgimg + '?h=320&w=320')"
        :key="bgimg"
        alt="Album Art"
      />
      <i class="absolute inset-0 flex items-center justify-center text-white text-5xl">
        <i class="fa-solid fa-play bg-black bg-opacity-50 p-4 rounded-full"></i>
      </i>
    </div>

    <div class="flex justify-between items-center w-full text-3xl space-x-6">
      <i class="fa-solid fa-backward-step" @click="audioStore.togglePrevious"></i>
      <i
        :class="[audioStore.isPlaying.value ? 'fa-circle-pause' : 'fa-circle-play']"
        class="fa-regular text-5xl"
        @click="audioStore.togglePlay"
      ></i>
      <i class="fa-solid fa-forward-step" @click="audioStore.toggleNext"></i>
    </div>

    <div class="text-center w-full px-2">
      <p class="truncate text-lg font-semibold">{{ title }}</p>
      <RouterLink :to="'search?a=' + artist" class="block text-sm text-blue-500 truncate">
        {{ artist }}
      </RouterLink>
    </div>

    <div class="flex justify-between items-center w-full px-4">
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

    <div class="w-full px-4">
      <input
        class="w-full appearance-none h-2 rounded-full bg-yellow-200 bg-opacity-20 accent-yellow-600 outline-none"
        type="range"
        @input="event => audioStore.updateTime(Number(event.target.value))"
        :max="100"
        step="0.001"
        :value="audioStore.percentDone.value"
      />
    </div>

    <div class="flex justify-between text-sm w-full px-4">
      <span>{{ audioStore.currentTime.value }}</span>
      <span>{{ audioStore.duration.value }}</span>
    </div>
  </div>
</main>
</template>

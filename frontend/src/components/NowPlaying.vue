<script setup lang="ts">
import { ref } from 'vue'
import { useAudioStore } from '@/stores/audioStore';

const audioStore = useAudioStore();
</script>


<template>
  <div>
    <hr>
    <div class="relative wrapper p-1 grow action">
      <img :src="encodeURI(audioStore.bgimg + '?h=150&w=400')" class="w-full absolute top-0 left-0 right-0 h-full"
        :style="{ 'filter': 'blur(2px)', 'opacity': '0.5' }" alt="Background Image" />

      <nav class="relative flex-col z-10">

        <div class="flex justify-between">
          <RouterLink to="/nowplaying" class="grow overflow-hidden">
            <p class="relative  text-sm text-left font-bold info overflow-hidden text-ellipsis text-nowrap">
              {{ audioStore.title }}
            </p>

            <p class="relative text-sm text-left font-bold info text-nowrap">
              {{ audioStore.artist }}
            </p>
          </RouterLink>
          <div class="flex flex-col text-center justify-center px-2" @click="audioStore.togglePlay">
            <i :class="[audioStore.isPlaying ? ' fa-circle-play' : 'fa-circle-pause']" class="text-3xl fa-regular"></i>
          </div>

        </div>

        <div class="w-full bg-gray-200 rounded-full h-0.5 dark:bg-gray-700">
          <div class="bg-blue-600 h-0.5 rounded-full dark:bg-yellow-500"
            :style="{ 'width': audioStore.percentDone + '%' }">
          </div>
        </div>

      </nav>
      <audio controls class="hidden" id="audio-player" :src="audioStore.songSrc"
        @timeupdate="audioStore.update"></audio>
    </div>
  </div>
</template>

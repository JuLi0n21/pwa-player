<script setup lang="ts">
import { ref } from 'vue'
import { useAudioStore } from '@/stores/audioStore';

const audioStore = useAudioStore();
</script>


<template>
  <div>
    <hr>
    <div class="relative wrapper p-1 grow text-yellow-500">

      <img :src="encodeURI(audioStore.bgimg)" class="absolute top-0 left-0 w-full h-full"
       :style="{ 'filter': 'blur(2px)', 'opacity': '0.5' }" alt="Background Image" />

      <nav class="relative flex justify-around my-2 z-10">

        <div class=" grow flex flex-col justify-around text-3xl text-center" to="/menu">
          <i @click="audioStore.toggleShuffle" :class="[audioStore.shuffle ? 'text-pink-500' : '']"
            class="fa-solid fa-shuffle"></i>
        </div>

        <div>
          <div class="flex-col grow text-center justify-center hover:text-pink-500" @click="audioStore.togglePlay">
            <i :class="[audioStore.isPlaying ? ' fa-circle-play' : 'fa-circle-pause']" class="text-7xl fa-regular"></i>
          </div>
        </div>

        <div class="grow flex flex-col justify-around text-3xl text-center hover:text-pink-500" to="/menu">
          <i @click="audioStore.toggleRepeat" :class="[audioStore.repeat ? 'text-pink-500' : '']"
            class="fa-solid fa-repeat "></i>
        </div>

      </nav>
      <RouterLink class="absolute right-2 bottom-2" to="/nowplaying">
        <i class="fa-solid fa-arrow-up"></i>
      </RouterLink>

      <marquee class="relative mx-16 text-2xl font-bold text-pink-500" behavior="scroll">{{ audioStore.artist }} - {{ audioStore.title
        }}</marquee>
      <audio controls class="hidden" id="audio-player" :src="audioStore.songSrc"
        @timeupdate="audioStore.update"></audio>
    </div>
  </div>
</template>

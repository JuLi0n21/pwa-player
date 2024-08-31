<script setup lang="ts">
import { ref } from 'vue'
import { useAudioStore } from '@/stores/audioStore';

const audioStore = useAudioStore();
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

  <main class="flex-1 flex justify-center text-center action">
    <div class="flex flex-col justify-around">
      <div class="relative">
        <i class="relative p-36 fa-solid fa-play">

          <img class="absolute top-4 left-0 bottom-0 right-0 bg-center bg-cover rounded-lg"
            :src="encodeURI(audioStore.bgimg + '?h=320&w=320')" :key="audioStore.bgimg" />
        </i>
      </div>

      <div class="h-1/3 flex flex-col justify-center">
        <div class="flex-1"></div>
        <div>
          <div class="flex w-screen justify-around">
            <i class="fa-solid fa-backward-step text-5xl self-center" @click="audioStore.togglePrev"></i>
            <i :class="[audioStore.isPlaying ? 'fa-circle-play' : 'fa-circle-pause']" class="fa-regular text-7xl "
              @click="audioStore.togglePlay"></i>
            <i class="fa-solid fa-forward-step text-5xl self-center" @click="audioStore.toggleNext"></i>
          </div>
        </div>
        <div class="flex flex-1 justify-around ml-4">
          <i @click="audioStore.toggleShuffle" :class="[audioStore.shuffle ? 'info' : '']"
            class="fa-solid fa-shuffle"></i>

          <div class="m-4 info flex-1 overflow-idden">
            <p>{{ audioStore.title }}</p>
            <RouterLink :to="'search?a=' + audioStore.artist">

              {{ audioStore.artist }}

            </RouterLink>
          </div>
          <div class="flex flex-col justify-between mb-4 mr-4">
            <i @click="audioStore.toggleRepeat" :class="[audioStore.repeat ? 'info' : '']"
              class="fa-solid fa-repeat"></i>
            <i @click="$router.go(-1);" class="fa-solid fa-arrow-down"></i>
          </div>
        </div>
        <div class="flex">
          <input
            class="appearance-none mx-4 flex-1 bg-yellow-200 bg-opacity-20 accent-yellow-600 rounded-lg outline-none slider "
            type="range" id="audio-slider" @change="audioStore.updateTime" max="100" :value="audioStore.percentDone">
        </div>
        <div class="flex justify-between mx-4">
          <span id="current-time" class="time">{{ audioStore.currentTime }}</span>
          <span id="duration" class="time ">{{ audioStore.duration }}</span>
        </div>
      </div>
    </div>

  </main>
</template>

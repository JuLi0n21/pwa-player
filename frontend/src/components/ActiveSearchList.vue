<script setup lang="ts">
import type { Song } from "../script/types";
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";

const audioStore = useAudio();
const userStore = useUser();

const props = defineProps<{
  songs: Song[];
  artist: string[];
  search: string;
}>();

function playSong(hash: string) {
  const selected = props.songs.find((s) => s.hash === hash);
  if (selected) {
    audioStore.setSong(selected);
  }
}

function highlightText(text: string, searchterm: string) {
  if (!searchterm) return text;
  const regex = new RegExp(`(${searchterm})`, "gi");
  return text.replace(regex, '<span class="font-bold text-yellow-400">$1</span>');
}
</script>

<template>
  <div class="w-full h-full text-xs">
    
    <div v-if="props.artist && props.artist.length > 0" class="mb-4">
      <h2 class="px-4 py-2 font-semibold text-lg uppercase tracking-wider action">Artists</h2>
      <ul class="space-y-1 px-2">
        <li v-for="(artist, index) in props.artist" :key="index">
          <RouterLink 
            class="flex items-center bg-white/5 hover:bg-white/10 p-3 border border-transparent rounded-xl transition-colors hover:bordercolor"
            :to="'/search?a=' + artist"
          >
            <div class="flex justify-center items-center bg-yellow-500/20 mr-3 rounded-full w-10 h-10">
              <i class="text-yellow-500 fa-solid fa-user"></i>
            </div>
            <span class="text-base info" v-html="highlightText(artist, props.search)"></span>
          </RouterLink>
        </li>
      </ul>
    </div>

    <div v-if="props.songs && props.songs.length > 0">
      <h2 class="px-4 py-2 font-semibold text-lg uppercase tracking-wider action">Songs</h2>
      <ul class="space-y-1 px-2 pb-20"> <li v-for="(song, index) in props.songs" :key="index">
          <button 
            @click="playSong(song.hash)" 
            class="flex items-center p-2 rounded-xl w-full text-left transition-colors"
          >
            <img
              :src="song.previewimage 
                ? `${userStore.cloudflareUrl.value}${song.previewimage}?h=120&w=120` 
                : '/default-bg.png'"
              class="rounded-lg w-12 h-12 object-cover shrink-0"
              loading="lazy"
            />
            <div class="flex-1 ml-3 overflow-hidden">
              <p class="text-base truncate info">
                <span v-html="highlightText(song.name, search)"></span>
              </p>
              <p class="opacity-60 text-sm truncate action">
                <span v-html="highlightText(song.artist, props.search)"></span>
              </p>
            </div>
            <i class="opacity-0 group-hover:opacity-100 ml-2 pr-2 transition-opacity fa-solid fa-play"></i>
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
:deep(.text-yellow-400) {
  color: #facc15;
}
</style>
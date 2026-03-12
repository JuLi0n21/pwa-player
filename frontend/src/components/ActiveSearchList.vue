<script setup lang="ts">
import type { Song, CollectionPreview } from "../script/types";
import { useAudio } from "@/composables/useAudio";
import { ref } from "vue";
import { RouterLink } from "vue-router";
import { useUser } from "@/composables/useUser";
const audioStore = useAudio();
const userStore = useUser();

const props = defineProps<{
  songs: Song[];
  artist: string[];
  search: string;
}>();

function update(hash: string) {
  audioStore.setSong(props.songs[props.songs.findIndex((s) => s.hash == hash)]);
}

function highlightText(text: string, searchterm: string) {
  if (!searchterm) return text;
  const regex = new RegExp(`(${searchterm})`, "gi");
  return text.replace(regex, '<span style="color: yellow;">$1</span>');
}
</script>

<template>
  <div class="border rounded-lg w-full h-full overflow-scroll overflow-x-hidden text-xs bordercolor bg">
    <div v-if="props.artist && props.artist.length > 0" class="border bordercolor">
      <h2 class="text-2xl action">Artists</h2>
      <ul>
        <li v-for="(artist, index) in props.artist" :key="index" class="rounded-lg">
          <RouterLink class="flex" :to="'/search?a=' + artist" v-html="highlightText(artist, props.search)">
          </RouterLink>
        </li>
      </ul>
    </div>
    <div v-if="props.songs && props.songs.length > 0" class="border bordercolor">
      <h2 class="text-2xl action">Songs</h2>
      <ul>
        <li v-for="(song, index) in props.songs" :key="index" class="rounded-lg">
          <button @click="update(song.hash)" class="flex">
            <img :src="encodeURI(`${userStore.cloudflareUrl.value}/${song.previewimage ? song.previewimage + '?h=120&w=120' : '/default-bg.png'}`)" class="w-12 h-12" />
            <p class="ml-2 overflow-hidden text-ellipsis text-nowrap">
              <span v-html="highlightText(song.name, search)"></span> -
              <span v-html="highlightText(song.artist, props.search)"></span>
            </p>
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

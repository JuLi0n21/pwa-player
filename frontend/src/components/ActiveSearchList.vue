<script setup lang="ts">
import { updateLanguageServiceSourceFile } from 'typescript';
import type { Song, CollectionPreview } from '../script/types'
import { useAudioStore } from '@/stores/audioStore';
import { ref } from 'vue';
import { RouterLink } from 'vue-router';

const audioStore = useAudioStore()

const props = defineProps<{
  songs: Song[];
  artist: string[];
  search: string;
}>();

function update(hash: string) {

  audioStore.setSong(props.songs.at(props.songs.findIndex(s => s.hash == hash)))
}

function highlightText(text: string, searchterm: string) {
  if (!searchterm) return text;
  const regex = new RegExp(`(${searchterm})`, 'gi');
  return text.replace(regex, '<span style="color: yellow;">$1</span>');
}

</script>

<template>
  <div class="h-full w-full overflow-scroll overflow-x-hidden border bordercolor rounded-lg text-xs bg">
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
            <img :src="encodeURI(song.previewimage + '?h=120&w=120')" class="h-12 w-12"></img>
            <p class="text-nowrap text-ellipsis overflow-hidden ml-2">
              <span v-html="highlightText(song.name, search)"></span> - <span
                v-html="highlightText(song.artist, props.search)"></span>
            </p>
          </button>
        </li>
      </ul>
    </div>

  </div>
</template>

<script setup lang="ts">
import { updateLanguageServiceSourceFile } from 'typescript';
import type { Song, CollectionPreview } from '../script/types'
import { useAudioStore } from '@/stores/audioStore';
import { ref } from 'vue';
import { RouterLink } from 'vue-router';

const audioStore = useAudioStore()

const props = defineProps<{
   songs :Song[];
   artist: string[];
}>();

function update(hash: string){

audioStore.setSong(props.songs.at(props.songs.findIndex(s => s.hash==hash)))
}
</script>

<template>
<div class="h-full overflow-scroll border border-pink-500 rounded-lg bg-gray-800">
    <h3>Artists</h3>
    <ul>
      <li v-for="(artist, index) in props.artist" :key="index">
        <RouterLink class="flex"  :to="'/search?a=' + artist">{{ artist }}</RouterLink>
      </li>
    </ul>
        <h3>Songs</h3>
    <ul>
      <li v-for="(song, index) in props.songs" :key="index">
        <button @click="update(song.hash)" class="flex">
            <img class="h-12 w-12" :src="song.previewimage">
            <div class="flex flex-col">
            {{ song.name }} - {{ song.artist }}
            </div>
        </button>
    </li>
    </ul>

</div>
</template>
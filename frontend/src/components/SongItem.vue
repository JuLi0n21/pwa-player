<script setup lang="ts">
import { ref, defineProps } from 'vue'
import { useAudioStore } from '@/stores/audioStore';
import { useUserStore } from '@/stores/userStore';
import type { Song } from '@/script/types';

const props = defineProps<{ song: Song }>();
const userStore = useUserStore();
const audioStore = useAudioStore();

function updateSong() {

  let updated = props.song;
  audioStore.setSong(updated)
}
</script>


<template>

  <div @click="updateSong" class="m-2 border border-pink-500 rounded-lg flex">
    <img class="h-16 w-16 m-1 rounded-lg" :src="encodeURI(props.song.previewimage + '?h=64&w=64')" loading="lazy" />
    <div class="flex flex-col">
      <h3 class="text-nowrap overflow-scroll">
        <slot name="songName">{{ props.song.name }}</slot>
      </h3>
      <h5 class="text-yellow-500 text-sm">
        <slot name="artist">{{ props.song.artist }}</slot>
      </h5>
      <h5 class="text-yellow-500 text-sm">
        <slot name="length">{{ props.song.length }}</slot>
      </h5>

    </div>
  </div>

</template>

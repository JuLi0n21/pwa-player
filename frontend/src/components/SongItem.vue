<script setup lang="ts">
import { useAudio } from '@/composables/useAudio';
import type { Song } from '@/script/types';

const props = defineProps<{
  song: Song,
  action?: string,
  info?: string,
  border?: string,
}>();
const audioStore = useAudio();

function updateSong() {

  let updated = props.song;
  audioStore.setSong(updated)
}
</script>


<template>

  <div @click="updateSong" :style="{ borderColor: border }" class="m-1 md:text-xl border bordercolor rounded-lg flex">
    <img class="h-14 w-14 md:w-24 md:h-24 m-1 rounded-lg"
      :src="encodeURI(props.song?.previewimage ? props.song?.previewimage + '?h=56&w=56' : '/default-bg.png')"
      loading="lazy" />
    <div class="flex flex-col overflow-hidden text-left">
      <p :style="{ color: info }" class="text-nowrap text-ellipsis overflow-hidden text-base info">
        <slot name="songName">{{ props.song?.name ? props.song?.name : 'Unknown Title' }}</slot>
      </p>
      <h5 :style="{ color: action }" class="action text-sm text-nowrap text-ellipsis overflow-hidden text-base">
        <slot name="artist">{{ props.song?.artist ? props.song.artist : 'Unknown Artist' }}</slot>
      </h5>
      <h5 :style="{ color: action }" class="action text-sm">
        <slot name="length">{{ Math.floor(props.song?.length / 60000 ?? 0) }}:{{ Math.floor((props.song?.length ?? 0 /
          1000)
          % 60).toString().padStart(2, '0') }}</slot>
      </h5>

    </div>
  </div>

</template>

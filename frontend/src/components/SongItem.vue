<script setup lang="ts">
import { useAudio } from "@/composables/useAudio";
import type { Song } from "@/script/types";

const props = defineProps<{
  song: Song;
  action?: string;
  info?: string;
  border?: string;
}>();
const audioStore = useAudio();

function updateSong() {
  let updated = props.song;
  audioStore.setSong(updated);
}
</script>

<template>
  <div @click="updateSong" :style="{ borderColor: border }" class="flex m-1 border rounded-lg md:text-xl bordercolor">
    <img
      class="m-1 rounded-lg w-14 md:w-24 h-14 md:h-24"
      :src="encodeURI(props.song?.previewimage ? props.song?.previewimage + '?h=56&w=56' : '/default-bg.png')"
      loading="lazy"
    />
    <div class="flex flex-col overflow-hidden text-left">
      <p :style="{ color: info }" class="overflow-hidden text-base text-ellipsis text-nowrap info">
        <slot name="songName">{{ props.song?.name ? props.song?.name : "Unknown Title" }}</slot>
      </p>
      <h5 :style="{ color: action }" class="overflow-hidden text-sm text-base text-ellipsis text-nowrap action">
        <slot name="artist">{{ props.song?.artist ? props.song.artist : "Unknown Artist" }}</slot>
      </h5>
      <h5 :style="{ color: action }" class="text-sm action">
        <slot name="length"
          >{{ Math.floor(props.song.length / 60000 || 0) }}:{{
            Math.floor((props.song.length ?? 0 / 1000) % 60)
              .toString()
              .padStart(2, "0")
          }}</slot
        >
      </h5>
    </div>
  </div>
</template>

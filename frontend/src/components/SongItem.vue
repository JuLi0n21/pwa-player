<script setup lang="ts">
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";
import type { Song } from "@/script/types";

const props = defineProps<{
  song: Song;
  action?: string;
  info?: string;
  border?: string;
}>();
const audioStore = useAudio();
const userStore = useUser();
function updateSong() {
  let updated = props.song;
  audioStore.setSong(updated);
}
</script>

<template>
  <div 
    @click="updateSong" 
    :style="{ borderColor: border }" 
    class="flex items-center m-1 border rounded-lg h-16 md:h-24 overflow-hidden md:text-xl cursor-pointer bordercolor"
  >
    <img
      class="p-1 rounded-lg w-14 md:w-24 h-16 md:h-24 object-cover shrink-0"
      :src="props.song?.previewimage 
        ? encodeURI(`${userStore.cloudflareUrl.value}${props.song.previewimage}?h=56&w=56`) 
        : '/default-bg.png'"
      loading="lazy"
    />
    
    <div class="flex flex-col flex-1 justify-center overflow-hidden text-left">
      <p :style="{ color: info }" class="overflow-hidden text-base text-ellipsis text-nowrap info">
        <slot name="songName">{{ props.song?.name || "Unknown Title" }}</slot>
      </p>
      <h5 :style="{ color: action }" class="overflow-hidden text-sm text-ellipsis text-nowrap action">
        <slot name="artist">{{ props.song?.artist || "Unknown Artist" }}</slot>
      </h5>
      <h5 :style="{ color: action }" class="text-sm action">
        <slot name="length">
          {{ Math.floor(props.song?.length / 60000 || 0) }}:{{
            Math.floor((props.song?.length / 1000 || 0) % 60)
              .toString()
              .padStart(2, "0")
          }}
        </slot>
      </h5>
    </div>
  </div>
</template>
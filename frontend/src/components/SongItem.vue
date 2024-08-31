<script setup lang="ts">
import { ref, defineProps } from 'vue'
import { useAudioStore } from '@/stores/audioStore';
import { useUserStore } from '@/stores/userStore';
import type { Song } from '@/script/types';

const props = defineProps<{
  song: Song,
  action?: string,
  info?: string,
  border?: string,
}>();
const userStore = useUserStore();
const audioStore = useAudioStore();

function updateSong() {

  let updated = props.song;
  audioStore.setSong(updated)
}
</script>


<template>

  <div @click="updateSong" :style="{ borderColor: border }" class="m-1 border bordercolor rounded-lg flex">
    <img class="h-14 w-14 m-1 rounded-lg"
      :src="encodeURI(props.song?.previewimage ? props.song?.previewimage + '?h=56&w=56' : '/default-bg.png')"
      loading="lazy" />
    <div class="flex flex-col overflow-hidden">
      <p :style="{ color: info }" class="text-nowrap text-ellipsis overflow-hidden text-base info">
        <slot name="songName">{{ props.song?.name ?? 'Title' }}</slot>
      </p>
      <h5 :style="{ color: action }" class="action text-sm text-nowrap text-ellipsis overflow-hidden text-base">
        <slot name="artist">{{ props.song?.artist ?? 'Artist' }}</slot>
      </h5>
      <h5 :style="{ color: action }" class="action text-sm">
        <slot name="length">{{ Math.floor(props.song?.length ?? 0 / 60000) }}:{{ Math.floor((props.song?.length ?? 0 /
          1000)
          % 60).toString().padStart(2, '0') }}</slot>
      </h5>

    </div>
  </div>

</template>

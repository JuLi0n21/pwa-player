<template>

  <main class="flex-1 flex-col">
    <div class="flex-1 flex-col h-full overflow-y-hidden song-container">
        <p class="p-1 rounded-full backdrop--light shadow-xl text-center">History</p>
        <p v-if="songs.length === 0">No songs found</p>
      <SongItem v-for="(song, index) in songs" :key="song.hash" :song="song" />
    </div>
  </main>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import SongItem from '../components/SongItem.vue'

import { useAudio } from '@/composables/useAudio';
const audioStore = useAudio();

const songs = computed(() => Array.from(audioStore.recentlyPlayed.value?.values()).reverse() || [])

</script>
<script setup lang="ts">
import SongItem from "../components/SongItem.vue";
import SongItemSkeleton from "../components/SongItemSkeleton.vue";

import { type Song, mapApiToSongs } from "../script/types";
import { ref, onMounted, nextTick } from "vue";
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";
import { useApi } from "@/composables/useApi";

const userStore = useUser();
const audioStore = useAudio();
const { musicApi } = useApi();
const api = musicApi.value;

const songs = ref<Song[]>([]);
const name = ref("name");
const containerRef = ref<HTMLElement | null>(null);

const limit = ref(100);
const offset = ref(0);
const isLoading = ref(false);

const fetchRecent = async () => {
  if (isLoading.value) return;

  isLoading.value = true;
  try {
    const response = await api.musicBackendRecent(limit.value, offset.value);
    if (response.data.songs) {
      let newSongs = mapApiToSongs(response.data.songs);

      offset.value += limit.value;
      songs.value = [...songs.value, ...newSongs];

      isLoading.value = false;
      audioStore.setCollection(songs.value);
    }
  } catch (error) {
    console.error("Failed to load songs:", error);
  }
};

onMounted(async () => {
  await fetchRecent();

  await nextTick();

  const container = containerRef.value;
  if (container) {
    container.addEventListener("scroll", async () => {
      const scrollTop = container.scrollTop;
      const scrollHeight = container.scrollHeight;
      const clientHeight = container.clientHeight;

      if (scrollTop + clientHeight >= scrollHeight * 0.9 && !isLoading.value) {
        await fetchRecent();
      }
    });
  }
});
</script>

<template>
  <div 
    ref="containerRef" 
    class="flex-1 gap-2 grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 p-1 overflow-y-scroll song-container"
  >
    <SongItem v-for="(song, index) in songs" :key="index" :song="song" />

    <template v-if="isLoading">
      <SongItemSkeleton v-for="i in 5" :key="'skel-' + i" />
    </template>
  </div>
</template>

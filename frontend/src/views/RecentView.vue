<script setup lang="ts">
import SongItem from '../components/SongItem.vue'

import type { Song, CollectionPreview } from '../script/types'
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore';
import { useRoute } from 'vue-router';
import CollectionListItem from '../components/CollectionListItem.vue'
import { useAudioStore } from '@/stores/audioStore';

const route = useRoute();
const userStore = useUserStore();
const audioStore = useAudioStore();

const songs = ref<Song[]>([]);
const name = ref('name');

const limit = ref(100);
const offset = ref(0);
const isLoading = ref(false);

const fetchRecent = async () => {
  if (isLoading.value) return;

  isLoading.value = true;
  const data = await userStore.fetchRecent(limit.value, offset.value);

  data.forEach(song => {
    song.previewimage = `${userStore.baseUrl}api/v1/images/${song.previewimage}`;
    song.url = `${userStore.baseUrl}api/v1/audio/${song.url}`;
  });

  offset.value += limit.value;
  console.log(data)
  songs.value = [...songs.value, ...data];

  isLoading.value = false;
  audioStore.setCollection(null);
}


onMounted(async () => {
  await fetchRecent();

  const container = document.querySelector('.song-container');
  if (container) {
    container.addEventListener('scroll', async () => {
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

  <main class="flex-1 flex-col overflow-scroll">
    <div class="flex-1 flex-col h-full overflow-scroll song-container">

      <SongItem v-for="(song, index) in songs" :key="index" :song="song" />
    </div>
  </main>
</template>

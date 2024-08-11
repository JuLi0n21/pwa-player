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

onMounted(async () => {
  const data = await userStore.fetchCollection(Number(route.params.id));
  console.log(data)
  
  data.songs.forEach(song => {
    song.previewimage = `${userStore.baseUrl}api/v1/images/${song.previewimage}`;
    song.url = `${userStore.baseUrl}api/v1/audio/${song.url}`;
  });

  name.value = data.name;
  songs.value = data.songs;

  audioStore.setCollection(songs.value)
});

</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start my-2 mx-1 space-x-1">
        <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/menu/collections"><i
            class="fa-solid fa-arrow-left"></i>
        </RouterLink>
        <h1 class="px-8 text-nowrap overflow-scroll absolute left-0 right-0 text-center"> {{ name }} </h1>

      </nav>
      <hr>
    </div>
  </header>

  <div class="flex-1 flex-col h-full overflow-scroll">

      <SongItem 
        v-for="(song, index) in songs"
        :key="index"
        :song="song"
      />
  </div>
</template>

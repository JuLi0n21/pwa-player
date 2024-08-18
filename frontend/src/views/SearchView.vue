<script setup lang="ts">
import type { Song, CollectionPreview } from '../script/types'
import { useUserStore } from '@/stores/userStore';
import { onMounted, ref, watch } from 'vue';
import ActiveSearchList from '../components/ActiveSearchList.vue'
import { useRoute, useRouter } from 'vue-router';
import SongItem from '../components/SongItem.vue'
import { useAudioStore } from '@/stores/audioStore';

const router = useRouter();
const route = useRoute();
const audioStore = useAudioStore();
const userStore = useUserStore();
const activesongs = ref<Song[]>([]);
const songs = ref<Song[]>([]);
const artists = ref<string[]>([]);
const showSearch = ref(false);



onMounted(async () => {

await loadartistifexist();

  const container = document.querySelector('.search') as HTMLInputElement;

  if (container) {
    container.addEventListener('input', async (event: Event) => {
      showSearch.value = true;

      const target = event.target as HTMLInputElement;
        if(target.value != undefined && target.value != ""){
         const data = await userStore.fetchActiveSearch(target.value)
         router.push({ query: {s: target.value } });

         data.songs.forEach(song => {
            song.previewimage = `${userStore.baseUrl}api/v1/images/${song.previewimage}`;
            song.url = `${userStore.baseUrl}api/v1/audio/${song.url}`;
          });
          
          activesongs.value = data.songs;
          audioStore.setCollection(data.songs)
          artists.value = data.artist;
      } else {
          activesongs.value = [];
          artists.value = [];
          showSearch.value = false;
         
      }
    }
    )}
      const s =  route.query.s as string;
      if(s){container.value = s; container.dispatchEvent(new Event('input'))}
  });

async function loadartistifexist(){
  const query = route.query.a as string;
  if (query) {
    showSearch.value = false;

    const data = await userStore.fetchSearchArtist(query)
    console.log(data);

    data.forEach(song => {
            song.previewimage = `${userStore.baseUrl}api/v1/images/${song.previewimage}`;
            song.url = `${userStore.baseUrl}api/v1/audio/${song.url}`;
          });

    songs.value = data;
  }
}

watch(() => route.query.a, async (newQuery) => {
  
  await loadartistifexist();

});

</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start my-2 mx-1 space-x-1">
        <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/"><i class="fa-solid fa-arrow-left"></i>
        </RouterLink>
        <h1 class="absolute left-0 right-0 text-center"> Search </h1>
      </nav>
      <hr>
    </div>
  </header>
  <main class="flex flex-col flex-1 flex-col w-full h-full overflow-scroll">
    <input placeholder="Type to Search..." class="flex-1 max-h-12 search border border-pink-500 accent-pink-800 bg-yellow-300 bg-opacity-20 rounded-lg m-2 p-2" />
    <div class="relative flex flex-col w-full h-full overflow-scroll">
      <div v-if="showSearch" class="absolute w-full text-center search-recommendations  z -20">
         <ActiveSearchList :songs="activesongs" :artist="artists"/>
      </div>    
    <SongItem v-for="(song, index) in songs" :key="index" :song="song" />
    </div>
  </main>
</template>

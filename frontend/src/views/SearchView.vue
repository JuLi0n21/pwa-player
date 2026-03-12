<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";
import { useApi } from "@/composables/useApi";
import { mapApiToSongs, type Song } from "../script/types";

// Components
import ActiveSearchList from "../components/ActiveSearchList.vue";
import SongItem from "../components/SongItem.vue";
import SongItemSkeleton from "../components/SongItemSkeleton.vue";
import ActiveSearchSkeleton from "@/components/ActiveSearchSkeleton.vue";

const router = useRouter();
const route = useRoute();
const audioStore = useAudio();
const userStore = useUser();
const { musicApi } = useApi();
const api = musicApi.value;

// State
const songs = ref<Song[]>([]);
const activesongs = ref<Song[]>([]);
const artists = ref<string[]>([]);
const searchInput = ref((route.query.s as string) || "");
const isLoading = ref(false);
const showSearch = ref(false);

async function fetchActiveSearch(term: string) {
  if (!term.trim()) return emptySearch();
  
  isLoading.value = true;
  try {
    const response = await api.musicBackendSearch(term);
    const songData = mapApiToSongs(response.data.songs ?? []);
    
    activesongs.value = songData;
    artists.value = response.data.artist ? [response.data.artist] : [];
    
    audioStore.setCollection(songData);
    showSearch.value = true;
    router.replace({ query: { ...route.query, s: term } });
  } finally {
    isLoading.value = false;
  }
}

/**
 * Fetch full song list for a specific Artist
 */
async function fetchSearchArtist(artist: string) {
  isLoading.value = true;
  showSearch.value = false; // Hide recommendations overlay when a choice is made
  
  try {
    const response = await api.musicBackendArtist(artist);
    const data = mapApiToSongs(response.data.songs || []);

    songs.value = data.map((song: Song) => ({
      ...song,
      previewimage: `${userStore.cloudflareUrl.value}/api/v1/images/${song.previewimage}`,
      url: `${userStore.cloudflareUrl.value}/api/v1/audio/${song.url}`
    }));
    
    router.replace({ query: { ...route.query, a: artist } });
  } finally {
    isLoading.value = false;
  }
}

async function emptySearch() {
  activesongs.value = [];
  artists.value = [];
  songs.value = [];
  showSearch.value = false;
  searchInput.value = "";
  router.replace({ query: {} });
}

let debounceTimeout: any;
watch(searchInput, (val) => {
  clearTimeout(debounceTimeout);
  if (val && val.trim() !== "") {
    debounceTimeout = setTimeout(() => fetchActiveSearch(val), 300);
  } else {
    emptySearch();
  }
});

watch(() => route.query.a, (newArtist) => {
  if (newArtist) fetchSearchArtist(newArtist as string);
});

onMounted(() => {
  if (route.query.a) fetchSearchArtist(route.query.a as string);
  else if (route.query.s) fetchActiveSearch(route.query.s as string);
});
</script>

<template>
  <header class="top-0 z-30 sticky bg-black/10 backdrop-blur-md">
    <div class="p-2 wrapper">
      <nav class="relative flex items-center h-10">
        <RouterLink class="z-10 bg-white/5 shadow-xl p-2 rounded-full" to="/">
          <i class="fa-arrow-left fa-solid"></i>
        </RouterLink>
        <h1 class="absolute inset-0 flex justify-center items-center font-bold text-xl">Search</h1>
      </nav>
      <hr class="opacity-10 mt-2" />
    </div>
  </header>

  <main class="flex flex-col flex-1 w-full h-full overflow-hidden">
    <div class="relative p-2">
      <input
        v-model="searchInput"
        placeholder="Type to Search..."
        class="flex-1 bg-white/5 p-4 border rounded-xl outline-none ring-yellow-500/50 focus:ring-2 w-full h-14 transition-all bordercolor"
      />
      <div 
        v-if="searchInput"
        class="top-1/2 right-6 absolute opacity-50 hover:opacity-100 -translate-y-1/2 cursor-pointer" 
        @click="emptySearch"
      >
        <i class="text-xl far fa-times-circle"></i>
      </div>
    </div>

    <div class="relative flex-1 overflow-y-auto">
      
      <div 
        v-if="showSearch && (activesongs.length || artists.length || isLoading)" 
        class="z-20 absolute backdrop-blur-xl w-full min-h-full"
      >
        <ActiveSearchSkeleton v-if="isLoading" />

        <ActiveSearchList 
          v-else 
          :songs="activesongs" 
          :artist="artists" 
          :search="searchInput" 
        />
      </div>

        <template v-else>
          <SongItem 
            v-for="(song, index) in songs" 
            :key="song.hash || index" 
            :song="song"
            class="song-render-node"
          />
        </template>
        
        <div v-if="!isLoading && songs.length === 0 && !showSearch" class="col-span-full opacity-30 py-20 text-center">
          <i class="mb-4 text-6xl fa-solid fa-music"></i>
          <p>Find your favorite music</p>
        </div>
      </div>
  </main>
</template>

<style scoped>
.song-render-node {
  content-visibility: auto;
  contain-intrinsic-size: 96px;
}
</style>
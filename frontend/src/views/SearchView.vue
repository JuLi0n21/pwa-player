<script setup lang="ts">
import { mapApiToSongs, mapToSong, type Song } from "../script/types";
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import ActiveSearchList from "../components/ActiveSearchList.vue";
import SongItem from "../components/SongItem.vue";
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";
import { useApi } from "@/composables/useApi";

const router = useRouter();
const route = useRoute();

const audioStore = useAudio();
const userStore = useUser();
const { musicApi } = useApi();
const api = musicApi.value;

const activesongs = ref<Song[]>([]);
const songs = ref<Song[]>([]);
const artists = ref<string[]>([]);
const showSearch = ref(false);
const searchTerm = ref("");

async function fetchActiveSearch(term: string) {
  const response = await api.musicBackendSearch(term);

  const songData = mapApiToSongs(response.data.songs ?? []);

  activesongs.value = songData;

  if (response.data.artist) artists.value = [response.data.artist];
  audioStore.setCollection(songData);
  showSearch.value = true;
  searchTerm.value = term;
  router.replace({ query: { s: term } });
}

async function fetchSearchArtist(artist: string) {
  const response = await api.musicBackendArtist(artist);

  const data = mapApiToSongs(response.data.songs || []);

  data.forEach((song: Song) => {
    song.previewimage = `${userStore.cloudflareUrl.value}/api/v1/images/${song.previewimage}`;
    song.url = `${userStore.cloudflareUrl.value}/api/v1/audio/${song.url}`;
  });

  songs.value = data;
  showSearch.value = false;
}

async function emptySearch() {
  activesongs.value = [];
  artists.value = [];
  songs.value = [];
  showSearch.value = false;
  searchTerm.value = "";
  router.replace({ query: {} });
  searchInput.value = "";
}

onMounted(async () => {
  if (route.query.a) {
    await fetchSearchArtist(route.query.a as string);
  }
  if (route.query.s) {
    await fetchActiveSearch(route.query.s as string);
  }
});

watch(
  () => route.query.a,
  async (newArtist) => {
    if (newArtist) {
      await fetchSearchArtist(newArtist as string);
    } else {
      songs.value = [];
    }
  },
);

const searchInput = ref(searchTerm.value);

watch(searchInput, async (val) => {
  if (val && val.trim() !== "") {
    await fetchActiveSearch(val);
  } else {
    showSearch.value = false;
    activesongs.value = [];
    artists.value = [];
    router.replace({ query: {} });
  }
});
</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start space-x-1 mx-1 my-2">
        <RouterLink class="shadow-xl backdrop--light p-1 rounded-full" to="/">
          <i class="fa-arrow-left fa-solid"></i>
        </RouterLink>
        <h1 class="right-0 left-0 absolute text-center">Search</h1>
      </nav>
      <hr />
    </div>
  </header>

  <main class="flex flex-col flex-1 w-full h-full">
    <div class="relative">
      <input
        v-model="searchInput"
        placeholder="Type to Search..."
        class="flex-1 bg-yellow-300 bg-opacity-20 p-2 border rounded-lg w-full max-h-12 accent-pink-800 search bordercolor"
      />
      <div class="top-4 right-4 absolute flex flex-col justify-center cursor-pointer" @click="emptySearch">
        <i class="opacity-50 far fa-times-circle"></i>
      </div>
    </div>

    <div class="relative flex flex-col w-full h-full overflow-y-scroll">
      <div v-if="showSearch" class="z-20 absolute w-full text-center search-recommendations">
        <ActiveSearchList :songs="activesongs" :artist="artists" :search="searchTerm" />
      </div>
      <SongItem v-for="(song, index) in songs" :key="index" :song="song" />
    </div>
  </main>
</template>

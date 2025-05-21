<script setup lang="ts">
import { mapApiToSongs, type Song } from '../script/types'
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ActiveSearchList from '../components/ActiveSearchList.vue'
import SongItem from '../components/SongItem.vue'
import { useAudio } from '@/composables/useAudio'
import { useUser } from '@/composables/useUser'
import { MusicBackendApi } from '@/generated'
import { useApi } from '@/composables/useApi'

const router = useRouter()
const route = useRoute()

const audioStore = useAudio()
const userStore = useUser()
const { musicApi } = useApi()
const api = musicApi()

const activesongs = ref<Song[]>([])
const songs = ref<Song[]>([])
const artists = ref<string[]>([])
const showSearch = ref(false)
const searchTerm = ref('')

async function fetchActiveSearch(term: string) {
  const response = await api.musicBackendSearch(term);

  const songData = mapApiToSongs(response.data.songs)
  
  songData.forEach((song: Song) => {
    song.previewimage = `${userStore.baseUrl.value}/api/v1/images/${song.previewimage}`
    song.url = `${userStore.baseUrl.value}/api/v1/audio/${song.url}`
  })

  activesongs.value = songData
  
  if (response.data.artist)  artists.value = [response.data.artist]
  audioStore.setCollection(songData)
  showSearch.value = true
  searchTerm.value = term
  router.replace({ query: { s: term } })
}

async function fetchSearchArtist(artist: string) {
  const data = await musicApi.MusicBackend_SearchArtists(artist)

  data.forEach((song: Song) => {
    song.previewimage = `${userStore.baseUrl.value}/api/v1/images/${song.previewimage}`
    song.url = `${userStore.baseUrl.value}/api/v1/audio/${song.url}`
  })
  
  songs.value = data
  showSearch.value = false
}

async function emptySearch() {
  activesongs.value = []
  artists.value = []
  songs.value = []
  showSearch.value = false
  searchTerm.value = ''
  router.replace({ query: {} })
}

onMounted(async () => {
  if (route.query.a) {
    await fetchSearchArtist(route.query.a as string)
  }
  if (route.query.s) {
    await fetchActiveSearch(route.query.s as string)
  }
})

// Watch for artist query changes
watch(() => route.query.a, async (newArtist) => {
  if (newArtist) {
    await fetchSearchArtist(newArtist as string)
  } else {
    songs.value = []
  }
})

// Search input model
const searchInput = ref(searchTerm.value)

watch(searchInput, async (val) => {
  if (val && val.trim() !== '') {
    await fetchActiveSearch(val)
  } else {
    showSearch.value = false
    activesongs.value = []
    artists.value = []
    router.replace({ query: {} })
  }
})
</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start my-2 mx-1 space-x-1">
        <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/">
          <i class="fa-solid fa-arrow-left"></i>
        </RouterLink>
        <h1 class="absolute left-0 right-0 text-center">Search</h1>
      </nav>
      <hr />
    </div>
  </header>

  <main class="flex flex-col flex-1 w-full h-full overflow-scroll">
    <div class="relative">
      <input
        v-model="searchInput"
        placeholder="Type to Search..."
        class="flex-1 max-h-12 search border bordercolor accent-pink-800 bg-yellow-300 bg-opacity-20 rounded-lg m-2 p-2"
      />
      <div class="absolute h-16 right-4 flex flex-col justify-center cursor-pointer" @click="emptySearch">
        <i class="far fa-times-circle opacity-50"></i>
      </div>
    </div>

    <div class="relative flex flex-col w-full h-full overflow-scroll">
      <div v-if="showSearch" class="absolute w-full text-center search-recommendations z-20">
        <ActiveSearchList :songs="activesongs" :artist="artists" :search="searchTerm" />
      </div>
      <SongItem v-for="(song, index) in songs" :key="index" :song="song" />
    </div>
  </main>
</template>

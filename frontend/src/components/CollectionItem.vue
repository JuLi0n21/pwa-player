<script setup lang="ts">
import SongItem from '../components/SongItem.vue'
import { type Song, mapApiToCollection } from '../script/types'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import CollectionListItem from '../components/CollectionListItem.vue'
import { useAudio } from '@/composables/useAudio'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const audioStore = useAudio()
const { musicApi } = useApi()

const songs = ref<Song[]>([])
const name = ref('name')
const limit = 50
let offset = ref(0)
let loading = ref(false)
let hasMore = ref(true)

const fetchMoreSongs = async () => {
  if (loading.value || !hasMore.value) return

  loading.value = true
  try {
    const response = await musicApi().musicBackendCollections(
      "",
      Array.isArray(route.params.id) ? route.params.id[0] : route.params.id,
      limit,
      offset.value
    )

    const col = mapApiToCollection(response.data)

    if (offset.value === 0) {
      name.value = col.name
      songs.value = col.songs
    } else {
      songs.value.push(...col.songs)
    }

    audioStore.setCollection(songs.value)

    if (col.songs.length < limit) {
      hasMore.value = false
    }

    offset.value += limit
  } catch (error) {
    console.error('Error fetching songs:', error)
  } finally {
    loading.value = false
  }
}

const onScroll = (e: Event) => {
  const target = e.target as HTMLElement
  const scrollTop = target.scrollTop
  const scrollHeight = target.scrollHeight
  const clientHeight = target.clientHeight

  if (scrollTop + clientHeight >= scrollHeight * 0.9) {
    fetchMoreSongs()
  }
}

onMounted(() => {
  fetchMoreSongs()
  const container = document.querySelector('.coll-container')
  container?.addEventListener('scroll', onScroll)
})

onBeforeUnmount(() => {
  const container = document.querySelector('.coll-container')
  container?.removeEventListener('scroll', onScroll)
})
</script>


<template>
 <main class="flex-1 text-center flex flex-col h-full overflow-y-hidden">
    <header v-show="true">
      <div class="wrapper">
        <nav class="flex justify-start my-2 mx-1 space-x-1 overflow-y-scroll flex-nowrap text-nowrap">
          <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/menu/collections"><i class="fa-solid fa-arrow-left"></i>
          </RouterLink>
        <p  class="p-1 rounded-full backdrop--light shadow-xl">{{ name }}</p>
        </nav>
        <hr>
      </div>
    </header>
    <div
      class="flex-1 flex-col  overflow-y-scroll coll-container"
    >
      <SongItem v-for="(song, index) in songs" :key="index" :song="song" />
    </div>
 </main>
</template>
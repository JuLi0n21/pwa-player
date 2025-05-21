<script setup lang="ts">
import SongItem from '../components/SongItem.vue'
import type { Song } from '../script/types'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import CollectionListItem from '../components/CollectionListItem.vue'
import { useAudio } from '@/composables/useAudio'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const audioStore = useAudio()
const { musicApi } = useApi()

const songs = ref<Song[]>([])
const name = ref('name')

onMounted(async () => {
  try {
    const response = await musicApi().musicBackendCollections(Array.isArray(route.params.id) ? route.params.id[0] : route.params.id) // Adjust method name if needed

    const base = import.meta.env.VITE_MUSIC_API_URL || 'http://localhost:8080'

    response.songs.forEach(song => {
      song.previewimage = `${base}/api/v1/images/${song.previewimage}`
      song.url = `${base}/api/v1/audio/${song.url}`
    })

    name.value = response.name
    songs.value = response.songs

    audioStore.setCollection(songs.value)
  } catch (error) {
    console.error('Error fetching collection:', error)
  }
})
</script>

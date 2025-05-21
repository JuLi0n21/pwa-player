import { ref, onMounted } from 'vue'
import { mapApiToSongs, type Song } from '@/script/types'
import { useApi } from './useApi'

let audioInstance: ReturnType<typeof createAudio> | null = null

function createAudio() {
  const { musicApi } = useApi()

  const songSrc = ref('https://cdn.pixabay.com/audio/2024/05/24/audio_46382ae035.mp3')
  const artist = ref('Artist')
  const title = ref('Title')
  const bgimg = ref('https://assets.ppy.sh/beatmaps/2197744/covers/cover@2x.jpg?1722207959')
  const hash = ref('0000')

  const isPlaying = ref(false)
  const duration = ref('0:00')
  const currentTime = ref('0:00')
  const percentDone = ref(0)

  const shuffle = ref(false)
  const repeat = ref(false)

  const activeCollection = ref<Song[]>([])
  const currentSong = ref<Song | null>(null)

  function saveSongToLocalStorage(song: Song) {
    console.log(song)
    localStorage.setItem('lastPlayedSong', JSON.stringify(song))
  }

  function loadSongFromLocalStorage(): Song | null {
    const song = localStorage.getItem('lastPlayedSong')
    return song ? JSON.parse(song) : null
  }

  function saveCollectionToLocalStorage(collection: Song[]) {
    localStorage.setItem('lastActiveCollection', JSON.stringify(collection))
  }

  function loadCollectionFromLocalStorage(): Song[] | null {
    const collection = localStorage.getItem('lastActiveCollection')
    return collection ? JSON.parse(collection) : null
  }

  function togglePlay() {
    const audio = document.getElementById('audio-player') as HTMLAudioElement
    if (!audio) return
    if (audio.paused) {
      audio.play()
    } else {
      audio.pause()
    }
  }

  function update() {
    const audio = document.getElementById('audio-player') as HTMLAudioElement
    if (!audio) return

    isPlaying.value = !audio.paused

    const current_min = Math.floor(audio.currentTime / 60)
    const current_sec = Math.floor(audio.currentTime % 60)
    if (!isNaN(current_min) && !isNaN(current_sec)) {
      currentTime.value = `${current_min}:${current_sec.toString().padStart(2, '0')}`
    }

    const duration_min = Math.floor(audio.duration / 60)
    const duration_sec = Math.floor(audio.duration % 60)
    if (!isNaN(duration_min) && !isNaN(duration_sec)) {
      duration.value = `${duration_min}:${duration_sec.toString().padStart(2, '0')}`
    }

    const percent = (audio.currentTime / audio.duration) * 100
    if (!isNaN(percent)) {
      percentDone.value = percent
    }

    if (audio.ended) {
      next()
    }
  }

  function next() {
    const audio = document.getElementById('audio-player') as HTMLAudioElement
    if (!audio) return

    if (repeat.value) {
      audio.pause()
      audio.currentTime = 0
      audio.play()
      return
    }

    if (shuffle.value) {
      audio.pause()
      const randomSong = activeCollection.value[Math.floor(Math.random() * activeCollection.value.length)]
      setSong(randomSong)
      audio.play()
      return
    }

    toggleNext()
  }

  function updateTime() {
    const audio = document.getElementById('audio-player') as HTMLAudioElement
    const slider = document.getElementById('audio-slider') as HTMLInputElement
    if (!audio || !slider) return

    audio.currentTime = (Number(slider.value) / 100) * audio.duration
  }

  function togglePrev() {
    const index = activeCollection.value.findIndex((s) => s.hash === hash.value)
    if (index === -1) return
    const prevIndex = (index - 1 + activeCollection.value.length) % activeCollection.value.length
    setSong(activeCollection.value[prevIndex])
  }

  function toggleNext() {
    let index = 0
    if (shuffle.value) {
      index = Math.floor(Math.random() * activeCollection.value.length)
    } else {
      index = activeCollection.value.findIndex((s) => s.hash === hash.value)
      if (index === -1) return
      index = (index + 1) % activeCollection.value.length
    }
    setSong(activeCollection.value[index])
  }

  function toggleShuffle() {
    shuffle.value = !shuffle.value
  }

  function toggleRepeat() {
    repeat.value = !repeat.value
  }

  function setSong(song: Song | null) {
    if (!song) return

    songSrc.value = song.url
    artist.value = song.artist
    title.value = song.name
    bgimg.value = song.previewimage
    hash.value = song.hash

    currentSong.value = song
    saveSongToLocalStorage(song)

    const audio = document.getElementById('audio-player') as HTMLAudioElement
    if (!audio) return
    if (!audio.paused) audio.pause()
    audio.src = song.url

    audio.addEventListener('canplaythrough', () => {
      audio.play().catch(console.error)
    }, { once: true })
  }

  function setCollection(songs: Song[] | null) {
    if (!songs) return
    activeCollection.value = songs
    saveCollectionToLocalStorage(songs)
  }

  async function loadInitialCollection() {
    try {
      const api = musicApi()
      const response = await api.musicBackendRecent()
      if (response.data.songs) {
        setCollection(mapApiToSongs(response.data.songs))
      }
    } catch (error) {
      console.error('Failed to load songs:', error)
    }
  }

  function init() {
    setSong(loadSongFromLocalStorage())
    setCollection(loadCollectionFromLocalStorage())

    console.log(activeCollection.value)
    console.log(currentSong.value)
    if (!currentSong.value || activeCollection.value.length === 0) {
      loadInitialCollection()
    }
  }

  onMounted(() => {
    init()
  })

  return {
    songSrc,
    artist,
    title,
    bgimg,
    hash,
    isPlaying,
    duration,
    currentTime,
    percentDone,
    shuffle,
    repeat,
    activeCollection,
    currentSong,
    togglePlay,
    update,
    next,
    updateTime,
    togglePrev,
    toggleNext,
    toggleShuffle,
    toggleRepeat,
    setSong,
    setCollection,
    init,
  }
}

export function useAudio() {
  if (!audioInstance) {
    audioInstance = createAudio()
  }
  return audioInstance
}
